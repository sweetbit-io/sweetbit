package main

import (
	"flag"
	"time"
	"log"
	"os"
	"os/signal"
	"net/url"
	"github.com/gorilla/websocket"
	"encoding/json"
	"fmt"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot"
	"github.com/matryer/try"
)

var addr = flag.String("addr", "35cqbbBap7trDfb18YExSBjaN8rTQ8CmhL", "receiving Bitcoin address")
var device = flag.String("device", "/dev/tty.usbmodem1411", "interface to the usb dispenser")
var initialDispense = flag.Duration("dispense", 0, "initial dispensing duration")

type AddrSubMessage struct {
	Op   string `json:"op"`
	Addr string `json:"addr"`
}

type UtxMessage struct {
	Op   string `json:"op"`
	X    struct {
		Out []struct {
			Addr  string `json:"addr"`
			Value int32  `json:"value"`
		} `json:"out"`
	} `json:"x"`
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	log.Println("Hello!")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: "ws.blockchain.info", Path: "/inv"}
	log.Printf("Connecting to %s", u.String())

	var c *websocket.Conn

	err := try.Do(func(attempt int) (bool, error) {
		var err error
		c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Println("dial failed, retrying in 5s:", err)
			time.Sleep(5 * time.Second)
		}
		return attempt < 5, err
	})
	if err != nil {
		log.Fatal("dial failed 5 times:", err)
	}

	defer c.Close()

	done := make(chan struct{})
	txns := make(chan UtxMessage)

	go func() {
		fmt.Println("running inline go func")

		defer close(txns)
		for {
			_, message, err := c.ReadMessage()

			fmt.Println("received message")

			if err != nil {
				log.Println("read:", err)
				return
			}

			fmt.Println("read", message)

			var msg UtxMessage

			err = json.Unmarshal(message, &msg)
			if err != nil {
				fmt.Println("error:", err)
			}

			log.Println("incoming tx", msg)

			txns <- msg
		}
	}()

	fmt.Println("listen to", *addr)

	subOp := AddrSubMessage{Op:"addr_sub", Addr:*addr}

	b, err := json.Marshal(subOp)

	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("subscribing", b)
	os.Stdout.Write(b)

	err = c.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Println("write:", err)
		return
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
			tx := <-txns

			value := int32(0)

			for _, out := range tx.X.Out {
				if out.Addr == *addr {
					value += out.Value
				}
			}

			fmt.Println("dispensing for tx", tx)
			fmt.Println("value is", value)

			dispense := time.Duration(value / 10) * time.Millisecond

			fmt.Println("dispensing time", dispense)

			pin.On()
			time.Sleep(dispense)
			pin.Off()
		}
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{pin},
		work,
	)

	robot.Start()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
