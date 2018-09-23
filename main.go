package main

import (
	"os"
	"os/signal"
	"github.com/the-lightning-land/sweetd/machine"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"github.com/the-lightning-land/sweetd/sweetrpc"
	"net"
	"github.com/the-lightning-land/sweetd/dnsmasq"
	"github.com/the-lightning-land/sweetd/hostapd"
	"time"
)

var (
	// Commit stores the current commit hash of this build. This should be set using -ldflags during compilation.
	Commit string
	// Version stores the version string of this build. This should be set using -ldflags during compilation.
	Version = "1.0.0"
	// Stores the configuration
	cfg *config
)

// TODO: remove me
type Dispenser struct {
	shouldBuzzOnDispense  bool
	shouldDispenseOnTouch bool
}

// TODO: Nest contents in another func so the defers will properly be executed in the case of a graceful shutdown.
func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	log.Info("Starting sweetd...")

	log.Info("Loading config...")

	cfg, err := loadConfig()
	if err != nil {
		log.WithError(err).Fatal("Could not read config.")
	}

	log.Info("Loaded config.")

	log.Infof("Version %s", Version)

	if cfg.RunAp {
		log.Infof("Setting up %s interface as access point...", cfg.Ap.Interface)

		if err := removeApInterface(cfg.Ap.Interface); err != nil {
			log.WithError(err).Fatal("Could not remove AP interface.")
		}

		if err := addApInterface(cfg.Ap.Interface); err != nil {
			log.WithError(err).Fatal("Could not add AP interface.")
		}

		if err := upApInterface(cfg.Ap.Interface); err != nil {
			log.WithError(err).Fatal("Could not up AP interface.")
		}

		if err := configureApInterface(cfg.Ap.Ip, cfg.Ap.Interface); err != nil {
			log.WithError(err).Fatal("Could not configure AP interface.")
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
			log.WithError(err).Fatal("Could not create hostapd service.")
		}

		if err := h.Start(); err != nil {
			log.WithError(err).Fatal("Could not start hostapd.")
		}

		defer h.Stop()

		log.Info("Started hostapd.")

		log.Info("Restarting dhcpd in order to reestablish previous connection...")

		if err := restartDhcp(); err != nil {
			log.WithError(err).Fatal("Could not restart dhcpd.")
		}

		log.Info("Restarted dhcpd.")
	} else {
		log.Info("Will not start access point according to configuration.")
	}

	if cfg.RunDnsmasq {
		log.Info("Starting dnsmasq for DNS and DHCP management...")

		d := dnsmasq.New(&dnsmasq.Config{
			Address:   cfg.Dnsmasq.Address,
			DhcpRange: cfg.Dnsmasq.DhcpRange,
			Log: func(s string) {
				log.WithField("service", "dnsmasq").Debug(s)
			},
		})

		if err := d.Start(); err != nil {
			log.WithError(err).Fatal("Could not start dnsmasq.")
		}

		defer d.Stop()

		log.Info("Started dnsmasq.")
	} else {
		log.Info("Will not start dnsmasq according to configuration.")
	}

	// TODO: remove me
	dispenser := Dispenser{
		shouldBuzzOnDispense:  false,
		shouldDispenseOnTouch: true,
	}

	var m machine.Machine

	if cfg.Machine == "raspberry" {
		log.Info("Using Raspberry Pi machine.")
		m = machine.NewDispenserMachine()
	} else if cfg.Machine == "mock" {
		log.Info("Using a mock machine.")
		m = machine.NewMockMachine(cfg.Mock.Listen)
	}

	log.Info("Starting machine...")
	if err := m.Start(); err != nil {
		log.WithError(err).Fatal("Could not start machine.")
	}
	defer m.Stop()

	signals := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(signals, os.Interrupt)

	go func() {
		sig := <-signals
		log.Info(sig)
		log.Info("Received an interrupt, stopping services...")
		done <- true
	}()

	grpcServer := grpc.NewServer()
	sweetrpc.RegisterSweetServer(grpcServer, newRPCServer(&rpcServerConfig{
		version: Version,
		commit:  Commit,
	}))

	// Next, Start the gRPC server listening for HTTP/2 connections.
	for _, listener := range cfg.Listeners {
		lis, err := net.Listen(listener.Network(), listener.String())
		if err != nil {
			log.Errorf("RPC server unable to listen on %s", listener)
			os.Exit(1)
		}
		defer lis.Close()
		go func() {
			log.Infof("RPC server listening on %s", lis.Addr())
			grpcServer.Serve(lis)
		}()
	}

	log.Info("sweetd started.")

	// Signal successful startup with two short buzzer noises
	m.ToggleBuzzer(true)
	time.Sleep(200 * time.Millisecond)
	m.ToggleBuzzer(false)
	time.Sleep(200 * time.Millisecond)
	m.ToggleBuzzer(true)
	time.Sleep(200 * time.Millisecond)
	m.ToggleBuzzer(false)

	for {
		select {
		case on := <-m.TouchEvents():
			log.Info("Touch event {}", on)

			if dispenser.shouldDispenseOnTouch {
				m.ToggleMotor(on)

				if dispenser.shouldBuzzOnDispense {
					m.ToggleBuzzer(on)
				}
			} else if !on {
				m.ToggleBuzzer(false)
				m.ToggleMotor(false)
			}
		case <-done:
			return
		}
	}
}
