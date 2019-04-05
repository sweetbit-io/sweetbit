package pairing

import (
	"bytes"
	"github.com/go-errors/errors"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/linux/btmgmt"
	"github.com/muka/go-bluetooth/service"
	"github.com/the-lightning-land/sweetd/ap"
	"github.com/the-lightning-land/sweetd/dispenser"
	"strings"
	"time"
)

const (
	// Unique UUID suffix for the candy dispenser
	uuidSuffix = "-75dd-4a0e-b688-66b7df342cc6"

	// Prefix of the candy service UUID
	candyServiceUuidPrefix = "CA00"

	// Where to expose the application
	objectName = "land.lightning"
	objectPath = "/sweet/pairing/service"

	// Local name of the application
	localName = "Candy"

	candyServiceUuid          = candyServiceUuidPrefix + "0000" + uuidSuffix
	networkAvailabilityStatus = candyServiceUuidPrefix + "0001" + uuidSuffix
	ipAddress                 = candyServiceUuidPrefix + "0002" + uuidSuffix
	wifiScanSignal            = candyServiceUuidPrefix + "0003" + uuidSuffix
	wifiScanList              = candyServiceUuidPrefix + "0004" + uuidSuffix
	wifiSsidString            = candyServiceUuidPrefix + "0005" + uuidSuffix
	wifiPskString             = candyServiceUuidPrefix + "0006" + uuidSuffix
	wifiConnectSignal         = candyServiceUuidPrefix + "0007" + uuidSuffix
)

type Controller struct {
	log                        Logger
	adapterId                  string
	accessPoint                ap.Ap
	apClient                   *ap.ApClient
	app                        *service.Application
	service                    *service.GattService1
	ssid                       string
	psk                        string
	dispenser                  *dispenser.Dispenser
	ssidChanges                chan []byte
	ipChanges                  chan []byte
	networkAvailabilityChanges chan []byte
}

func NewController(config *Config) (*Controller, error) {
	controller := &Controller{}

	if config.Logger != nil {
		controller.log = config.Logger
	} else {
		controller.log = noopLogger{}
	}

	// Assign the device adapter id (ex. hci0)
	controller.adapterId = config.AdapterId

	// Assign the depending access point
	controller.accessPoint = config.AccessPoint

	// Most pairing actions rely on functions the dispenser is providing
	controller.dispenser = config.Dispenser

	controller.ipChanges = make(chan []byte)
	controller.ssidChanges = make(chan []byte)
	controller.networkAvailabilityChanges = make(chan []byte)

	var err error

	app := GattApp(objectName, objectPath, localName)
	service := app.Service(Primary, candyServiceUuid, Advertised)

	service.DeviceNameCharacteristic().
		WithValue("Candy").
		Create().
		UserDescriptionDescriptor("Device Name").
		PresentationDescriptor()

	service.ManufacturerNameCharacteristic().
		WithValue("The Lightning Land").
		Create().
		UserDescriptionDescriptor("Manufacturer Name").
		PresentationDescriptor()

	service.SerialNumberCharacteristic().
		WithValue("123456789").
		Create().
		UserDescriptionDescriptor("Serial Number").
		PresentationDescriptor()

	service.ModelNumberCharacteristic().
		WithValue("moon").
		Create().
		UserDescriptionDescriptor("Model Number").
		PresentationDescriptor()

	service.Characteristic(networkAvailabilityStatus).
		WithReadHandler(controller.readNetworkAvailabilityStatus).
		WithNotifications(controller.networkAvailabilityChanges).
		Create().
		UserDescriptionDescriptor("Network Availability Status")

	service.Characteristic(ipAddress).
		WithReadHandler(controller.readIpAddress).
		WithNotifications(controller.ipChanges).
		Create().
		UserDescriptionDescriptor("IP Address")

	service.Characteristic(wifiScanSignal).
		WithWriteHandler(controller.writeWifiScanSignal).
		Create().
		UserDescriptionDescriptor("Wi-Fi Scan Signal")

	service.Characteristic(wifiScanList).
		WithReadHandler(controller.readWifiScanList).
		Create().
		UserDescriptionDescriptor("Wi-Fi Scan List")

	service.Characteristic(wifiSsidString).
		WithReadHandler(controller.readWifiSsidString).
		WithWriteHandler(controller.writeWifiSsidString).
		WithNotifications(controller.ssidChanges).
		Create().
		UserDescriptionDescriptor("Wi-Fi SSID")

	service.Characteristic(wifiPskString).
		WithWriteHandler(controller.writeWifiPskString).
		Create().
		UserDescriptionDescriptor("Wi-Fi PSK")

	service.Characteristic(wifiConnectSignal).
		WithWriteHandler(controller.writeWifiConnectSignal).
		Create().
		UserDescriptionDescriptor("Wi-Fi Connect Signal")

	controller.app, err = app.Run()
	if err != nil {
		return nil, errors.Errorf("Could not start app: %v", err)
	}

	return controller, nil
}

func (c *Controller) Start() error {
	mgmt := btmgmt.NewBtMgmt(c.adapterId)
	err := mgmt.Reset()
	if err != nil {
		return errors.Errorf("Reset %s: %v", c.adapterId, err)
	}

	// Sleep to give the device some time after the reset
	time.Sleep(time.Millisecond * 500)

	gattManager, err := api.GetGattManager(c.adapterId)
	if err != nil {
		return errors.Errorf("Get gatt manager failed: %v", err)
	}

	err = gattManager.RegisterApplication(c.app.Path(), map[string]interface{}{})
	if err != nil {
		return errors.Errorf("Register failed: %v", err)
	}

	err = c.app.StartAdvertising(c.adapterId)
	if err != nil {
		return errors.Errorf("Failed to advertise: %v", err)
	}

	c.apClient = c.accessPoint.SubscribeUpdates()
	go func() {
		for {
			update := <-c.apClient.Updates

			c.log.Infof("Got update: %v", update)

			c.ssidChanges <- []byte(update.Ssid)
			c.ipChanges <- []byte(update.Ip)

			if update.Connected {
				c.networkAvailabilityChanges <- []byte{1}
			} else {
				c.networkAvailabilityChanges <- []byte{0}
			}
		}
	}()

	return nil
}

func (c *Controller) Stop() error {
	c.apClient.Cancel()

	err := c.app.StopAdvertising()
	if err != nil {
		return errors.Errorf("Could not stop advertising: %v", err)
	}

	gattManager, err := api.GetGattManager(c.adapterId)
	if err != nil {
		return errors.Errorf("Get gatt manager failed: %v", err)
	}

	err = gattManager.UnregisterApplication(c.app.Path())
	if err != nil {
		return errors.Errorf("Unregister failed: %v", err)
	}

	return nil
}

func (c *Controller) readNetworkAvailabilityStatus() ([]byte, error) {
	c.log.Infof("Reading network availability...")

	status, err := c.accessPoint.GetConnectionStatus()
	if err != nil {
		return nil, errors.Errorf("Could not get wifi status: %v", err)
	}

	var connected uint8

	if status.State == "COMPLETED" {
		connected = 1
	} else {
		connected = 0
	}

	return []byte{connected}, nil
}

func (c *Controller) readIpAddress() ([]byte, error) {
	c.log.Infof("Reading ip address...")

	status, err := c.accessPoint.GetConnectionStatus()
	if err != nil {
		return nil, errors.Errorf("Could not get wifi status: %v", err)
	}

	return []byte(status.Ip), nil
}

type WifiScanListItem struct {
	Ssid string `json:"ssid"`
}

func (c *Controller) writeWifiScanSignal(value []byte) error {
	c.log.Infof("Writing wifi scan signal to %v", value)

	if bytes.Equal(value, []byte{1}) {
		err := c.accessPoint.ScanWifiNetworks()
		if err != nil {
			return errors.Errorf("Could not scan wifi networks: %v", err)
		}
	}

	return nil
}

func (c *Controller) readWifiScanList() ([]byte, error) {
	c.log.Infof("Reading wifi scan list...")

	networks, err := c.accessPoint.ListWifiNetworks()
	if err != nil {
		return nil, errors.Errorf("Could not get wifi scan list: %v", err)
	}

	// Map wifi's, so only one entry per SSID is left
	wifiScanMap := make(map[string]*WifiScanListItem)
	for _, net := range networks {
		wifiScanMap[net.Ssid] = &WifiScanListItem{
			Ssid: net.Ssid,
		}
	}

	var list string

	// Use literal instead of declaration so it serializes into empty json array
	wifiScanList := []*WifiScanListItem{}
	for _, net := range wifiScanMap {
		wifiScanList = append(wifiScanList, net)
	}

	for i, net := range wifiScanList {
		if i == 0 {
			list = net.Ssid
		} else {
			list = net.Ssid + "\t" + list
		}
	}

	//payload, err := json.Marshal(wifiScanList)
	//if err != nil {
	//	return nil, errors.Errorf("Could not serialize wifi scan list: %v", err)
	//}

	c.log.Infof("Returning Wi-Fi networks: %s", list)

	return []byte(list), nil
}

func (c *Controller) readWifiSsidString() ([]byte, error) {
	c.log.Infof("Reading wifi ssid...")

	status, err := c.accessPoint.GetConnectionStatus()
	if err != nil {
		return nil, errors.Errorf("Could not get wifi status: %v", err)
	}

	return []byte(status.Ssid), nil
}

func (c *Controller) writeWifiSsidString(value []byte) error {
	ssid := string(value)

	c.log.Infof("Writing wifi ssid to %v", ssid)

	c.ssid = ssid

	return nil
}

func (c *Controller) writeWifiPskString(value []byte) error {
	psk := string(value)
	stars := strings.Repeat("*", len(psk))

	c.log.Infof("Writing wifi psk to %v", stars)

	c.psk = psk

	return nil
}

func (c *Controller) writeWifiConnectSignal(value []byte) error {
	c.log.Infof("Writing wifi connect signal to %v", value)

	if bytes.Equal(value, []byte{1}) {
		err := c.dispenser.ConnectToWifi(c.ssid, c.psk)
		if err != nil {
			return errors.Errorf("Could not connect to wifi: %v", err)
		}
	}

	return nil
}
