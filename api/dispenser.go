package api

import (
	"encoding/json"
	"fmt"
	"github.com/the-lightning-land/sweetd/state"
	"net/http"
)

type dispenserUpdateResponse struct {
	Id string `json:"id"`
}

type dispenserResponse struct {
	Name            string                   `json:"name"`
	Api             string                   `json:"api"`
	Pos             string                   `json:"pos"`
	Version         string                   `json:"version"`
	State           string                   `json:"state"`
	DispenseOnTouch bool                     `json:"dispenseOnTouch"`
	Update          *dispenserUpdateResponse `json:"update"`
}

type patchDispenserOp struct {
	Op    string      `json:"op"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type patchDispenserRequest []patchDispenserOp

func (a *Handler) getDispenser() *dispenserResponse {
	var currentUpdateRes *dispenserUpdateResponse
	currentUpdate, err := a.dispenser.GetCurrentUpdate()
	if err != nil {
		a.log.Errorf("unable to get current update: %v", err)
	}

	if currentUpdate != nil {
		currentUpdateRes = &dispenserUpdateResponse{
			Id: currentUpdate.Id,
		}
	}

	return &dispenserResponse{
		Name:            a.dispenser.GetName(),
		Api:             a.dispenser.GetApiOnionID(),
		Pos:             a.dispenser.GetPosOnionID(),
		State:           state.String(a.dispenser.GetState()),
		DispenseOnTouch: a.dispenser.ShouldDispenseOnTouch(),
		Update:          currentUpdateRes,
	}
}

func (a *Handler) handleGetDispenser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := a.getDispenser()
		a.jsonResponse(w, res, http.StatusOK)
	}
}

func (a *Handler) handlePatchDispenser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := patchDispenserRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := a.getDispenser()

		shutdownChan := make(chan struct{})

		for _, op := range req {
			if op.Op == "set" {
				if op.Name == "dispenseOnTouch" {
					if value, ok := op.Value.(bool); ok {
						err := a.dispenser.SetDispenseOnTouch(value)
						if err != nil {
							a.jsonError(w, "Could not set dispense on touch", http.StatusInternalServerError)
							return
						}

						res.DispenseOnTouch = value
					} else {
						a.jsonError(w, fmt.Sprintf("%s value not a boolean, but %T", op.Name, op.Value), http.StatusBadRequest)
						return
					}
				} else if op.Name == "name" {
					if value, ok := op.Value.(string); ok {
						err := a.dispenser.SetName(value)
						if err != nil {
							a.jsonError(w, "Could not set name", http.StatusInternalServerError)
							return
						}

						res.Name = value
					} else {
						a.jsonError(w, fmt.Sprintf("%s value not a string, but %T", op.Name, op.Value), http.StatusInternalServerError)
						return
					}
				} else {
					a.jsonError(w, fmt.Sprintf("unknown field %s", op.Name), http.StatusBadRequest)
					return
				}
			} else if op.Op == "reboot" {
				res.State = state.String(state.StateStopping)

				go func() {
					<-shutdownChan
					err := a.dispenser.Reboot()
					if err != nil {
						a.log.Errorf("unable to reboot: %v", err)
						return
					}
				}()
			} else if op.Op == "shutdown" {
				res.State = state.String(state.StateStopping)

				go func() {
					<-shutdownChan
					err := a.dispenser.ShutDown()
					if err != nil {
						a.log.Errorf("unable to shutdown: %v", err)
						return
					}
				}()
			} else {
				a.jsonError(w, fmt.Sprintf("unknown op %s", op.Op), http.StatusBadRequest)
				return
			}
		}

		a.jsonResponse(w, res, http.StatusOK)

		close(shutdownChan)
	}
}
