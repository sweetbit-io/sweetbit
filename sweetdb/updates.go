package sweetdb

import "time"

var (
	updatesBucket = []byte("updates")
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
