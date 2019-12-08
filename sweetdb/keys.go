package sweetdb

import (
	"github.com/go-errors/errors"
	"go.etcd.io/bbolt"
)

func (db *DB) getKeys(bucketName []byte) ([][]byte, error) {
	keys := [][]byte{}

	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return nil
		}

		err := bucket.ForEach(func(k, v []byte) error {
			keys = append(keys, k)
			return nil
		})
		if err != nil {
			return errors.Errorf("unable to loop through keys: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, errors.Errorf("unable to view: %v", err)
	}

	return keys, nil
}
