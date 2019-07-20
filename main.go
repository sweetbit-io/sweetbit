package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/the-lightning-land/sweetd/ap"
	"github.com/the-lightning-land/sweetd/dispenser"
	"github.com/the-lightning-land/sweetd/machine"
	"github.com/the-lightning-land/sweetd/pairing"
	"github.com/the-lightning-land/sweetd/pos"
	"github.com/the-lightning-land/sweetd/sweetdb"
	"github.com/the-lightning-land/sweetd/sweetlog"
	"github.com/the-lightning-land/sweetd/sweetrpc"
	"github.com/the-lightning-land/sweetd/updater"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

var (
	// commit stores the current commit hash of this build. This should be set using -ldflags during compilation.
	commit string
	// version stores the version string of this build. This should be set using -ldflags during compilation.
	version string
	// date stores the date of this build. This should be set using -ldflags during compilation.
	date string
)

// sweetdMain is the true entry point for sweetd. This is required since defers
// created in the top-level scope of a main method aren't executed if os.Exit() is called.
func sweetdMain() error {
	sweetLog := sweetlog.New()

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.AddHook(sweetLog)

	log.Debug("Starting sweetd...")
	log.Debug("Loading config...")

	// Load CLI configuration and defaults
	cfg, err := loadConfig()
	if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
		return nil
	} else if err != nil {
		return errors.Errorf("Failed parsing arguments: %v", err)
	}

	// Set logger into debug mode if called with --debug
	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
		log.Info("Setting debug mode.")
	}

	log.Debug("Loaded config.")

	// Print version of the daemon
	log.Infof("Version %s (commit %s)", version, commit)
	log.Infof("Built on %s", date)

	// Stop here if only version was requested
	if cfg.ShowVersion {
		return nil
	}

	// sweet.db persistently stores all dispenser configurations and settings
	sweetDB, err := sweetdb.Open(cfg.DataDir)
	if err != nil {
		return errors.Errorf("Could not open sweet.db: %v", err)
	}

	log.Infof("Opened sweet.db")

	// The network access point, which acts as the core connectivity
	// provider for all other components
	var a ap.Ap

	switch cfg.Net {
	case "dispenser":
		a, err = ap.NewDispenserAp(&ap.DispenserApConfig{
			Interface: "wlan0",
		})

		log.Info("Created Candy Dispenser access point.")
	case "mock":
		a = ap.NewMockAp()

		log.Info("Created a mock access point.")
	default:
		return errors.Errorf("Unknown networking type %v", cfg.Machine)
	}

	err = a.Start()
	if err != nil {
		return errors.Errorf("Could not start access point: %v", err)
	}

	defer func() {
		err := a.Stop()
		if err != nil {
			log.Errorf("Could not properly shut down access point: %v", err)
		}
	}()

	wifiConnection, err := sweetDB.GetWifiConnection()
	if err != nil {
		log.Warnf("Could not retrieve saved wifi connection: %v", err)
	}

	if wifiConnection != nil {
		log.Infof("Will attempt connecting to Wifi %v.", wifiConnection.Ssid)

		err := a.ConnectWifi(wifiConnection.Ssid, wifiConnection.Psk)
		if err != nil {
			log.Warnf("Whoops, couldn't connect to wifi: %v", err)
		}
	} else {
		log.Infof("No saved Wifi connection available. Not connecting.")
	}

	// The updater
	var u updater.Updater

	switch cfg.Updater {
	case "none":
		u = updater.NewNoopUpdater()

		log.Info("Created noop updater.")
	case "mender":
		u, err = updater.NewMenderUpdater(&updater.MenderUpdaterConfig{
			Logger: log.New().WithField("system", "updater"),
		})

		log.Info("Created Mender updater.")
	default:
		return errors.Errorf("Unknown updater type %v", cfg.Updater)
	}

	artifactName, err := u.GetArtifactName()
	if err != nil {
		log.Error("Could not obtain artifact name. Continuing though...")
	} else {
		log.Infof("Updater returned artifact name %v", artifactName)
	}

	// The hardware controller
	var m machine.Machine

	switch cfg.Machine {
	case "raspberry":
		m = machine.NewDispenserMachine(&machine.DispenserMachineConfig{
			TouchPin:  cfg.Raspberry.TouchPin,
			MotorPin:  cfg.Raspberry.MotorPin,
			BuzzerPin: cfg.Raspberry.BuzzerPin,
		})

		log.Infof("Created Raspberry Pi machine on touch pin %v, motor pin %v and buzzer pin %v.",
			cfg.Raspberry.TouchPin, cfg.Raspberry.MotorPin, cfg.Raspberry.BuzzerPin)
	case "mock":
		m = machine.NewMockMachine(cfg.Mock.Listen)

		log.Info("Created a mock machine.")
	default:
		return errors.Errorf("Unknown machine type %v", cfg.Machine)
	}

	log.Infof("Creating PoS...")

	// create subsystem responsible for the point of sale app
	pos, err := pos.NewPos(&pos.Config{
		Logger:     log.New().WithField("system", "pos"),
		TorDataDir: cfg.Tor.DataDir,
	})
	if err != nil {
		return errors.Errorf("Could not create PoS: %v", err)
	}

	log.Infof("Created PoS")

	// central controller for everything the dispenser does
	dispenser := dispenser.NewDispenser(&dispenser.Config{
		Machine:     m,
		AccessPoint: a,
		DB:          sweetDB,
		MemoPrefix:  cfg.MemoPrefix,
		Updater:     u,
		Pos:         pos,
		SweetLog:    sweetLog,
		Logger:      log.New().WithField("system", "dispenser"),
	})

	log.Infof("Created dispenser.")

	log.Infof("Creating pairing controller...")

	// create subsystem responsible for pairing
	pairingController, err := pairing.NewController(&pairing.Config{
		Logger:      log.New().WithField("system", "pairing"),
		AdapterId:   "hci0",
		AccessPoint: a,
		Dispenser:   dispenser,
	})
	if err != nil {
		return errors.Errorf("Could not create pairing controller: %v", err)
	}

	log.Infof("Created pairing controller")

	log.Infof("Starting pairing controller...")

	err = pairingController.Start()
	if err != nil {
		return errors.Errorf("Could not start pairing controller: %v", err)
	}

	log.Infof("Started pairing controller")

	defer func() {
		log.Infof("Stopping pairing controller...")

		err := pairingController.Stop()
		if err != nil {
			log.Errorf("Could not properly shut down pairing controller: %v", err)
		}
	}()

	// Create a gRPC server for remote control of the dispenser
	if len(cfg.Listeners) > 0 {
		grpcServer := grpc.NewServer()

		sweetrpc.RegisterSweetServer(grpcServer, newRPCServer(dispenser, &rpcServerConfig{
			version: version,
			commit:  commit,
		}))

		// Next, Start the gRPC server listening for HTTP/2 connections.
		for _, listener := range cfg.Listeners {
			lis, err := net.Listen(listener.Network(), listener.String())
			if err != nil {
				return errors.Errorf("RPC server unable to listen on %s", listener.String())
			}

			defer func() {
				err := lis.Close()
				if err != nil {
					log.Errorf("Could not properly close: %v", err)
				}
			}()

			go func() {
				log.Infof("RPC server listening on %s", lis.Addr())
				err := grpcServer.Serve(lis)
				if err != nil {
					log.Errorf("Could not serve: %v", err)
				}
			}()
		}
	}

	// Handle interrupt signals correctly
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		sig := <-signals
		log.Info(sig)
		log.Info("Received an interrupt, stopping dispenser...")
		dispenser.Shutdown()
	}()

	// blocks until the dispenser is shut down
	err = dispenser.Run()
	if err != nil {
		return errors.Errorf("Failed running dispenser: %v", err)
	}

	// finish with no error
	return nil
}

func main() {
	// Call the "real" main in a nested manner so the defers will properly
	// be executed in the case of a graceful shutdown.
	if err := sweetdMain(); err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
		} else {
			log.WithError(err).Println("Failed running sweetd.")
		}
		os.Exit(1)
	}
}
