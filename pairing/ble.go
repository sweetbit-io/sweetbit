package pairing

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	"github.com/the-lightning-land/sweetd/network"
	"github.com/the-lightning-land/sweetd/pairing/ble"
	"time"
)

// check compliance to interface during compile time
var _ Controller = (*BLEController)(nil)

const (
	appPath           = "/io/sweetbit/app"
	advertisementPath = "/io/sweetbit/advertisement"

	// unique uuid suffix for the candy dispenser
	uuidSuffix = "-75dd-4a0e-b688-66b7df342cc6"

	candyServiceUuid       = "ca000000" + uuidSuffix
	statusCharUuid         = "ca001000" + uuidSuffix
	scanWifiCharUuid       = "ca002000" + uuidSuffix
	discoveredWifiCharUuid = "ca003000" + uuidSuffix
	connectWifiCharUuid    = "ca004000" + uuidSuffix
	onionApiCharUuid       = "ca005000" + uuidSuffix
)

type Dispenser interface {
	ConnectWifi(network.Connection) error
	GetApiOnionID() string
	GetName() string
	ScanWifi() (*network.ScanClient, error)
}

type BLEControllerConfig struct {
	Logger    Logger
	AdapterId string
	Dispenser Dispenser
}

type BLEController struct {
	log                  Logger
	adapterId            string
	dispenser            Dispenser
	app                  *ble.GattApp
	notifyDisocveredWifi ble.Writer
	discoveredWifis      chan *network.Wifi
}

func NewController(config *BLEControllerConfig) (*BLEController, error) {
	controller := &BLEController{}

	if config.Logger != nil {
		controller.log = config.Logger
	} else {
		controller.log = noopLogger{}
	}

	controller.discoveredWifis = make(chan *network.Wifi, 10)

	// Assign the device adapter id (ex. hci0)
	controller.adapterId = config.AdapterId

	// Most pairing actions rely on functions the dispenser is providing
	controller.dispenser = config.Dispenser

	controller.app = ble.NewGattApp(
		config.AdapterId,
		appPath,
		ble.WithAppService(candyServiceUuid,
			ble.WithServiceCharacteristic(
				statusCharUuid,
				ble.WithCharacteristicReadHandler(controller.status),
				ble.WithCharacteristicUserDescriptionDescriptor("Status"),
				ble.WithCharacteristicPresentationFormatDescriptor(),
			),
			ble.WithServiceCharacteristic(
				scanWifiCharUuid,
				ble.WithCharacteristicWriteHandler(controller.scanWifi),
				ble.WithCharacteristicUserDescriptionDescriptor("Scan Wi-Fi"),
			),
			ble.WithServiceCharacteristic(
				discoveredWifiCharUuid,
				ble.WithCharacteristicWriter(&controller.notifyDisocveredWifi),
				ble.WithCharacteristicUserDescriptionDescriptor("Discovered Wi-Fi"),
			),
			ble.WithServiceCharacteristic(
				connectWifiCharUuid,
				ble.WithCharacteristicWriteHandler(controller.connectWifi),
				ble.WithCharacteristicUserDescriptionDescriptor("Connect Wi-Fi"),
			),
			ble.WithServiceCharacteristic(
				onionApiCharUuid,
				ble.WithCharacteristicReadHandler(controller.onionApi),
				ble.WithCharacteristicUserDescriptionDescriptor("Onion API"),
			),
		),
	)

	return controller, nil
}

func (c *BLEController) Advertise() error {
	adapter, err := api.GetAdapter(c.adapterId)
	if err != nil {
		return errors.Errorf("unable to get adapter: %v", err)
	}

	err = adapter.SetAlias("Candy")
	if err != nil {
		return errors.Errorf("unable to set alias: %v", err)
	}

	err = adapter.SetDiscoverable(true)
	if err != nil {
		return errors.Errorf("unable to make discoverable: %v", err)
	}

	err = adapter.SetDiscoverableTimeout(0)
	if err != nil {
		return errors.Errorf("unable to set discoverable timeout: %v", err)
	}

	err = adapter.SetPowered(true)
	if err != nil {
		return errors.Errorf("unable to set powered: %v", err)
	}

	advertisementProps := &advertising.LEAdvertisement1Properties{
		Type:      advertising.AdvertisementTypePeripheral,
		LocalName: "Candy",
		ServiceUUIDs: []string{
			candyServiceUuid,
		},
		Timeout:             0,
		Duration:            0,
		DiscoverableTimeout: 0,
		Discoverable:        true,
	}

	// advertisementPath := dbus.ObjectPath(advertisementPath)

	advertisement, err := api.NewAdvertisement(c.adapterId, advertisementProps)
	if err != nil {
		return errors.Errorf("unable to create advertisement: %v", err)
	}

	err = api.ExposeDBusService(advertisement)
	if err != nil {
		return errors.Errorf("unable to expose advertisement: &%v", err)
	}

	advertisingManager, err := advertising.NewLEAdvertisingManager1FromAdapterID(c.adapterId)
	if err != nil {
		return errors.Errorf("unable to create advertisement manager: %v", err)
	}

	err = advertisingManager.RegisterAdvertisement(advertisement.Path(), map[string]interface{}{})
	if err != nil {
		return errors.Errorf("unable to register advertisement manager: %v", err)
	}

	return nil
}

func (c *BLEController) Start() error {
	err := c.app.Start()
	if err != nil {
		return errors.Errorf("unable to start app: %v", err)
	}

	go func() {
		type discoveredWifi struct {
			Ssid       string `json:"ssid"`
			Encryption string `json:"encryption"`
		}

		for {
			w, ok := <-c.discoveredWifis
			if !ok {
				break
			}

			wifi := discoveredWifi{
				Ssid: w.Ssid,
			}

			if w.Encryption == network.EncryptionPersonal {
				wifi.Encryption = "personal"
			} else if w.Encryption == network.EncryptionEnterprise {
				wifi.Encryption = "enterprise"
			} else {
				wifi.Encryption = "none"
			}

			value, err := json.Marshal(wifi)
			if err != nil {
				break
			}

			err = c.notifyDisocveredWifi(value)
			if err != nil {
				c.log.Errorf("unable to write discovered wifi: %v", err)
			}

			// wait 100ms between wifi discovery notifications so that the client can keep up
			// receiving them
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return nil
}

func (c *BLEController) Stop() error {
	err := c.app.Stop()
	if err != nil {
		return errors.Errorf("unable to stop app: %v", err)
	}

	return nil
}

func (c *BLEController) status() ([]byte, error) {
	type status struct {
		Name string `json:"name"`
	}

	s := status{
		Name: c.dispenser.GetName(),
	}

	value, err := json.Marshal(s)
	if err != nil {
		return nil, errors.Errorf("unable to serialize: %v", err)
	}

	return value, nil
}

func (c *BLEController) scanWifi(value []byte) error {
	client, err := c.dispenser.ScanWifi()
	if err != nil {
		return errors.Errorf("unable to scan: %v", err)
	}

	for {
		// this channel will close when the scan is finished
		wifi, ok := <-client.Wifis
		if !ok {
			break
		}

		c.discoveredWifis <- wifi
	}

	return nil
}

func (c *BLEController) connectWifi(value []byte) error {
	type config struct {
		Encryption string `json:"encryption"`
		Ssid       string `json:"ssid"`
		Psk        string `json:"psk"`
		Identity   string `json:"identity"`
		Password   string `json:"password"`
	}

	var cfg config

	err := json.Unmarshal(value, &cfg)
	if err != nil {
		return errors.Errorf("unable to deserialize: %v", err)
	}

	var conn network.Connection

	switch cfg.Encryption {
	case "none":
		conn = &network.WpaConnection{
			Ssid: cfg.Ssid,
		}
	case "personal":
		conn = &network.WpaPersonalConnection{
			Ssid: cfg.Ssid,
			Psk:  cfg.Psk,
		}
	case "enterprise":
		conn = &network.WpaEnterpriseConnection{
			Ssid:     cfg.Ssid,
			Identity: cfg.Identity,
			Password: cfg.Password,
		}
	default:
		return errors.Errorf("unknown encryption type %s", cfg.Encryption)
	}

	err = c.dispenser.ConnectWifi(conn)
	if err != nil {
		return errors.Errorf("unable to connect: %v", err)
	}

	return nil
}

func (c *BLEController) onionApi() ([]byte, error) {
	return []byte(c.dispenser.GetApiOnionID()), nil
}
