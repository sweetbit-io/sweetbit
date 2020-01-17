package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/the-lightning-land/sweetd/lightning"
	"github.com/the-lightning-land/sweetd/nodeman"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	postNodesTypeRemoteLnd = "remote-lnd"
	postNodesTypeLocal     = "local"
)

type postNodesRequest struct {
	Type string `json:"type"`
}

type postNodesRemoteLndRequest struct {
	Name     string `json:"name"`
	Uri      string `json:"uri"`
	Macaroon string `json:"macaroon"`
	Cert     string `json:"cert"`
}

type postNodesLocalRequest struct {
	Name string `json:"name"`
}

type postNodesRemoteLndResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Uri     string `json:"uri"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Status  string `json:"status"`
}

type postNodesLocalResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Uri     string `json:"uri"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Status  string `json:"status"`
}

type getNodesRemoteLndResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Uri     string `json:"uri"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Status  string `json:"status"`
}

type getNodesLocalLndResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Uri     string `json:"uri"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Status  string `json:"status"`
}

type getNodesResponse []interface{}

type patchNodeRequest struct {
	Op string `json:"op"`
}

type patchNodeRenameRequest struct {
	Name string `json:"name"`
}

type patchNodeEnableRequest struct {
	Enabled bool `json:"enabled"`
}

type patchNodeUnlockRequest struct {
	Password string `json:"password"`
}

type patchNodeInitRequest struct {
	Password string   `json:"password"`
	Mnemonic []string `json:"mnemonic"`
}

type nodeStatusResponse struct {
	Status string `json:"status"`
}

type postNodeSeedRequest struct {
}

type nodeSeedResponse struct {
	Mnemonic []string `json:"mnemonic"`
}

type postNodeConnectionRequest struct {
}

type nodeConnectionResponse struct {
	Uri      string `json:"uri"`
	Cert     string `json:"cert"`
	Macaroon string `json:"macaroon"`
}

func nodeStatusString(status lightning.Status) string {
	switch status {
	case lightning.StatusStopped:
		return "stopped"
	case lightning.StatusUninitialized:
		return "uninitialized"
	case lightning.StatusLocked:
		return "locked"
	case lightning.StatusStarted:
		return "started"
	case lightning.StatusFailed:
		return "failed"
	default:
		return ""
	}
}

func (a *Handler) postNodes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			a.jsonError(w, fmt.Sprintf("unable to read body: %v", err), http.StatusInternalServerError)
			return
		}

		req := postNodesRequest{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		switch req.Type {
		case postNodesTypeRemoteLnd:
			req := postNodesRemoteLndRequest{}
			err := json.Unmarshal(body, &req)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			macaroonBytes, err := base64.StdEncoding.DecodeString(req.Macaroon)
			if err != nil {
				a.jsonError(w, fmt.Sprintf("unable to decode macaroon: %v", err), http.StatusBadRequest)
				return
			}

			node, err := a.dispenser.AddNode(&nodeman.RemoteLndNodeConfig{
				Name:     req.Name,
				Uri:      req.Uri,
				Macaroon: macaroonBytes,
				Cert:     []byte(req.Cert),
			})
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			remoteLndNode := node.(*nodeman.RemoteLndNode)

			a.jsonResponse(w, &postNodesRemoteLndResponse{
				ID:      remoteLndNode.ID(),
				Type:    postNodesTypeRemoteLnd,
				Uri:     remoteLndNode.Uri,
				Name:    remoteLndNode.Name(),
				Enabled: remoteLndNode.Enabled(),
				Status:  nodeStatusString(remoteLndNode.Status()),
			}, http.StatusOK)
		case postNodesTypeLocal:
			req := postNodesLocalRequest{}
			err := json.Unmarshal(body, &req)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			node, err := a.dispenser.AddNode(&nodeman.LocalNodeConfig{
				Name: req.Name,
			})
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			localNode := node.(*nodeman.LocalNode)

			a.jsonResponse(w, &postNodesLocalResponse{
				ID:      localNode.ID(),
				Type:    postNodesTypeLocal,
				Uri:     localNode.Uri(),
				Name:    localNode.Name(),
				Enabled: localNode.Enabled(),
				Status:  nodeStatusString(localNode.Status()),
			}, http.StatusOK)
		default:
			a.jsonError(w, fmt.Sprintf("unknown type \"%s\"", req.Type), http.StatusBadRequest)
		}
	}
}

func (a *Handler) getNodes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results := getNodesResponse{}

		nodes := a.dispenser.GetNodes()

		for _, node := range nodes {
			switch node := node.(type) {
			case *nodeman.RemoteLndNode:
				results = append(results, &getNodesRemoteLndResponse{
					ID:      node.ID(),
					Type:    postNodesTypeRemoteLnd,
					Uri:     node.Uri,
					Name:    node.Name(),
					Enabled: node.Enabled(),
					Status:  nodeStatusString(node.Status()),
				})
			case *nodeman.LocalNode:
				results = append(results, &getNodesLocalLndResponse{
					ID:      node.ID(),
					Type:    postNodesTypeLocal,
					Uri:     node.Uri(),
					Name:    node.Name(),
					Enabled: node.Enabled(),
					Status:  nodeStatusString(node.Status()),
				})
			default:
				a.log.Warnf("got unknown type of node %T", node)
			}
		}

		a.jsonResponse(w, &results, http.StatusOK)
	}
}

func (a *Handler) deleteNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		err := a.dispenser.RemoveNode(id)
		if err != nil {
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		a.emptyResponse(w, http.StatusNoContent)
	}
}

func (a *Handler) patchNode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			a.jsonError(w, fmt.Sprintf("unable to read body: %v", err), http.StatusInternalServerError)
			return
		}

		req := patchNodeRequest{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		node := a.dispenser.GetNode(id)

		switch req.Op {
		case "rename":
			req := patchNodeRenameRequest{}
			err := json.Unmarshal(body, &req)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = a.dispenser.RenameNode(id, req.Name)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case "enable":
			req := patchNodeEnableRequest{}
			err := json.Unmarshal(body, &req)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if req.Enabled {
				err = a.dispenser.EnableNode(id)
				if err != nil {
					a.jsonError(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				err = a.dispenser.DisableNode(id)
				if err != nil {
					a.jsonError(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		case "unlock":
			req := patchNodeUnlockRequest{}
			err := json.Unmarshal(body, &req)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = node.Unlock(req.Password)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case "init":
			req := patchNodeInitRequest{}
			err := json.Unmarshal(body, &req)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = node.Init(req.Password, req.Mnemonic)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			a.jsonError(w, "Can only rename, enable, disable, init and unlock node.", http.StatusBadRequest)
			return
		}

		switch node := node.(type) {
		case *nodeman.RemoteLndNode:
			a.jsonResponse(w, &getNodesRemoteLndResponse{
				ID:      node.ID(),
				Type:    postNodesTypeRemoteLnd,
				Uri:     node.Uri,
				Name:    node.Name(),
				Enabled: node.Enabled(),
				Status:  nodeStatusString(node.Status()),
			}, http.StatusOK)
			return
		case *nodeman.LocalNode:
			a.jsonResponse(w, &getNodesLocalLndResponse{
				ID:      node.ID(),
				Type:    postNodesTypeLocal,
				Uri:     node.Uri(),
				Name:    node.Name(),
				Enabled: node.Enabled(),
				Status:  nodeStatusString(node.Status()),
			}, http.StatusOK)
			return
		default:
			a.jsonError(w, fmt.Sprintf("unknown node type %T", node), http.StatusBadRequest)
			return
		}
	}
}

func (a *Handler) handleGetNodeStatusEvents() http.HandlerFunc {
	upgrader := &websocket.Upgrader{
		CheckOrigin: checkOrigin,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		node := a.dispenser.GetNode(id)
		if node == nil {
			a.jsonError(w, fmt.Sprintf("No node with id %s found", id), http.StatusNotFound)
			return
		}

		client := node.SubscribeStatus()

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			client.Cancel()
			a.log.Errorf("unable to upgrade: %v", err)
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
				case status, ok := <-client.Status:
					c.SetWriteDeadline(time.Now().Add(10 * time.Second))

					if !ok {
						c.WriteMessage(websocket.CloseMessage, []byte{})
						return
					}

					err := c.WriteJSON(&nodeStatusResponse{
						Status: nodeStatusString(status),
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

func (a *Handler) handlePostNodeSeed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		req := postNodeSeedRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		node := a.dispenser.GetNode(id)
		if node == nil {
			a.jsonError(w, fmt.Sprintf("No node with id %s found", id), http.StatusNotFound)
			return
		}

		mnemonic, err := node.GenerateSeed()
		if err != nil {
			a.jsonError(w, fmt.Sprintf("Unable to generate seed: %v", err), http.StatusInternalServerError)
			return
		}

		a.jsonResponse(w, &nodeSeedResponse{
			Mnemonic: mnemonic,
		}, http.StatusOK)
	}
}

func (a *Handler) handlePostNodeConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		req := postNodeConnectionRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		node := a.dispenser.GetNode(id)
		if node == nil {
			a.jsonError(w, fmt.Sprintf("No node with id %s found", id), http.StatusNotFound)
			return
		}

		// TODO support all types
		localNode, ok := node.(*nodeman.LocalNode)
		if !ok {
			a.jsonError(w, fmt.Sprintf("No connections available for node type %T", node), http.StatusNotFound)
			return
		}

		// TODO create a new macaroon here

		a.jsonResponse(w, &nodeConnectionResponse{
			Uri:      localNode.Uri(),
			Cert:     localNode.Cert(),
			Macaroon: localNode.AdminMacaroon(),
		}, http.StatusOK)
	}
}
