package sweetdb

import (
	bolt "go.etcd.io/bbolt"
)

func (db *DB) setString(bucket []byte, bucketKey []byte, v string) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		if err := bucket.Put(bucketKey, []byte(v)); err != nil {
			return err
		}

		return nil
	})
}

func (db *DB) getString(bucketName []byte, bucketKey []byte) (string, error) {
	var v string

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return nil
		}

		nameBytes := bucket.Get(bucketKey)
		v = string(nameBytes)

		return nil
	})

	if err != nil {
		return "", err
	}

	return v, nil
}
