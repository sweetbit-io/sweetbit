// Convenient methods for populating Gatt services,
// characteristics and descriptors

package pairing

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/service"
)

type PrimaryType bool

const Primary = PrimaryType(true)
const Secondary = PrimaryType(false)

type AdvertisedType bool

const Advertised = AdvertisedType(true)
const AdvertisedOptional = AdvertisedType(false)

type HandleRead = func() ([]byte, error)
type HandleWrite = func(value []byte) error
type Notifier <-chan []byte

type gattApp struct {
	app           *service.Application
	err           error
	readHandlers  map[string]HandleRead
	writeHandlers map[string]HandleWrite
	notifiers     map[*service.GattCharacteristic1]Notifier
	done          chan struct{}
}

type gattService struct {
	*gattApp
	service *service.GattService1
}

// characteristic that is still constructed, before .Create() is called
type gattPendingCharacteristic struct {
	*gattService
	characteristicUuid          string
	characteristicValue         []byte
	characteristicRead          HandleRead
	characteristicWrite         HandleWrite
	characteristicNotifications Notifier
}

type gattCharacteristic struct {
	*gattService
	characteristic *service.GattCharacteristic1
}

func GattApp(objectName string, objectPath string, localName string) *gattApp {
	a := &gattApp{}
	var err error

	a.readHandlers = make(map[string]HandleRead)
	a.writeHandlers = make(map[string]HandleWrite)
	a.notifiers = make(map[*service.GattCharacteristic1]Notifier)
	a.done = make(chan struct{})

	a.app, err = service.NewApplication(&service.ApplicationConfig{
		ObjectName: objectName,
		ObjectPath: dbus.ObjectPath(objectPath),
		LocalName:  localName,
		ReadFunc:   a.handleRead,
		WriteFunc:  a.handleWrite,
	})
	if err != nil {
		return &gattApp{
			err: errors.Errorf("Could not create app: %v", err),
		}
	}

	return a
}

func (a *gattApp) handleRead(app *service.Application, serviceUuid string, characteristicUuid string) ([]byte, error) {
	if readHandler, ok := a.readHandlers[characteristicUuid]; ok {
		return readHandler()
	}

	return nil, service.NewCallbackError(service.CallbackNotRegistered, "")
}

func (a *gattApp) handleWrite(app *service.Application, serviceUuid string, characteristicUuid string, value []byte) error {
	if writeHandler, ok := a.writeHandlers[characteristicUuid]; ok {
		return writeHandler(value)
	}

	return service.NewCallbackError(service.CallbackNotRegistered, "")
}

func (a *gattApp) Run() (*service.Application, error) {
	if a.err != nil {
		return nil, a.err
	}

	err := a.app.Run()
	if err != nil {
		return nil, errors.Errorf("Could not run app: %v", err)
	}

	for char, notifier := range a.notifiers {
		go a.handleNotifications(char, notifier)
	}

	return a.app, nil
}

func (a *gattApp) Stop() {
	close(a.done)
}

func (a *gattApp) handleNotifications(characteristic *service.GattCharacteristic1, notifier Notifier) {
	for {
		select {
		case <-a.done:
			return
		case value := <-notifier:
			err := characteristic.WriteValue(value, nil)
			if err != nil {
				fmt.Printf("Could not write: %v", err)
			}
		}
	}
}

func (a *gattApp) Service(primaryType PrimaryType, uuid string, advertised AdvertisedType) *gattService {
	if a.err != nil {
		return &gattService{gattApp: a}
	}

	svc, err := a.app.CreateService(&profile.GattService1Properties{
		Primary: bool(primaryType),
		UUID:    uuid,
	}, bool(advertised))

	if err != nil {
		a.err = errors.Errorf("Failed to create service: %v", err)
		return &gattService{gattApp: a}
	}

	err = a.app.AddService(svc)
	if err != nil {
		a.err = errors.Errorf("Failed to add service: %v", err)
		return &gattService{gattApp: a}
	}

	return &gattService{
		gattApp: a,
		service: svc,
	}
}

func (s *gattService) DeviceNameCharacteristic() *gattPendingCharacteristic {
	return s.Characteristic("2A00")
}

func (s *gattService) ManufacturerNameCharacteristic() *gattPendingCharacteristic {
	return s.Characteristic("2A29")
}

func (s *gattService) SerialNumberCharacteristic() *gattPendingCharacteristic {
	return s.Characteristic("2A25")
}

func (s *gattService) ModelNumberCharacteristic() *gattPendingCharacteristic {
	return s.Characteristic("2A24")
}

func (s *gattService) Characteristic(uuid string) *gattPendingCharacteristic {
	if s.err != nil {
		return &gattPendingCharacteristic{gattService: s}
	}

	return &gattPendingCharacteristic{
		gattService:        s,
		characteristicUuid: uuid,
	}
}

func (c *gattPendingCharacteristic) WithValue(value string) *gattPendingCharacteristic {
	if c.err != nil {
		c.err = errors.Errorf("Failed to set characteristic value: %v", c.err)
		return c
	}

	c.characteristicValue = []byte(value)
	return c
}

func (c *gattPendingCharacteristic) WithReadHandler(read HandleRead) *gattPendingCharacteristic {
	if c.err != nil {
		c.err = errors.Errorf("Failed to set characteristic read handler: %v", c.err)
		return c
	}

	c.characteristicRead = read
	return c
}

func (c *gattPendingCharacteristic) WithWriteHandler(write HandleWrite) *gattPendingCharacteristic {
	if c.err != nil {
		c.err = errors.Errorf("Failed to set characteristic write handler: %v", c.err)
		return c
	}

	c.characteristicWrite = write
	return c
}

func (c *gattPendingCharacteristic) WithNotifications(changes Notifier) *gattPendingCharacteristic {
	if c.err != nil {
		c.err = errors.Errorf("Failed to set characteristic notifications: %v", c.err)
		return c
	}

	c.characteristicNotifications = changes
	return c
}

func (c *gattPendingCharacteristic) Create() *gattCharacteristic {
	if c.err != nil {
		c.err = errors.Errorf("Failed to create characteristic: %v", c.err)
		return &gattCharacteristic{gattService: c.gattService}
	}

	var inferredFlags []string

	if c.characteristicRead != nil || c.characteristicValue != nil {
		inferredFlags = append(inferredFlags, bluez.FlagCharacteristicRead)
	}

	if c.characteristicRead != nil {
		// TODO: Mapping by characteristic UUID only makes this work for one service
		c.readHandlers[c.characteristicUuid] = c.characteristicRead
	}

	if c.characteristicWrite != nil {
		inferredFlags = append(inferredFlags, bluez.FlagCharacteristicWrite)

		// TODO: Mapping by characteristic UUID only makes this work for one service
		c.writeHandlers[c.characteristicUuid] = c.characteristicWrite
	}

	if c.characteristicNotifications != nil {
		inferredFlags = append(inferredFlags, bluez.FlagCharacteristicNotify)
	}

	characteristic, err := c.service.CreateCharacteristic(&profile.GattCharacteristic1Properties{
		UUID:  c.characteristicUuid,
		Value: c.characteristicValue,
		Flags: inferredFlags,
	})

	if err != nil {
		c.err = errors.Errorf("Failed to create characteristic: %v", err)
		return &gattCharacteristic{gattService: c.gattService}
	}

	if c.characteristicNotifications != nil {
		c.notifiers[characteristic] = c.characteristicNotifications
	}

	err = c.service.AddCharacteristic(characteristic)
	if err != nil {
		c.err = errors.Errorf("Failed to add characteristic: %v", err)
		return &gattCharacteristic{gattService: c.gattService}
	}

	return &gattCharacteristic{
		gattService:    c.gattService,
		characteristic: characteristic,
	}
}

func (c *gattCharacteristic) UserDescriptionDescriptor(value string) *gattCharacteristic {
	return c.descriptor("2901", []byte(value))
}

func (c *gattCharacteristic) PresentationDescriptor() *gattCharacteristic {
	return c.descriptor("2904", []byte{25})
}

func (c *gattCharacteristic) descriptor(uuid string, value []byte) *gattCharacteristic {
	if c.err != nil {
		return c
	}

	descriptor, err := c.characteristic.CreateDescriptor(&profile.GattDescriptor1Properties{
		UUID:  uuid,
		Value: value,
		Flags: []string{
			bluez.FlagDescriptorRead,
		},
	})

	if err != nil {
		c.err = errors.Errorf("Failed to create descriptor: %v", err)
		return c
	}

	err = c.characteristic.AddDescriptor(descriptor)
	if err != nil {
		c.err = errors.Errorf("Failed to add descriptor: %v", err)
		return c
	}

	return c
}
