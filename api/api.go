package api

import (
	"github.com/gorilla/mux"
	"github.com/the-lightning-land/sweetd/network"
	"github.com/the-lightning-land/sweetd/nodeman"
	"github.com/the-lightning-land/sweetd/state"
	"github.com/the-lightning-land/sweetd/sweetdb"
	"github.com/the-lightning-land/sweetd/updater"
	"net/http"
	"regexp"
)

var localhostOriginPattern = regexp.MustCompile(`^https?://localhost(:\d+)?$`)

type Config struct {
	Dispenser Dispenser
	Log       Logger
}

type Handler struct {
	http.Handler
	dispenser Dispenser
	log       Logger
}

func NewHandler(config *Config) http.Handler {
	api := &Handler{
		dispenser: config.Dispenser,
	}

	if config.Log != nil {
		api.log = config.Log
	} else {
		api.log = noopLogger{}
	}

	router := mux.NewRouter()
	router.Use(api.loggingMiddleware)
	router.Use(api.localhostMiddleware)

	router.Handle("/dispenser", api.handleGetDispenser()).Methods(http.MethodGet, http.MethodOptions)
	router.Handle("/dispenser", api.handlePatchDispenser()).Methods(http.MethodPatch)
	router.Handle("/dispenser/events", api.handleGetDispenser()).Methods(http.MethodGet, http.MethodOptions)

	router.Handle("/updates", api.handlePostUpdate()).Methods(http.MethodPost, http.MethodOptions)
	router.Handle("/updates/{id}", api.handleGetUpdate()).Methods(http.MethodGet, http.MethodOptions)
	router.Handle("/updates/{id}", api.handlePatchUpdate()).Methods(http.MethodPatch)
	router.Handle("/updates/{id}/events", api.handleGetUpdateEvents()).Methods(http.MethodGet, http.MethodOptions)

	router.Handle("/nodes", api.noContent()).Methods(http.MethodOptions)
	router.Handle("/nodes", api.getNodes()).Methods(http.MethodGet)
	router.Handle("/nodes", api.postNodes()).Methods(http.MethodPost)
	router.Handle("/nodes/{id}", api.noContent()).Methods(http.MethodOptions)
	router.Handle("/nodes/{id}", api.getNodes()).Methods(http.MethodGet)
	router.Handle("/nodes/{id}", api.patchNode()).Methods(http.MethodPatch)
	router.Handle("/nodes/{id}", api.deleteNode()).Methods(http.MethodDelete)

	router.Handle("/networks", api.noContent()).Methods(http.MethodOptions)
	router.Handle("/networks", api.handlePostUpdate()).Methods(http.MethodPost)
	router.Handle("/networks/{id}", api.noContent()).Methods(http.MethodOptions)
	router.Handle("/networks/{id}", api.handlePostUpdate()).Methods(http.MethodPatch)
	router.Handle("/networks/events", api.noContent()).Methods(http.MethodOptions)
	router.Handle("/networks/events", api.handlePostUpdate()).Methods(http.MethodGet)

	router.Use(mux.CORSMethodMiddleware(router))

	api.Handler = router

	return api
}

func (a *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.log.Infof("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (a *Handler) noContent() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		w.WriteHeader(http.StatusNoContent)
	})
}

func (a *Handler) localhostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if localhostOriginPattern.MatchString(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Max-Age", "1")
		}
		next.ServeHTTP(w, r)
	})
}

type LightningNode interface {
	ID() string
	Name() string
	Enabled() bool
}

type Dispenser interface {
	GetNodes() []nodeman.LightningNode
	GetNode(id string) nodeman.LightningNode
	AddNode(config nodeman.NodeConfig) (nodeman.LightningNode, error)
	RemoveNode(id string) error
	EnableNode(id string) error
	DisableNode(id string) error
	RenameNode(id string, name string) error
	GetApiOnionID() string
	GetPosOnionID() string
	ToggleDispense(on bool)
	SetWifiConnection(connection sweetdb.Wifi) error
	GetState() state.State
	GetName() string
	ShouldDispenseOnTouch() bool
	ShouldBuzzOnDispense() bool
	SetName(name string) error
	SetDispenseOnTouch(dispenseOnTouch bool) error
	SetBuzzOnDispense(buzzOnDispense bool) error
	ConnectToWifi(connection network.Connection) error
	Reboot() error
	ShutDown() error
	Stop()
	// SubscribeDispenses() *dispenser.DispenseClient
	StartUpdate(url string) (*updater.Update, error)
	GetUpdate(id string) (*updater.Update, error)
	GetCurrentUpdate() (*updater.Update, error)
	CancelUpdate(id string) (*updater.Update, error)
	SubscribeUpdate(id string) (*updater.UpdateClient, error)
	CommitUpdate(id string) (*updater.Update, error)
	RejectUpdate(id string) (*updater.Update, error)
}
