package main

import (
	"flag"
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
	"time"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type Payment struct {
	Type  string
	Value int
}

func main() {
	var candyEndpoints arrayFlags
	flag.Var(&candyEndpoints, "lightning.subscription", "subscription endpoint to paid lightning invoices")
	var bitcoinAddress = flag.String("bitcoin.address", "", "receiving Bitcoin address")
	var initialDispense = flag.Duration("debug.dispense", 0, "dispensing duration on startup")

	flag.Parse()
	log.SetFlags(0)

	done := make(chan struct{})
	transactions := make(chan UtxMessage)
	invoices := make(chan Invoice)
	stop := make(chan bool)

	if *bitcoinAddress != "" {
		go listenForBlockchainTxns(*bitcoinAddress, transactions)
	}

	for _, endpoint := range candyEndpoints {
		go listenForCandyPayments(endpoint, invoices, stop)
	}

	r := raspi.NewAdaptor()
	motorPin := gpio.NewDirectPinDriver(r, "2")
	vibratorPin := gpio.NewDirectPinDriver(r, "0")
	touchSensor := gpio.NewButtonDriver(r, "7")

	work := func() {
		button.On(gpio.ButtonPush, func(data interface{}) {
			fmt.Println("button pressed")
		})

		button.On(gpio.ButtonRelease, func(data interface{}) {
			fmt.Println("button released")
		})

		if *initialDispense > 0 {
			fmt.Println("Initial dispensing is on")
			fmt.Println("Dispensing for", *initialDispense)

			motorPin.On()
			time.Sleep(*initialDispense)
			motorPin.Off()
		}

		for {
			var payment Payment

			select {
			case tx := <-transactions:
				value := 0

				for _, out := range tx.X.Out {
					if out.Addr == *bitcoinAddress {
						value += out.Value
					}
				}
				payment = Payment{Value: value, Type: "lightning"}

				dispense := time.Duration(payment.Value/2) * time.Millisecond

				log.Println("Dispensing for a duration of", dispense)

				motorPin.On()
				time.Sleep(dispense)
				motorPin.Off()
			case invoice := <-invoices:
				payment = Payment{Value: invoice.Value, Type: "bitcoin"}

				dispense := time.Duration(payment.Value/2) * time.Millisecond

				log.Println("Dispensing for a duration of", dispense)

				motorPin.On()
				time.Sleep(dispense)
				motorPin.Off()
			}
		}
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{r},
		[]gobot.Device{motorPin, vibratorPin, touchSensor},
		work,
	)

	robot.Start()

	<-done
}
