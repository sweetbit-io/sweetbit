package sweetdb

import (
	"github.com/go-errors/errors"
	"time"
)

var (
	updatesBucket       = []byte("updates")
	currentUpdateBucket = []byte("currentUpdate")
	updateIdKey         = []byte("id")
)

type Update struct {
	Id      string    `json:"id"`
	Started time.Time `json:"started"`
	Url     string    `json:"url"`
	State   string    `json:"state"`
}

func (db *DB) SaveUpdate(update *Update) error {
	return db.setJSON(updatesBucket, []byte(update.Id), update)
}

func (db *DB) GetUpdate(id string) (*Update, error) {
	var update = &Update{}

	if err := db.getJSON(updatesBucket, []byte(id), &update); err != nil {
		return nil, err
	}

	return update, nil
}

func (db *DB) GetCurrentUpdate() (*Update, error) {
	id, err := db.getString(currentUpdateBucket, updateIdKey)
	if err != nil {
		return nil, errors.Errorf("unable to get current update id: %v", err)
	}

	if id == "" {
		// no current update available
		return nil, nil
	}

	update, err := db.GetUpdate(id);
	if err != nil {
		return nil, errors.Errorf("unable to get update: %v", err)
	}

	return update, nil
}

func (db *DB) SetCurrentUpdate(id string) error {
	err := db.setString(currentUpdateBucket, updateIdKey, id)
	if err != nil {
		return errors.Errorf("unable to set current update: %v", err)
	}

	return nil
}

func (db *DB) ClearCurrentUpdate() error {
	err := db.setString(currentUpdateBucket, updateIdKey, "")
	if err != nil {
		return errors.Errorf("unable to clear current update: %v", err)
	}

	return nil
}
