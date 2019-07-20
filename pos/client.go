package pos

import (
  "github.com/gorilla/websocket"
  "log"
)

type client struct {
  conn               *websocket.Conn
  subscribedInvoices []string
}

func (c *client) processIncoming() {
  for {
    _, message, err := c.conn.ReadMessage()
    if err != nil {
      log.Println("read:", err)
      break
    }

    //err := json.Unmarshal(message, &lightningNode)
    //if err != nil {
    //  return errors.Errorf("Could not unmarshal data: %v", err)
    //}

    log.Printf("recv: %s", message)
  }
}

func (c *client) send() error {
  //payload, err := json.Marshal(&invoiceMessage{
  //  settled: true,
  //})
  //if err != nil {
  //  return errors.Errorf("Could not marshall %v", 0)
  //}
  //
  //err = c.conn.WriteMessage(websocket.TextMessage, payload)
  //if err != nil {
  //  return errors.Errorf("Could not write message to client: %v", payload)
  //}

  return nil
}
