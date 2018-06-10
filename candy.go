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
	"github.com/go-errors/errors"
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

func establishConnection(candySubscriptionEndpoint string) (*websocket.Conn, error) {
	var conn *websocket.Conn

	// Establish connection to candy subscription web socket endpoint
	err := try.Do(func(attempt int) (bool, error) {
		var err error

		var header = http.Header{}
		header.Add("Sec-WebSocket-Protocol", "graphql-ws")

		conn, _, err = websocket.DefaultDialer.Dial(candySubscriptionEndpoint, header)
		if err != nil {
			log.Println("Dial failed. Retrying in 5s: ", err)
			time.Sleep(5 * time.Second)
		}
		return attempt < 5, err
	})

	if err != nil {
		log.Fatal("Dial failed 5 times: ", err)
		return nil, errors.New("Dial failed 5 times")
	}

	// Send connection initialization message
	err = conn.WriteMessage(websocket.TextMessage, []byte(connectionInitPayload))
	if err != nil {
		log.Fatal("Init WriteMessage: ", err)
		return nil, errors.New("Sending connection init message failed")
	}

	// Wait for ack
	_, _, err = conn.ReadMessage()
	if err != nil {
		log.Fatal("Ack ReadMessage:", err)
		return nil, errors.New("Receiving connection ack message failed")
	}

	return conn, nil
}

func subscribeToPaidInvoices(conn *websocket.Conn) error {
	// Create start payload template
	startPayloadTmpl := template.New("start")
	startPayloadTmpl, err := startPayloadTmpl.Parse(startPayloadTmplString)
	if err != nil {
		log.Fatal("Template Parse: ", err)
		return err
	}

	// Start listening for incoming paid invoices
	var payload bytes.Buffer
	err = startPayloadTmpl.Execute(&payload, StartPayloadData{
		Id: 1,
	})
	if err != nil {
		log.Fatal("Template Execute: ", err)
		return err
	}

	// Send subscription message
	err = conn.WriteMessage(websocket.TextMessage, payload.Bytes())
	if err != nil {
		log.Fatal("Subscription WriteMessage: ", err)
		return err
	}

	return nil
}

func establishConnectionAndSubscribeToPaidInvoices(candySubscriptionEndpoint string) (*websocket.Conn, error) {
	conn, err := establishConnection(candySubscriptionEndpoint)
	if err != nil {
		log.Fatal("Connection failed: ", err)
		return nil, err
	}

	log.Println("Connected")

	err = subscribeToPaidInvoices(conn)
	if err != nil {
		log.Fatal("Subscription: ", err)
		return conn, err
	}

	log.Println("Subscribed to invoices")

	return conn, nil
}

func listenForCandyPayments(candySubscriptionEndpoint string, paidInvoices chan<- Invoice) {
	var conn *websocket.Conn

	conn, err := establishConnectionAndSubscribeToPaidInvoices(candySubscriptionEndpoint)
	if err != nil {
		log.Fatal("Connection and subscription failed: ", err)
		return
	}

	defer conn.Close()

	for {
		// Read incoming message
		_, message, err := conn.ReadMessage()
		if websocket.IsCloseError(err, 1006) {
			conn, err = establishConnectionAndSubscribeToPaidInvoices(candySubscriptionEndpoint)
			if err != nil {
				log.Fatal("Connection and subscription failed: ", err)
				return
			}

			defer conn.Close()
		} else if err != nil {
			log.Fatal("Incoming ReadMessage: ", err)
			return
		} else {
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
}