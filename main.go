package main

import (
	"os"
	"os/signal"
	"github.com/the-lightning-land/sweetd/machine"
	log "github.com/sirupsen/logrus"
)

type Dispenser struct {
	shouldBuzzOnDispense  bool
	shouldDispenseOnTouch bool
}

func main() {
	log.SetOutput(os.Stdout)

	dispenser := Dispenser{
		shouldBuzzOnDispense:  false,
		shouldDispenseOnTouch: true,
	}

	machine := machine.NewMachine()

	defer machine.Stop()
	machine.Start()

	signals := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(signals, os.Interrupt)

	go func() {
		sig := <-signals
		log.Info(sig)
		log.Info("Received an interrupt, stopping services...")
		done <- true
	}()

	for {
		select {
		case on := <-machine.TouchEvents:
			log.Info("Touch event {}", on)

			if dispenser.shouldDispenseOnTouch {
				machine.ToggleMotor(on)

				if dispenser.shouldBuzzOnDispense {
					machine.ToggleBuzzer(on)
				}
			} else if !on {
				machine.ToggleBuzzer(false)
				machine.ToggleMotor(false)
			}
		case <-done:
			return
		}
	}
}
