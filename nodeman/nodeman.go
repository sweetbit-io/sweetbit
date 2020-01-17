package nodeman

import (
	"crypto/x509"
	"github.com/cretz/bine/tor"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/the-lightning-land/sweetd/lightning"
	"github.com/the-lightning-land/sweetd/onion"
	"github.com/the-lightning-land/sweetd/sweetdb"
	"path/filepath"
)

type Nodeman struct {
	// nodes that are set up for generating Lightning invoices and
	// accepting payments
	nodes []LightningNode

	// nodesDataDir directory where all node data is saved
	nodesDataDir string

	// db where node connection data is persisted
	db *sweetdb.DB

	// logCreator
	logCreator LogCreator

	// log
	log Logger

	// tor instance for nodes to expose services through
	tor *tor.Tor
}

type Config struct {
	// NodesDataDir directory where all node data is saved
	NodesDataDir string

	// DB where node connection data is persisted
	DB *sweetdb.DB

	// Tor instance for nodes to expose services through
	Tor *tor.Tor

	// LogCreator
	LogCreator LogCreator
}

type LogCreator func(node string) Logger

func New(config *Config) *Nodeman {
	nodeman := &Nodeman{
		nodes:        nil,
		nodesDataDir: config.NodesDataDir,
		db:           config.DB,
		tor:          config.Tor,
		logCreator:   config.LogCreator,
	}

	if config.LogCreator != nil {
		nodeman.log = config.LogCreator("")
	} else {
		nodeman.log = noopLogger{}
	}

	return nodeman
}

func (n *Nodeman) Load() {
	nodes, err := n.db.GetNodes()
	if err != nil {
		n.log.Errorf("unable to get nodes: %v", err)
	}

	for _, node := range nodes {
		switch node := node.(type) {
		case *sweetdb.RemoteLndNode:
			lndNode, err := lightning.NewLndNode(&lightning.LndNodeConfig{
				Uri:           node.Url,
				CertBytes:     node.Cert,
				MacaroonBytes: node.Macaroon,
				Logger:        n.log,
			})
			if err != nil {
				n.log.Errorf("unable to create node: %v", err)
				continue
			}

			n.nodes = append(n.nodes, &RemoteLndNode{
				LndNode: lndNode,
				id:      node.Id,
				name:    node.Name,
				enabled: node.Enabled,
				Uri:     node.Url,
			})
		case *sweetdb.LocalNode:
			key, err := x509.ParsePKCS1PrivateKey(node.OnionKey)
			if err != nil {
				n.log.Errorf("unable to parse onion key: %v", err)
				continue
			}

			onionSvc := onion.NewService(&onion.ServiceConfig{
				Tor:    n.tor,
				Logger: n.log,
				Port:   8080,
				Key:    key,
			})

			localNode, err := lightning.NewLocalNode(&lightning.LocalNodeConfig{
				DataDir:  filepath.Join(n.nodesDataDir, node.Id),
				Logger:   n.logCreator(node.Id),
				OnionSvc: onionSvc,
			})
			if err != nil {
				n.log.Errorf("unable to create node: %v", err)
				continue
			}

			n.nodes = append(n.nodes, &LocalNode{
				LocalNode: localNode,
				id:        node.Id,
				name:      node.Name,
				enabled:   node.Enabled,
			})
		default:
			n.log.Errorf("unknown node type %T", node)
		}
	}
}

func (n *Nodeman) GetNodes() []LightningNode {
	return n.nodes
}

func (n *Nodeman) GetNode(id string) LightningNode {
	for _, node := range n.nodes {
		if node.ID() == id {
			return node
		}
	}

	return nil
}

func (n *Nodeman) AddNode(config NodeConfig) (LightningNode, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Errorf("unable to generate uuid: %v", err)
	}

	switch config := config.(type) {
	case *RemoteLndNodeConfig:
		n.log.Infof("adding remote lnd node with id %s", id)

		err := n.db.SaveNode(&sweetdb.RemoteLndNode{
			Id:       id.String(),
			Name:     config.Name,
			Url:      config.Uri,
			Cert:     config.Cert,
			Macaroon: config.Macaroon,
			Enabled:  false,
		})
		if err != nil {
			return nil, errors.Errorf("unable to save: %v", err)
		}

		lndNode, err := lightning.NewLndNode(&lightning.LndNodeConfig{
			Uri:           config.Uri,
			CertBytes:     config.Cert,
			MacaroonBytes: config.Macaroon,
			Logger:        n.logCreator(id.String()),
		})
		if err != nil {
			return nil, errors.Errorf("unable to create: %v", err)
		}

		node := &RemoteLndNode{
			LndNode: lndNode,
			id:      id.String(),
			name:    config.Name,
			enabled: false,
			Uri:     config.Uri,
		}

		n.nodes = append(n.nodes, node)

		return node, nil
	case *LocalNodeConfig:
		n.log.Infof("adding local node with id %s", id)

		onionKey, err := onion.GeneratePrivateKey(onion.V2)
		if err != nil {
			return nil, errors.Errorf("unable to generate onion key: %v", err)
		}

		payload := x509.MarshalPKCS1PrivateKey(onionKey)

		err = n.db.SaveNode(&sweetdb.LocalNode{
			Id:       id.String(),
			Name:     config.Name,
			Enabled:  false,
			OnionKey: payload,
		})
		if err != nil {
			return nil, errors.Errorf("unable to save: %v", err)
		}

		onionSvc := onion.NewService(&onion.ServiceConfig{
			Tor:    n.tor,
			Logger: n.log,
			Port:   8080,
			Key:    onionKey,
		})

		localNode, err := lightning.NewLocalNode(&lightning.LocalNodeConfig{
			DataDir:  filepath.Join(n.nodesDataDir, id.String()),
			Logger:   n.logCreator(id.String()),
			OnionSvc: onionSvc,
		})
		if err != nil {
			return nil, errors.Errorf("unable to create: %v", err)
		}

		node := &LocalNode{
			LocalNode: localNode,
			id:        id.String(),
			name:      config.Name,
			enabled:   false,
		}

		n.nodes = append(n.nodes, node)

		return node, nil
	default:
		return nil, errors.Errorf("unknown config type %T", config)
	}
}

func (n *Nodeman) RemoveNode(id string) error {
	err := n.db.RemoveNode(id)
	if err != nil {
		return errors.Errorf("unable to delete: %v", err)
	}

	index := -1
	for i, node := range n.nodes {
		if node.ID() == id {
			index = i
		}
	}

	if index < 0 {
		return errors.Errorf("unavailable")
	}

	copy(n.nodes[index:], n.nodes[index+1:]) // shift a[i+1:] left one index
	n.nodes[len(n.nodes)-1] = nil            // erase last element
	n.nodes = n.nodes[:len(n.nodes)-1]       // truncate slice

	return nil
}

func (n *Nodeman) EnableNode(id string) error {
	node, err := n.db.GetNode(id)
	if err != nil {
		return errors.Errorf("unable to get node: %v", err)
	}

	switch node := node.(type) {
	case sweetdb.RemoteLndNode:
		node.Enabled = true
	case sweetdb.LocalNode:
		node.Enabled = true
	}

	err = n.db.SaveNode(node)
	if err != nil {
		return errors.Errorf("unable to save node: %v", err)
	}

	for _, node := range n.nodes {
		if node.ID() == id {
			node.setEnabled(true)

			return nil
		}
	}

	return errors.Errorf("node with id %s not found", id)
}

func (n *Nodeman) DisableNode(id string) error {
	node, err := n.db.GetNode(id)
	if err != nil {
		return errors.Errorf("unable to get node: %v", err)
	}

	switch node := node.(type) {
	case sweetdb.RemoteLndNode:
		node.Enabled = false
	case sweetdb.LocalNode:
		node.Enabled = false
	}

	err = n.db.SaveNode(node)
	if err != nil {
		return errors.Errorf("unable to save node: %v", err)
	}

	for _, node := range n.nodes {
		if node.ID() == id {
			node.setEnabled(false)

			return nil
		}
	}

	return errors.Errorf("node with id %s not found", id)
}

func (n *Nodeman) RenameNode(id string, name string) error {
	node, err := n.db.GetNode(id)
	if err != nil {
		return errors.Errorf("unable to get node: %v", err)
	}

	switch node := node.(type) {
	case sweetdb.RemoteLndNode:
		node.Name = name
	case sweetdb.LocalNode:
		node.Name = name
	}

	err = n.db.SaveNode(node)
	if err != nil {
		return errors.Errorf("unable to save node: %v", err)
	}

	for _, node := range n.nodes {
		if node.ID() == id {
			node.setName(name)
		}
	}

	return errors.Errorf("node with id %s not found", id)
}
