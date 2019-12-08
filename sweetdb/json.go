package sweetdb

import (
	"bytes"
	"encoding/json"
	"github.com/go-errors/errors"
	"go.etcd.io/bbolt"
)

func (db *DB) setJSON(bucket []byte, bucketKey []byte, v interface{}) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		if err := bucket.Put(bucketKey, payload); err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) getJSON(bucketName []byte, bucketKey []byte, v interface{}) error {
	err := db.View(func(tx *bbolt.Tx) error {
		// First fetch the bucket
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return nil
		}

		payload := bucket.Get(bucketKey)
		if payload == nil || bytes.Equal(payload, []byte("null")) {
			v = nil
			return nil
		}

		err := json.Unmarshal(payload, v)
		if err != nil {
			return errors.Errorf("Could not unmarshal data: %v", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
