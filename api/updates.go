package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/the-lightning-land/sweetd/updater"
	"net/http"
	"time"
)

type postUpdateRequest struct {
	Url string `json:"url"`
}

type patchUpdateRequest struct {
	State string `json:"state"`
}

type updateResponse struct {
	Id           string    `json:"id"`
	Started      time.Time `json:"started"`
	Url          string    `json:"url"`
	State        string    `json:"state"`
	Progress     uint8     `json:"progress"`
	ShouldReboot bool      `json:"reboot"`
	ShouldCommit bool      `json:"commit"`
}

func (a *Handler) handlePostUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := postUpdateRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		update, err := a.dispenser.StartUpdate(req.Url)
		if err != nil {
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		a.jsonResponse(w, &updateResponse{
			Id:           update.Id,
			Started:      update.Started,
			Url:          update.Url,
			State:        update.State,
			Progress:     update.Progress,
			ShouldReboot: update.ShouldReboot,
			ShouldCommit: update.ShouldCommit,
		}, http.StatusOK)
	}
}

func (a *Handler) handleGetUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		update, err := a.dispenser.GetUpdate(id)
		if err != nil {
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if update == nil {
			a.jsonError(w, fmt.Sprintf("No update with id %s found.", id), http.StatusNotFound)
			return
		}

		a.jsonResponse(w, &updateResponse{
			Id:           update.Id,
			Started:      update.Started,
			Url:          update.Url,
			State:        update.State,
			Progress:     update.Progress,
			ShouldReboot: update.ShouldReboot,
			ShouldCommit: update.ShouldCommit,
		}, http.StatusOK)
	}
}

func (a *Handler) handlePatchUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		req := patchUpdateRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var update *updater.Update

		if req.State == updater.StateCancelled {
			update, err = a.dispenser.CancelUpdate(id)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if req.State == updater.StateCompleted {
			update, err = a.dispenser.CommitUpdate(id)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if req.State == updater.StateRejected {
			update, err = a.dispenser.RejectUpdate(id)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			a.jsonError(w, fmt.Sprint("Can only cancel, commit or reject an update."), http.StatusBadRequest)
			return
		}

		a.jsonResponse(w, &updateResponse{
			Id:           update.Id,
			Started:      update.Started,
			Url:          update.Url,
			State:        update.State,
			Progress:     update.Progress,
			ShouldReboot: update.ShouldReboot,
			ShouldCommit: update.ShouldCommit,
		}, http.StatusOK)
	}
}

func (a *Handler) handleGetUpdateEvents() http.HandlerFunc {
	upgrader := &websocket.Upgrader{}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		client, err := a.dispenser.SubscribeUpdate(id)
		if err != nil {
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if client == nil {
			a.jsonError(w, fmt.Sprintf("No update with id %s found", id), http.StatusNotFound)
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			client.Cancel()
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
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
						a.log.Errorf("unexpected websocket closure: %v", err)
					}
					break
				}
			}
		}()

		// write pump
		go func() {
			defer c.Close()
			defer client.Cancel()

			ticker := time.NewTicker(54 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case update, ok := <-client.Update:
					c.SetWriteDeadline(time.Now().Add(10 * time.Second))

					if !ok {
						c.WriteMessage(websocket.CloseMessage, []byte{})
						return
					}

					err := c.WriteJSON(&updateResponse{
						Id:           update.Id,
						Started:      update.Started,
						Url:          update.Url,
						State:        update.State,
						Progress:     update.Progress,
						ShouldReboot: update.ShouldReboot,
						ShouldCommit: update.ShouldCommit,
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
