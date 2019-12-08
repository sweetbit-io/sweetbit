package sweetdb

import (
	"crypto/rsa"
	bolt "go.etcd.io/bbolt"
)

var (
	settingsBucket     = []byte("settings")
	lightningNodeKey   = []byte("lightningNode")
	nameKey            = []byte("name")
	dispenseOnTouchKey = []byte("dispenseOnTouch")
	buzzOnDispenseKey  = []byte("buzzOnDispense")
	posPrivateKeyKey   = []byte("posPrivateKey")
	apiPrivateKeyKey   = []byte("apiPrivateKey")
)

func (db *DB) SetPosPrivateKey(key *rsa.PrivateKey) error {
	return db.setPrivateKey(settingsBucket, posPrivateKeyKey, key)
}

func (db *DB) GetPosPrivateKey() (*rsa.PrivateKey, error) {
	return db.getPrivateKey(settingsBucket, posPrivateKeyKey)
}

func (db *DB) SetApiPrivateKey(key *rsa.PrivateKey) error {
	return db.setPrivateKey(settingsBucket, apiPrivateKeyKey, key)
}

func (db *DB) GetApiPrivateKey() (*rsa.PrivateKey, error) {
	return db.getPrivateKey(settingsBucket, apiPrivateKeyKey)
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
	return db.setJSON(settingsBucket, dispenseOnTouchKey, dispenseOnTouch)
}

func (db *DB) GetDispenseOnTouch() (bool, error) {
	var dispenseOnTouch bool

	if err := db.getJSON(settingsBucket, dispenseOnTouchKey, &dispenseOnTouch); err != nil {
		return false, err
	}

	return dispenseOnTouch, nil
}

func (db *DB) SetBuzzOnDispense(buzzOnDispense bool) error {
	return db.setJSON(settingsBucket, buzzOnDispenseKey, buzzOnDispense)
}

func (db *DB) GetBuzzOnDispense() (bool, error) {
	var buzzOnDispense bool

	if err := db.getJSON(settingsBucket, buzzOnDispenseKey, &buzzOnDispense); err != nil {
		return false, err
	}

	return buzzOnDispense, nil
}
