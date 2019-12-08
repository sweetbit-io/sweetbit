package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/the-lightning-land/sweetd/nodeman"
	"io/ioutil"
	"net/http"
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
}

type postNodesLocalResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type getNodesRemoteLndResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Uri     string `json:"uri"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type getNodesLocalLndResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type getNodesResponse []interface{}

type patchNodeRequest struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
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
				Type:    postNodesTypeRemoteLnd,
				Name:    localNode.Name(),
				Enabled: localNode.Enabled(),
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
				})
			case *nodeman.LocalNode:
				results = append(results, &getNodesLocalLndResponse{
					ID:      node.ID(),
					Type:    postNodesTypeLocal,
					Name:    node.Name(),
					Enabled: node.Enabled(),
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

		req := patchNodeRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			a.jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		node := a.dispenser.GetNode(id)

		if req.Enabled != node.Enabled() {
			var err error
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
		} else if req.Name != node.Name() {
			err = a.dispenser.RenameNode(id, req.Name)
			if err != nil {
				a.jsonError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			a.jsonError(w, "Can only enable, disable and rename node.", http.StatusBadRequest)
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
			}, http.StatusOK)
			return
		case *nodeman.LocalNode:
			a.jsonResponse(w, &getNodesLocalLndResponse{
				ID:      node.ID(),
				Type:    postNodesTypeLocal,
				Name:    node.Name(),
				Enabled: node.Enabled(),
			}, http.StatusOK)
			return
		default:
			a.jsonError(w, fmt.Sprintf("unknown node type %T", node), http.StatusBadRequest)
			return
		}
	}
}
