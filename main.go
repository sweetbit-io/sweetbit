package main

import (
	"flag"
	"log"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/drivers/gpio"
	"fmt"
	"time"
	"gobot.io/x/gobot"
)

var candyEndpoint = flag.String("lightning.subscription", "", "subscription endpoint to paid lightning invoices")
var bitcoinAddress = flag.String("bitcoin.address", "", "receiving Bitcoin address")
var device = flag.String("device.path", "/dev/tty.usbmodem1411", "path to the USB device")
var devicePin = flag.String("device.pin", "3", "dispensing GPIO pin on device")
var initialDispense = flag.Duration("debug.dispense", 0, "dispensing duration on startup")

type Payment struct {
	Type  string
	Value int
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	done := make(chan struct{})
	transactions := make(chan UtxMessage)
	invoices := make(chan Invoice)
	stop:= make(chan bool)

	if *bitcoinAddress != "" {
		go listenForBlockchainTxns(*bitcoinAddress, transactions)
	}
	if *candyEndpoint != "" {
		go listenForCandyPayments(*candyEndpoint, invoices, stop)
	}

	firmataAdaptor := firmata.NewAdaptor(*device)
	pin := gpio.NewDirectPinDriver(firmataAdaptor, *devicePin)

	work := func() {
		if *initialDispense > 0 {
			fmt.Println("Initial dispensing is on")
			fmt.Println("Dispensing for", *initialDispense)

			pin.On()
			time.Sleep(*initialDispense)
			pin.Off()
		}

		for {
			var payment Payment

			select {
			case tx := <- transactions:
				value := 0

				for _, out := range tx.X.Out {
					if out.Addr == *bitcoinAddress {
						value += out.Value
					}
				}
				payment = Payment{Value:value, Type:"lightning"}

				dispense := time.Duration(payment.Value / 2) * time.Millisecond

				log.Println("Dispensing for a duration of", dispense)

				pin.On()
				time.Sleep(dispense)
				pin.Off()
			case invoice := <- invoices:
				payment = Payment{Value:invoice.Value, Type:"bitcoin"}

				dispense := time.Duration(payment.Value / 2) * time.Millisecond

				log.Println("Dispensing for a duration of", dispense)

				pin.On()
				time.Sleep(dispense)
				pin.Off()
			}
		}
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{pin},
		work,
	)

	robot.Start()

	<-done
}