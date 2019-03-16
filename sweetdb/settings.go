package sweetdb

import (
	"bytes"
	"encoding/json"
	"github.com/go-errors/errors"
	bolt "go.etcd.io/bbolt"
)

var (
	settingsBucket     = []byte("settings")
	lightningNodeKey   = []byte("lightningNode")
	nameKey            = []byte("name")
	dispenseOnTouchKey = []byte("dispenseOnTouch")
	buzzOnDispenseKey  = []byte("buzzOnDispense")
	wifiConnectionKey  = []byte("wifi")
)

type LightningNode struct {
	Uri      string `json:"uri"`
	Cert     []byte `json:"cert"`
	Macaroon []byte `json:"macaroon"`
}

type WifiConnection struct {
	Ssid string `json:ssid`
	Psk  string `json:psk`
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

func (db *DB) SetWifiConnection(wifiConnection *WifiConnection) error {
	payload, err := json.Marshal(wifiConnection)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		// First grab the settings bucket
		bucket, err := tx.CreateBucketIfNotExists(settingsBucket)
		if err != nil {
			return err
		}

		if err := bucket.Put(wifiConnectionKey, payload); err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) GetWifiConnection() (*WifiConnection, error) {
	var wifiConnection = &WifiConnection{}

	err := db.View(func(tx *bolt.Tx) error {
		// First fetch the bucket
		bucket := tx.Bucket(settingsBucket)
		if bucket == nil {
			return nil
		}

		wifiConnectionBytes := bucket.Get(wifiConnectionKey)
		if wifiConnectionBytes == nil || bytes.Equal(wifiConnectionBytes, []byte("null")) {
			wifiConnection = nil
			return nil
		}

		err := json.Unmarshal(wifiConnectionBytes, &wifiConnection)
		if err != nil {
			return errors.Errorf("Could not unmarshal data: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return wifiConnection, nil
}

func (db *DB) SetName(name string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(settingsBucket)
		if err != nil {
			return err
		}

		if err := bucket.Put(nameKey, []byte(name)); err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) GetName() (string, error) {
	var name string

	err := db.View(func(tx *bolt.Tx) error {
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

func (db *DB) SetDispenseOnTouch(dispenseOnTouch bool) error {
	payload, err := json.Marshal(dispenseOnTouch)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(settingsBucket)
		if err != nil {
			return err
		}

		if err := bucket.Put(dispenseOnTouchKey, payload); err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) GetDispenseOnTouch() (bool, error) {
	var dispenseOnTouch bool

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(settingsBucket)
		if bucket == nil {
			return nil
		}

		dispenseOnTouchBytes := bucket.Get(dispenseOnTouchKey)

		err := json.Unmarshal(dispenseOnTouchBytes, &dispenseOnTouch)
		if err != nil {
			return errors.Errorf("Could not unmarshal data: %v", err)
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return dispenseOnTouch, nil
}

func (db *DB) SetBuzzOnDispense(buzzOnDispense bool) error {
	payload, err := json.Marshal(buzzOnDispense)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(settingsBucket)
		if err != nil {
			return err
		}

		if err := bucket.Put(buzzOnDispenseKey, payload); err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) GetBuzzOnDispense() (bool, error) {
	var buzzOnDispense bool

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(settingsBucket)
		if bucket == nil {
			return nil
		}

		buzzOnDispenseBytes := bucket.Get(buzzOnDispenseKey)

		err := json.Unmarshal(buzzOnDispenseBytes, &buzzOnDispense)
		if err != nil {
			return errors.Errorf("Could not unmarshal data: %v", err)
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return buzzOnDispense, nil
}
