package sweetdb

import (
	"bytes"
	"encoding/json"
	"github.com/go-errors/errors"
	bolt "go.etcd.io/bbolt"
)

var (
	settingsBucket   = []byte("settings")
	lightningNodeKey = []byte("lightningNode")
	nameKey          = []byte("name")
)

type LightningNode struct {
	Uri      string `json:"uri"`
	Cert     []byte `json:"cert"`
	Macaroon []byte `json:"macaroon"`
}

func (db *DB) SetLightningNode(lightningNode *LightningNode) error {
	payload, err := json.Marshal(lightningNode)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		// First grab the settings bucket
		bucket, err := tx.CreateBucketIfNotExists(settingsBucket)
		if err != nil {
			return err
		}

		// Set the lightning node
		if err := bucket.Put(lightningNodeKey, payload); err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) GetLightningNode() (*LightningNode, error) {
	var lightningNode = &LightningNode{}

	err := db.View(func(tx *bolt.Tx) error {
		// First fetch the bucket
		bucket := tx.Bucket(settingsBucket)
		if bucket == nil {
			return nil
		}

		lightningNodeBytes := bucket.Get(lightningNodeKey)
		if lightningNodeBytes == nil || bytes.Equal(lightningNodeBytes, []byte("null")) {
			lightningNode = nil
			return nil
		}

		err := json.Unmarshal(lightningNodeBytes, &lightningNode)
		if err != nil {
			return errors.Errorf("Could not unmarshal data: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return lightningNode, nil
}

func (db *DB) SetName(name string) error {
	return db.Update(func(tx *bolt.Tx) error {
		// First grab the settings bucket
		bucket, err := tx.CreateBucketIfNotExists(settingsBucket)
		if err != nil {
			return err
		}

		// Set the name
		if err := bucket.Put(nameKey, []byte(name)); err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) GetName() (string, error) {
	var name string

	err := db.View(func(tx *bolt.Tx) error {
		// First fetch the bucket
		bucket := tx.Bucket(settingsBucket)
		if bucket == nil {
			return nil
		}

		nameBytes := bucket.Get(nameKey)
		name = string(nameBytes)

		return nil
	})

	if err != nil {
		return "", err
	}

	return name, nil
}
