package main

import (
	"os"
	"os/signal"
	"log"
	"github.com/davidknezic/sweetd/machine"
)

func main() {
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
			machine.ToggleBuzzer(on)
			machine.ToggleMotor(on)
		case <-done:
			return
		}
	}
}