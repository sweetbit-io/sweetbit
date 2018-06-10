package main

import (
	"text/template"
	"log"
	"github.com/gorilla/websocket"
	"github.com/matryer/try"
	"time"
	"encoding/json"
	"net/http"
	"bytes"
)

type InvoiceMessage struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Payload struct {
		Data struct {
			InvoicesPaid Invoice `json:"invoicesPaid"`
		} `json:"data"`
	} `json:"payload"`
}

type Invoice struct {
	RHash   string `json:"r_hash"`
	Settled bool   `json:"settled"`
	Value   int    `json:"value"`
}

const connectionInitPayload = `{ "type": "connection_init", "payload": {} }`

type StartPayloadData struct {
	Id int
}

const startPayloadTmplString = `{
	"id": "{{.Id}}",
	"type": "start",
	"payload": {
		"variables": {},
		"extensions": {},
		"operationName": null,
		"query": "subscription { invoicesPaid { r_hash, value, settled } }"
	}
}`

func listenForCandyPayments(candySubscriptionEndpoint string, paidInvoices chan<- Invoice) {
	var conn *websocket.Conn

	// Establish connection to candy subscription web socket endpoint
	err := try.Do(func(attempt int) (bool, error) {
		var err error

		var header = http.Header{}
		// header.Add("Sec-WebSocket-Extensions", "permessage-deflate; client_max_window_bits")
		// header.Add("Sec-WebSocket-Key", "kBzKaCDYRKnLaTGDzsFeXg==")
		header.Add("Sec-WebSocket-Protocol", "graphql-ws")
		// header.Add("Sec-WebSocket-Version", "13")

		conn, _, err = websocket.DefaultDialer.Dial(candySubscriptionEndpoint, header)
		if err != nil {
			log.Fatal("Dial failed. Retrying in 5s: ", err)
			time.Sleep(5 * time.Second)
		}
		return attempt < 5, err
	})

	if err != nil {
		log.Fatal("Dial failed 5 times: ", err)
	}

	// Close established web socket connection when done
	defer conn.Close()

	// Send connection initialization message
	err = conn.WriteMessage(websocket.TextMessage, []byte(connectionInitPayload))
	if err != nil {
		log.Fatal("Init WriteMessage: ", err)
		return
	}

	// Wait for ack
	_, _, err = conn.ReadMessage()
	if err != nil {
		log.Fatal("Ack ReadMessage:", err)
		return
	}

	// Create start payload template
	startPayloadTmpl := template.New("start")
	startPayloadTmpl, err = startPayloadTmpl.Parse(startPayloadTmplString)
	if err != nil {
		log.Fatal("Template Parse: ", err)
		return
	}

	// Create a connection writer
	// writer, err := conn.NextWriter(websocket.TextMessage)
	// if err != nil {
	// 	log.Fatal("NextWriter: ", err)
	// 	return
	// }

	// Start listening for incoming paid invoices
	var payload bytes.Buffer
	err = startPayloadTmpl.Execute(&payload, StartPayloadData{
		Id: 1,
	})
	if err != nil {
		log.Fatal("Template Execute: ", err)
		return
	}

	// Send subscription message
	err = conn.WriteMessage(websocket.TextMessage, payload.Bytes())
	if err != nil {
		log.Fatal("Subscription WriteMessage: ", err)
		return
	}

	log.Println("Subscribed to invoices")

	for {
		// Read incoming message
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Incoming ReadMessage: ", err)
			return
		}

		log.Println("Got a message", string(message))

		var msg InvoiceMessage

		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Fatal("Incoming Unmarshal: ", err)
		}

		log.Println("Got paid invoice", msg)

		paidInvoices <- msg.Payload.Data.InvoicesPaid
	}
}