package main

import (
	"os"
	"os/signal"
	"log"
)

func main() {
	machine := NewMachine()

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

	<-done
}