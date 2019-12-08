package pos

import (
	"encoding/json"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/the-lightning-land/sweetd/lightning"
	"github.com/the-lightning-land/sweetd/nodeman"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var localhostOriginPattern = regexp.MustCompile(`^https?://localhost(:\d+)?$`)

type Dispenser interface {
	GetNodes() []nodeman.LightningNode
	GetNode(id string) nodeman.LightningNode
}

type Config struct {
	Logger    Logger
	Dispenser Dispenser
}

type Node struct {
	lightning.Node
	Id          string
	Name        string
	Description string
	Order       int
}

type Handler struct {
	http.Handler
	log       Logger
	dispenser Dispenser
}

func NewHandler(config *Config) *Handler {
	pos := &Handler{}

	if config.Logger != nil {
		pos.log = config.Logger
	} else {
		pos.log = noopLogger{}
	}

	pos.dispenser = config.Dispenser

	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	api.Use(pos.createLoggingMiddleware(pos.log.Infof))
	api.Use(pos.localhostMiddleware)
	api.Use(pos.availabilityMiddleware)
	api.Handle("/invoices/{rHash}/status", pos.handleStreamInvoiceStatus()).Methods(http.MethodGet, http.MethodOptions)
	api.Handle("/invoices/{rHash}", pos.handleGetInvoice()).Methods(http.MethodGet, http.MethodOptions)
	api.Handle("/invoices", pos.handleAddInvoice()).Methods(http.MethodPost, http.MethodOptions)
	api.Use(mux.CORSMethodMiddleware(api))

	box := packr.New("web", "./out")
	router.Use(pos.createLoggingMiddleware(pos.log.Debugf))
	router.PathPrefix("/").Handler(pos.handleStatic(box)).Methods(http.MethodGet)

	pos.Handler = router

	return pos
}

func (p *Handler) createLoggingMiddleware(log func(string, ...interface{})) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log("Accessing %v", r.RequestURI)
			next.ServeHTTP(w, r)
		})
	}
}

func (p *Handler) localhostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if localhostOriginPattern.MatchString(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Max-Age", "1")
		}
		next.ServeHTTP(w, r)
	})
}

func (p *Handler) getActiveNode() nodeman.LightningNode {
	nodes := p.dispenser.GetNodes()

	if len(nodes) > 0 {
		return nodes[0]
	}

	return nil
}

func (p *Handler) availabilityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p.getActiveNode() == nil {
			p.log.Errorf("PoS request failed due to unavailable node")
			p.jsonError(w, "No node is available at the moment", http.StatusServiceUnavailable)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (p *Handler) handleStatic(box *packr.Box) http.Handler {
	return http.FileServer(box)
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}

	if strings.Contains(origin[0], "http://localhost:3001") {
		return true
	}

	u, err := url.Parse(origin[0])
	if err != nil {
		return false
	}

	return strings.EqualFold(u.Host, r.Host)
}

func (p *Handler) handleStreamInvoiceStatus() http.HandlerFunc {
	upgrader := &websocket.Upgrader{
		CheckOrigin: checkOrigin,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rHash := vars["rHash"]

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			p.log.Errorf("Could not upgrade: %v", err)
			return
		}

		// read pump
		go func() {
			defer c.Close()

			c.SetReadLimit(512)
			c.SetReadDeadline(time.Now().Add(60 * time.Second))
			c.SetPongHandler(func(string) error {
				c.SetReadDeadline(time.Now().Add(60 * time.Second))
				return nil
			})

			for {
				_, _, err := c.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						p.log.Errorf("Unexpected websocket closure: %v", err)
					}
					break
				}
			}
		}()

		// write pump
		go func() {
			defer c.Close()

			ticker := time.NewTicker(54 * time.Second)
			defer ticker.Stop()

			client, err := p.getActiveNode().SubscribeInvoices()
			if err != nil {
				p.log.Errorf("Could not subscribe to invoices: %v", err)
				return
			}

			defer client.Cancel()

			for {
				select {
				case invoice, ok := <-client.Invoices:
					c.SetWriteDeadline(time.Now().Add(10 * time.Second))

					if !ok {
						c.WriteMessage(websocket.CloseMessage, []byte{})
						return
					}

					if invoice.RHash != rHash {
						continue
					}

					err := c.WriteJSON(&invoiceStatusMessage{
						Settled: invoice.Settled,
					})
					if err != nil {
						return
					}
				case <-ticker.C:
					c.SetWriteDeadline(time.Now().Add(10 * time.Second))
					if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
						return
					}
				}
			}
		}()
	}
}

func (p *Handler) handleGetInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rHash := vars["rHash"]

		invoice, err := p.getActiveNode().GetInvoice(rHash)
		if err != nil {
			p.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&invoiceMessage{
			Settled:        invoice.Settled,
			RHash:          invoice.RHash,
			PaymentRequest: invoice.PaymentRequest,
		})
		if err != nil {
			p.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (p *Handler) handleAddInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invoice, err := p.getActiveNode().AddInvoice(&lightning.InvoiceRequest{
		})
		if err != nil {
			p.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&invoiceMessage{
			Settled:        invoice.Settled,
			RHash:          invoice.RHash,
			PaymentRequest: invoice.PaymentRequest,
		})
		if err != nil {
			p.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type invoiceMessage struct {
	RHash          string `json:"r_hash"`
	PaymentRequest string `json:"payment_request"`
	Settled        bool   `json:"settled"`
}

type invoiceStatusMessage struct {
	Settled bool `json:"settled"`
}
