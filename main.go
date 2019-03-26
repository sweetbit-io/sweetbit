package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/the-lightning-land/sweetd/ap"
	"github.com/the-lightning-land/sweetd/machine"
	"github.com/the-lightning-land/sweetd/sweetdb"
	"github.com/the-lightning-land/sweetd/sweetrpc"
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
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

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

	// The network access point, which acts as the core connectivity
	// provider for all other components
	var a ap.Ap

	if cfg.RunAp {
		a, err = ap.NewDispenserAp(&ap.DispenserApConfig{
			Interface: "wlan0",
			Hotspot: &ap.DispenserApHotspotConfig{
				Interface:  cfg.Ap.Interface,
				Ip:         cfg.Ap.Ip,
				DhcpRange:  cfg.Ap.DhcpRange,
				Ssid:       cfg.Ap.Ssid,
				Passphrase: cfg.Ap.Passphrase,
			},
		})

		log.Info("Created Candy Dispenser access point.")
	} else {
		a = ap.NewMockAp()

		log.Info("Created a mock access point.")
	}

	defer a.Stop()

	err = a.Start()
	if err != nil {
		return errors.Errorf("Could not start access point: %v", err)
	}

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

	err = a.StartHotspot()
	if err != nil {
		log.Warnf("Whoops, couldn't start hotspot: %v", err)
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

	log.Infof("Opened sweet.db")

	// central controller for everything the dispenser does
	dispenser := newDispenser(&dispenserConfig{
		machine:     m,
		accessPoint: a,
		db:          sweetDB,
		memoPrefix:  cfg.MemoPrefix,
	})

	log.Infof("Created dispenser.")

	// Handle interrupt signals correctly
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		sig := <-signals
		log.Info(sig)
		log.Info("Received an interrupt, stopping dispenser...")
		dispenser.shutdown()
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
				return errors.New("RPC server unable to listen on %s")
			}

			defer lis.Close()

			go func() {
				log.Infof("RPC server listening on %s", lis.Addr())
				grpcServer.Serve(lis)
			}()
		}
	}

	// blocks until the dispenser is shut down
	err = dispenser.run()
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
