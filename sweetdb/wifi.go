package sweetdb

import (
	"github.com/go-errors/errors"
)

var (
	wifiBucket = []byte("wifi")
	wifiKey    = []byte("default")
)

type wifiEncryption string

const (
	wifiEncryptionNone       wifiEncryption = "none"
	wifiEncryptionPersonal                  = "personal"
	wifiEncryptionEnterprise                = "enterprise"
)

type wifiBase struct {
	Encryption wifiEncryption `json:"encryption"`
}

type Wifi interface{}

type WifiPublic struct {
	wifiBase
	Ssid string `json:"ssid"`
}

type WifiPersonal struct {
	wifiBase
	Ssid string `json:"ssid"`
	Psk  string `json:"psk"`
}

type WifiEnterprise struct {
	wifiBase
	Ssid     string `json:"ssid"`
	Identity string `json:"identity"`
	Password string `json:"password"`
}

func (db *DB) SaveWifi(wifi Wifi) error {
	switch w := wifi.(type) {
	case *WifiPublic:
		w.Encryption = wifiEncryptionNone
		return db.setJSON(wifiBucket, wifiKey, w)
	case *WifiPersonal:
		w.Encryption = wifiEncryptionPersonal
		return db.setJSON(wifiBucket, wifiKey, w)
	case *WifiEnterprise:
		w.Encryption = wifiEncryptionEnterprise
		return db.setJSON(wifiBucket, wifiKey, w)
	default:
		return errors.Errorf("Can only save wifi, got %T", w)
	}
}

func (db *DB) GetWifi() (Wifi, error) {
	var wifiBase *wifiBase

	if err := db.getJSON(wifiBucket, wifiKey, &wifiBase); err != nil {
		return nil, err
	}

	if wifiBase == nil {
		return nil, nil
	}

	switch wifiBase.Encryption {
	case wifiEncryptionNone:
		var wifi *WifiPublic
		if err := db.getJSON(wifiBucket, wifiKey, &wifi); err != nil {
			return nil, err
		}
		return wifi, nil
	case wifiEncryptionPersonal:
		var wifi *WifiPersonal
		if err := db.getJSON(wifiBucket, wifiKey, &wifi); err != nil {
			return nil, err
		}
		return wifi, nil
	case wifiEncryptionEnterprise:
		var wifi *WifiEnterprise
		if err := db.getJSON(wifiBucket, wifiKey, &wifi); err != nil {
			return nil, err
		}
		return wifi, nil
	default:
		return nil, errors.Errorf("unknown wifi type %s", wifiBase.Encryption)
	}
}
