package main

import (
	"os"
	"os/signal"
	"github.com/the-lightning-land/sweetd/machine"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"github.com/the-lightning-land/sweetd/sweetrpc"
	"net"
	"github.com/the-lightning-land/sweetd/hostapd"
	"github.com/the-lightning-land/sweetd/dnsmasq"
)

var (
	// Commit stores the current commit hash of this build. This should be set using -ldflags during compilation.
	Commit string
	// Version stores the version string of this build. This should be set using -ldflags during compilation.
	Version = "1.0.0"
	// Stores the configuration
	cfg *config
)

// Todo: remove me
type Dispenser struct {
	shouldBuzzOnDispense  bool
	shouldDispenseOnTouch bool
}

func main() {
	log.SetOutput(os.Stdout)

	log.Info("Starting sweetd...")

	log.Info("Loading config...")

	cfg, err := loadConfig()
	if err != nil {
		log.WithError(err).Fatal("Could not read config.")
	}

	log.Info("Loaded config.")

	log.Infof("Version %s", Version)

	if err := removeApInterface(); err != nil {
		log.WithError(err).Fatal("Could not remove AP interface.")
	}

	if err := addApInterface(); err != nil {
		log.WithError(err).Fatal("Could not add AP interface.")
	}

	if err := upApInterface(); err != nil {
		log.WithError(err).Fatal("Could not up AP interface.")
	}

	if err := configureApInterface(cfg.Ap.Ip); err != nil {
		log.WithError(err).Fatal("Could not configure AP interface.")
	}

	log.Info("Starting hostapd for access point management...")

	h := hostapd.New(&hostapd.Config{
		Ssid:       "hello",
		Passphrase: "world",
	})

	if err := h.Start(); err != nil {
		log.WithError(err).Fatal("Could not start hostapd.")
	}

	defer h.Stop()

	log.Info("Started hostapd.")

	if cfg.RunDnsmasq {
		log.Info("Starting dnsmasq for DNS and DHCP management...")

		d := dnsmasq.New(&dnsmasq.Config{
			Address:   cfg.Dnsmasq.Address,
			DhcpRange: cfg.Dnsmasq.DhcpRange,
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
	sweetrpc.RegisterSweetServer(grpcServer, newRPCServer(m.TouchEvents()))

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
