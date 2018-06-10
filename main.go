package main

import (
	"flag"
	"time"
	"log"
	"fmt"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot"
)

var candyEndpoint = flag.String("candyendpoint", "", "where to listen for paid invoices")
var bitcoinAddress = flag.String("bitcoinaddress", "", "receiving Bitcoin address")
var device = flag.String("device", "/dev/tty.usbmodem1411", "interface to the usb dispenser")
var initialDispense = flag.Duration("dispense", 0, "initial dispensing duration")

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

	if *bitcoinAddress != "" {
		go listenForBlockchainTxns(*bitcoinAddress, transactions)
	}
	if *candyEndpoint != "" {
		go listenForCandyPayments(*candyEndpoint, invoices)
	}

	firmataAdaptor := firmata.NewAdaptor(*device)
	pin := gpio.NewDirectPinDriver(firmataAdaptor, "3")

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

				dispense := time.Duration(payment.Value / 10) * time.Millisecond

				log.Println("Dispensing for a duration of", dispense)

				pin.On()
				time.Sleep(dispense)
				pin.Off()
			case invoice := <- invoices:
				payment = Payment{Value:invoice.Value, Type:"bitcoin"}

				dispense := time.Duration(payment.Value / 10) * time.Millisecond

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