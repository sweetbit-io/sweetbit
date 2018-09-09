package main

import (
	"os"
	"os/signal"
	"log"
	"github.com/davidknezic/sweetd/machine"
)

type Dispenser struct {
	shouldBuzzOnDispense  bool
	shouldDispenseOnTouch bool
}

func main() {
	dispenser := Dispenser{
		shouldBuzzOnDispense:  true,
		shouldDispenseOnTouch: true,
	}

	machine := machine.NewMachine()

	defer machine.Stop()
	machine.Start()

	signals := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(signals, os.Interrupt)

	go func() {
		sig := <-signals
		log.Println(sig)
		log.Println("Received an interrupt, stopping services...")
		close(done)
	}()

	for {
		select {
		case on := <-machine.TouchEvents:
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
