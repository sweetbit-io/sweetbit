package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/the-lightning-land/sweetd/dnsmasq"
	"github.com/the-lightning-land/sweetd/hostapd"
	"github.com/the-lightning-land/sweetd/machine"
	"github.com/the-lightning-land/sweetd/sweetdb"
	"github.com/the-lightning-land/sweetd/sweetrpc"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strings"
)

var (
	// Commit stores the current commit hash of this build. This should be set using -ldflags during compilation.
	commit string
	// Version stores the version string of this build. This should be set using -ldflags during compilation.
	version string
	// Stores the date of this build. This should be set using -ldflags during compilation.
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

	// Should run access point, so that the dispenser can be paired to
	// and controlled through an app?
	if cfg.RunAp {
		h, d, err := setUpAccessPoint(cfg)

		if h != nil {
			defer h.Stop()
		}

		if d != nil {
			defer d.Stop()
		}

		if err != nil {
			return errors.Errorf("Could not start access point: %v", err)
		}
	} else {
		log.Info("Will not start access point according to configuration.")
	}

	// The hardware controller
	var m machine.Machine

	switch cfg.Machine {
	case "raspberry":
		m = machine.NewDispenserMachine()

		log.Info("Created Raspberry Pi machine.")
	case "mock":
		m = machine.NewMockMachine(cfg.Mock.Listen)

		log.Info("Created a mock machine.")
	default:
		return errors.Errorf("Unknown machine type %v", cfg.Machine)
	}

	// sweet.db persistently stores all dispenser configurations and settings
	sweetDB, err := sweetdb.Open(cfg.DataDir)

	log.Infof("Opened sweet.db")

	// central controller of everything the dispenser does
	dispenser := newDispenser(m, sweetDB)

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

		sweetrpc.RegisterSweetServer(grpcServer, newRPCServer(&rpcServerConfig{
			version:   version,
			commit:    commit,
			dispenser: dispenser,
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

func setUpAccessPoint(cfg *config) (*hostapd.Hostapd, *dnsmasq.Dnsmasq, error) {
	log.Infof("Setting up %s interface as access point...", cfg.Ap.Interface)

	if err := removeApInterface(cfg.Ap.Interface); err != nil {
		return nil, nil, errors.Errorf("Could not remove AP interface: %v", err)
	}

	if err := addApInterface(cfg.Ap.Interface); err != nil {
		return nil, nil, errors.Errorf("Could not add AP interface: %v", err)
	}

	if err := upApInterface(cfg.Ap.Interface); err != nil {
		return nil, nil, errors.Errorf("Could not up AP interface: %v", err)
	}

	if err := configureApInterface(cfg.Ap.Ip, cfg.Ap.Interface); err != nil {
		return nil, nil, errors.Errorf("Could not configure AP interface: %v", err)
	}

	log.Info("Starting hostapd for access point management...")

	h, err := hostapd.New(&hostapd.Config{
		Ssid:       cfg.Ap.Ssid,
		Passphrase: cfg.Ap.Passphrase,
		Log: func(s string) {
			log.WithField("service", "hostapd").Debug(s)
		},
	})
	if err != nil {
		return nil, nil, errors.Errorf("Could not create hostapd service: %v", err)
	}

	if err := h.Start(); err != nil {
		return nil, nil, errors.Errorf("Could not start hostapd: %v", err)
	}

	log.Info("Started hostapd.")

	log.Info("Restarting dhcpd in order to reestablish previous connection...")

	if err := restartDhcp(); err != nil {
		return h, nil, errors.Errorf("Could not restart dhcpd: %v", err)
	}

	log.Info("Restarted dhcpd.")

	log.Info("Starting dnsmasq for DNS and DHCP management...")

	d, err := dnsmasq.New(&dnsmasq.Config{
		Address:   "/#/" + strings.Split(cfg.Ap.Ip, "/")[0],
		DhcpRange: cfg.Ap.DhcpRange,
		Log: func(s string) {
			log.WithField("service", "dnsmasq").Debug(s)
		},
	})

	if err != nil {
		return h, nil, errors.Errorf("Could not create dnsmasq service: %v", err)
	}

	if err := d.Start(); err != nil {
		return h, nil, errors.Errorf("Could not start dnsmasq: %v", err)
	}

	log.Info("Started dnsmasq.")

	return h, d, nil
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
