package sweetdb

import (
	"github.com/go-errors/errors"
	"go.etcd.io/bbolt"
)

var (
	nodesBucket = []byte("nodes")
)

type lightningNodeKind string

const (
	lightningNodeKindLocal  lightningNodeKind = "local"
	lightningNodeKindRemote                   = "remote"
)

type lightningNode struct {
	Kind lightningNodeKind `json:"kind"`
}

type LightningNode interface{}

type RemoteLndNode struct {
	lightningNode
	Id       string `json:"id"`
	Name     string `json:"name"`
	Enabled  bool   `json:"enabled"`
	Mainnet  bool   `json:"mainnet"`
	Url      string `json:"url"`
	Cert     []byte `json:"cert"`
	Macaroon []byte `json:"macaroon"`
}

type LocalNode struct {
	lightningNode
	Id       string `json:"id"`
	Name     string `json:"name"`
	Enabled  bool   `json:"enabled"`
	Mainnet  bool   `json:"mainnet"`
	OnionKey []byte `json:"onionkey"`
}

func (db *DB) GetNodes() ([]LightningNode, error) {
	keys, err := db.getKeys(nodesBucket)
	if err != nil {
		return nil, errors.Errorf("unable to get keys: %v", err)
	}

	nodes := []LightningNode{}

	for _, k := range keys {
		id := string(k)
		node, err := db.GetNode(id)
		if err != nil {
			return nil, errors.Errorf("unable to get node with id %s: %v", id, err)
		}

		if node == nil {
			return nil, errors.Errorf("unable to find node with id %s", id)
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (db *DB) SaveNode(node LightningNode) error {
	switch n := node.(type) {
	case *RemoteLndNode:
		n.Kind = lightningNodeKindRemote
		return db.setJSON(nodesBucket, []byte(n.Id), n)
	case *LocalNode:
		n.Kind = lightningNodeKindLocal
		return db.setJSON(nodesBucket, []byte(n.Id), n)
	default:
		return errors.Errorf("Can only save nodes, got %T", node)
	}
}

func (db *DB) RemoveNode(id string) error {
	return db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(nodesBucket)
		if err != nil {
			return err
		}

		if err := bucket.Delete([]byte(id)); err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) GetNode(id string) (LightningNode, error) {
	var node *lightningNode

	if err := db.getJSON(nodesBucket, []byte(id), &node); err != nil {
		return nil, err
	}

	if node == nil {
		return nil, nil
	}

	switch node.Kind {
	case lightningNodeKindRemote:
		var node *RemoteLndNode
		if err := db.getJSON(nodesBucket, []byte(id), &node); err != nil {
			return nil, err
		}
		return node, nil
	case lightningNodeKindLocal:
		var node *LocalNode
		if err := db.getJSON(nodesBucket, []byte(id), &node); err != nil {
			return nil, err
		}
		return node, nil
	default:
		return nil, errors.Errorf("unknown node type %s", node.Kind)
	}
}
