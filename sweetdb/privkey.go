package sweetdb

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"github.com/go-errors/errors"
	"go.etcd.io/bbolt"
)

func (db *DB) setPrivateKey(bucket []byte, bucketKey []byte, key *rsa.PrivateKey) error {
	payload := x509.MarshalPKCS1PrivateKey(key)

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

func (db *DB) getPrivateKey(bucket []byte, bucketKey []byte) (*rsa.PrivateKey, error) {
	var key *rsa.PrivateKey

	err := db.View(func(tx *bbolt.Tx) error {
		// First fetch the bucket
		bucket := tx.Bucket(bucket)
		if bucket == nil {
			return nil
		}

		posPrivateKeyBytes := bucket.Get(bucketKey)
		if posPrivateKeyBytes == nil || bytes.Equal(posPrivateKeyBytes, []byte("null")) {
			return nil
		}

		var err error
		key, err = x509.ParsePKCS1PrivateKey(posPrivateKeyBytes)
		if err != nil {
			return errors.Errorf("Could not unmarshal data: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return key, nil
}
