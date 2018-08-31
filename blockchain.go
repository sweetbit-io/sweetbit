package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/matryer/try"
	"log"
	"net/url"
	"time"
)

type AddrSubMessage struct {
	Op   string `json:"op"`
	Addr string `json:"addr"`
}

type UtxMessage struct {
	Op string `json:"op"`
	X  struct {
		Out []struct {
			Addr  string `json:"addr"`
			Value int    `json:"value"`
		} `json:"out"`
	} `json:"x"`
}

func listenForBlockchainTxns(bitcoinAddr string, transactions chan<- UtxMessage) {
	var conn *websocket.Conn

	u := url.URL{Scheme: "wss", Host: "ws.blockchain.info", Path: "/inv"}

	// Establish connection to candy subscription web socket endpoint
	err := try.Do(func(attempt int) (bool, error) {
		var err error

		conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Println("Dial failed. Retrying in 5s: ", err)
			time.Sleep(5 * time.Second)
		}
		return attempt < 12, err
	})

	if err != nil {
		log.Println("Dial failed 12 times: ", err)
	}

	// Close established web socket connection when done
	defer conn.Close()

	// Create subscribe payload
	payload, err := json.Marshal(AddrSubMessage{Op: "addr_sub", Addr: bitcoinAddr})
	if err != nil {
		log.Println("Marshal: ", err)
	}

	// Start listening for incoming paid invoices
	err = conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		log.Println("WriteMessage: ", err)
		return
	}

	for {
		// Read incoming message
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage: ", err)
			return
		}

		var msg UtxMessage

		// Parse message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("Unmarshal: ", err)
		}

		// Send incoming transaction to channel
		transactions <- msg
	}
}
