package dispenser

import (
	"github.com/go-errors/errors"
	"github.com/the-lightning-land/sweetd/network"
	"github.com/the-lightning-land/sweetd/sweetdb"
	"sync"
)

// maybeAttemptSavedWifiConnection is run as a goroutine and attempts a connection
// to the most recently persisted wifi connection, if no network connection is available yet
func (d *Dispenser) maybeAttemptSavedWifiConnection(wg sync.WaitGroup) {
	wg.Add(1)

	wifiConnection, err := d.db.GetWifi()
	if err != nil {
		d.log.Warnf("could not get wifi connection: %v", err)
	}

	if wifiConnection != nil {
		var err error

		switch conn := wifiConnection.(type) {
		case *sweetdb.WifiPublic:
			err = d.network.Connect(&network.WpaConnection{
				Ssid: conn.Ssid,
			})
		case *sweetdb.WifiPersonal:
			err = d.network.Connect(&network.WpaPersonalConnection{
				Ssid: conn.Ssid,
				Psk:  conn.Psk,
			})
		case *sweetdb.WifiEnterprise:
			err = d.network.Connect(&network.WpaEnterpriseConnection{
				Ssid:     conn.Ssid,
				Identity: conn.Identity,
				Password: conn.Password,
			})
		default:
			d.log.Errorf("unknown connection type %T", wifiConnection)
			wg.Done()
			return
		}
		if err != nil {
			d.log.Errorf("could not connect to saved wifi: %v", err)
		}
	} else {
		d.log.Debugf("no saved wifi connection was found")
	}

	wg.Done()
}

func (d *Dispenser) ConnectToWifi(connection network.Connection) error {
	switch conn := connection.(type) {
	case *network.WpaPersonalConnection:
		d.log.Infof("Connecting to personal wifi %v", conn.Ssid)

		err := d.network.Connect(conn)
		if err != nil {
			return errors.Errorf("unable to connect: %v", err)
		}

		err = d.SetWifiConnection(&sweetdb.WifiPersonal{
			Ssid: conn.Ssid,
			Psk:  conn.Psk,
		})
		if err != nil {
			d.log.Errorf("Could not save wifi connection: %v", err)
		}
	case *network.WpaConnection:
		d.log.Infof("Connecting to public wifi %v", conn.Ssid)

		err := d.network.Connect(conn)
		if err != nil {
			return errors.Errorf("unable to connect: %v", err)
		}

		err = d.SetWifiConnection(&sweetdb.WifiPublic{
			Ssid: conn.Ssid,
		})
		if err != nil {
			d.log.Errorf("Could not save wifi connection: %v", err)
		}
	case *network.WpaEnterpriseConnection:
		d.log.Infof("Connecting to enterprise wifi %v", conn.Ssid)

		err := d.network.Connect(conn)
		if err != nil {
			return errors.Errorf("unable to connect: %v", err)
		}

		err = d.SetWifiConnection(&sweetdb.WifiEnterprise{
			Ssid:     conn.Ssid,
			Identity: conn.Identity,
			Password: conn.Password,
		})
		if err != nil {
			d.log.Errorf("Could not save wifi connection: %v", err)
		}
	default:
		return errors.Errorf("unsupported connection type %T", connection)
	}

	return nil
}

func (d *Dispenser) SetWifiConnection(connection sweetdb.Wifi) error {
	d.log.Infof("Setting Wifi connection")

	err := d.db.SaveWifi(connection)
	if err != nil {
		return errors.Errorf("Failed setting Wifi connection: %v", err)
	}

	return nil
}

func (d *Dispenser) ScanWifi() (*network.ScanClient, error) {
	return d.network.Scan()
}
