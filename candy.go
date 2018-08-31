package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/gorilla/websocket"
	"github.com/matryer/try"
	"log"
	"net/http"
	"text/template"
	"time"
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
		return attempt < 12, err
	})

	if err != nil {
		log.Println("Dial failed 5 times: ", err)
		return nil, errors.New("Dial failed 12 times")
	}

	// Send connection initialization message
	err = conn.WriteMessage(websocket.TextMessage, []byte(connectionInitPayload))
	if err != nil {
		log.Println("Init WriteMessage: ", err)
		return nil, errors.New("Sending connection init message failed")
	}

	// Wait for ack
	_, _, err = conn.ReadMessage()
	if err != nil {
		log.Println("Ack ReadMessage:", err)
		return nil, errors.New("Receiving connection ack message failed")
	}

	return conn, nil
}

func subscribeToPaidInvoices(conn *websocket.Conn) error {
	// Create start payload template
	startPayloadTmpl := template.New("start")
	startPayloadTmpl, err := startPayloadTmpl.Parse(startPayloadTmplString)
	if err != nil {
		log.Println("Template Parse: ", err)
		return err
	}

	// Start listening for incoming paid invoices
	var payload bytes.Buffer
	err = startPayloadTmpl.Execute(&payload, StartPayloadData{
		Id: 1,
	})
	if err != nil {
		log.Println("Template Execute: ", err)
		return err
	}

	// Send subscription message
	err = conn.WriteMessage(websocket.TextMessage, payload.Bytes())
	if err != nil {
		log.Println("Subscription WriteMessage: ", err)
		return err
	}

	return nil
}

func listenForCandyPayments(candySubscriptionEndpoint string, paidInvoices chan<- Invoice, stop <-chan bool) {
	error := make(chan bool)

	client := &Client{error: error, paidInvoices: paidInvoices}
	go client.listen(candySubscriptionEndpoint, stop)

	for {
		select {
		case <-error:
			retry := 2 * time.Second
			log.Println("Connection closed. Establishing a new one in", retry)
			time.Sleep(retry)

			client := &Client{error: error, paidInvoices: paidInvoices}
			go client.listen(candySubscriptionEndpoint, stop)
		case <-stop:
			return
		}
	}
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	paidInvoices chan<- Invoice
	error        chan bool
}

func (c *Client) listen(candySubscriptionEndpoint string, stop <-chan bool) {
	var conn *websocket.Conn
	var err error

	if conn, err = establishConnection(candySubscriptionEndpoint); err != nil {
		log.Println("Connection failed: ", err)
		c.error <- true
		return
	}

	defer conn.Close()

	if err = subscribeToPaidInvoices(conn); err != nil {
		log.Println("Subscription to paid invoices failed: ", err)
		c.error <- true
		return
	}

	log.Println("Subscription established.")

	go c.read(conn)
	go c.write(conn)

	<-stop
}

func (c *Client) read(conn *websocket.Conn) {
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read Message: ", err)
			c.error <- true
			return
		}

		log.Println("Got a message", string(message))

		var msg InvoiceMessage

		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("Incoming Unmarshal: ", err)
		}

		log.Println("Got paid invoice", msg)

		c.paidInvoices <- msg.Payload.Data.InvoicesPaid
	}
}

func (c *Client) write(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)

	defer ticker.Stop()

	for {
		<-ticker.C
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			c.error <- true
			return
		}
	}
}
