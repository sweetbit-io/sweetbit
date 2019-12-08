package ble

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
	"github.com/sirupsen/logrus"
	"strings"
)

const UserDescriptionDescriptorUuid = "2901"
const PresentationFormatDescriptorUuid = "2904"

var replacer = strings.NewReplacer("-", "_")

type AppOption func(app *GattApp) error
type ServiceOption func(service *gattService) error
type CharacteristicOption func(characteristic *gattCharacteristic) error
type DescriptorOption func(descriptor *gattDescriptor) error

type GattApp struct {
	options           []AppOption
	dbusConn          *dbus.Conn
	adapterId         string
	gattManager       *gatt.GattManager1
	objManager        *api.DBusObjectManager
	adapterObjectPath dbus.ObjectPath
	appObjectPath     dbus.ObjectPath
}

type gattService struct {
	*GattApp
	serviceObjectPath dbus.ObjectPath
	dbusProperties    *api.DBusProperties
	properties        *gatt.GattService1Properties
}

type Writer func([]byte) error
type WriteHandler func([]byte) error
type ReadHandler func() ([]byte, error)

type gattCharacteristic struct {
	*gattService
	dbusProperties           *api.DBusProperties
	properties               *gatt.GattCharacteristic1Properties
	characteristicObjectPath dbus.ObjectPath
	writeHandler             WriteHandler
	readHandler              ReadHandler
}

func (c *gattCharacteristic) ReadValue(options map[string]interface{}) ([]byte, *dbus.Error) {
	if c.readHandler != nil {
		value, err := c.readHandler()
		if err != nil {
			return nil, dbus.MakeFailedError(err)
		}

		return value, nil
	}

	return c.properties.Value, nil
}

func (c *gattCharacteristic) WriteValue(value []byte, options map[string]interface{}) *dbus.Error {
	err := c.dbusProperties.Instance().Set(gatt.GattCharacteristic1Interface, "Value", dbus.MakeVariant(value))
	if err != nil {
		return dbus.MakeFailedError(errors.Errorf("unable to write: %v", err))
	}

	c.properties.Lock()
	c.properties.Value = value
	c.properties.Unlock()

	if c.writeHandler != nil {
		err := c.writeHandler(value)
		if err != nil {
			return dbus.MakeFailedError(err)
		}
	}

	return nil
}

type gattDescriptor struct {
	*gattCharacteristic
	dbusProperties       *api.DBusProperties
	properties           *gatt.GattDescriptor1Properties
	descriptorObjectPath dbus.ObjectPath
}

func (d *gattDescriptor) ReadValue(options map[string]interface{}) ([]byte, *dbus.Error) {
	return d.properties.Value, nil
}

func (d *gattDescriptor) WriteValue(value []byte, options map[string]interface{}) *dbus.Error {
	d.properties.Value = value
	return nil
}

func WithAppService(serviceUuid string, opts ...ServiceOption) AppOption {
	return func(app *GattApp) error {
		var err error

		dbusProperties, err := api.NewDBusProperties(app.dbusConn)
		if err != nil {
			logrus.Fatalf("unable to create dbus properties: %v", err)
		}

		_, err = app.dbusConn.RequestName(bluez.OrgBluezInterface, dbus.NameFlagDoNotQueue&dbus.NameFlagReplaceExisting)
		if err != nil {
			logrus.Fatalf("unable to request name: %v", err)
		}

		serviceObjectPath := dbus.ObjectPath(fmt.Sprintf("%s/%s", app.appObjectPath, replacer.Replace(serviceUuid)))

		properties := &gatt.GattService1Properties{
			Primary:         true,
			UUID:            serviceUuid,
			Characteristics: []dbus.ObjectPath{},
		}

		service := &gattService{
			GattApp:           app,
			serviceObjectPath: serviceObjectPath,
			dbusProperties:    dbusProperties,
			properties:        properties,
		}

		for _, opt := range opts {
			err := opt(service)
			if err != nil {
				return errors.Errorf("option execution failed %v", err)
			}
		}

		err = service.dbusConn.Export(service, serviceObjectPath, gatt.GattService1Interface)
		if err != nil {
			logrus.Fatalf("unable to export service: %v", err)
		}

		err = dbusProperties.AddProperties(gatt.GattService1Interface, properties)
		if err != nil {
			logrus.Fatalf("unable to create add properties: %v", err)
		}

		dbusProperties.Expose(serviceObjectPath)

		err = app.objManager.AddObject(serviceObjectPath, map[string]bluez.Properties{
			gatt.GattService1Interface: properties,
		})
		if err != nil {
			logrus.Fatalf("unable to add service: %v", err)
		}

		return nil
	}
}

func WithServiceCharacteristic(characteristicUuid string, opts ...CharacteristicOption) ServiceOption {
	return func(service *gattService) error {
		var err error

		dbusProperties, err := api.NewDBusProperties(service.dbusConn)
		if err != nil {
			logrus.Fatalf("unable to create dbus properties: %v", err)
		}

		characteristicObjectPath := dbus.ObjectPath(fmt.Sprintf("%s/%s", service.serviceObjectPath, replacer.Replace(characteristicUuid)))

		properties := &gatt.GattCharacteristic1Properties{
			UUID:        characteristicUuid,
			Service:     service.serviceObjectPath,
			Flags:       []string{},
			Descriptors: []dbus.ObjectPath{},
		}

		characteristic := &gattCharacteristic{
			gattService:              service,
			dbusProperties:           dbusProperties,
			properties:               properties,
			characteristicObjectPath: characteristicObjectPath,
		}

		for _, opt := range opts {
			err := opt(characteristic)
			if err != nil {
				return errors.Errorf("option execution failed %v", err)
			}
		}

		err = service.dbusConn.Export(characteristic, characteristicObjectPath, gatt.GattCharacteristic1Interface)
		if err != nil {
			logrus.Fatalf("unable to export char service: %v", err)
		}

		err = dbusProperties.AddProperties(gatt.GattCharacteristic1Interface, properties)
		if err != nil {
			logrus.Fatalf("unable to create add properties: %v", err)
		}

		dbusProperties.Expose(characteristicObjectPath)

		err = service.objManager.AddObject(characteristicObjectPath, map[string]bluez.Properties{
			gatt.GattCharacteristic1Interface: properties,
		})
		if err != nil {
			logrus.Fatalf("unable to add char: %v", err)
		}

		// add newly created characteristic to parent service
		service.properties.Lock()
		service.properties.Characteristics = append(service.properties.Characteristics, characteristicObjectPath)
		service.properties.Unlock()

		return nil
	}
}

func hasFlag(flag string, flags []string) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}

	return false
}

func WithCharacteristicWriteHandler(handler WriteHandler) CharacteristicOption {
	return func(characteristic *gattCharacteristic) error {
		characteristic.writeHandler = handler

		if !hasFlag(gatt.FlagCharacteristicWrite, characteristic.properties.Flags) {
			characteristic.properties.Lock()
			characteristic.properties.Flags = append(characteristic.properties.Flags, gatt.FlagCharacteristicWrite)
			characteristic.properties.Unlock()
		}

		return nil
	}
}

func WithCharacteristicReadHandler(handler ReadHandler) CharacteristicOption {
	return func(characteristic *gattCharacteristic) error {
		characteristic.readHandler = handler

		if !hasFlag(gatt.FlagCharacteristicRead, characteristic.properties.Flags) {
			characteristic.properties.Lock()
			characteristic.properties.Flags = append(characteristic.properties.Flags, gatt.FlagCharacteristicRead)
			characteristic.properties.Unlock()
		}

		return nil
	}
}

func WithCharacteristicWriter(notify *Writer) CharacteristicOption {
	return func(characteristic *gattCharacteristic) error {

		if !hasFlag(gatt.FlagCharacteristicNotify, characteristic.properties.Flags) {
			characteristic.properties.Lock()
			characteristic.properties.Flags = append(characteristic.properties.Flags, gatt.FlagCharacteristicNotify)
			characteristic.properties.Unlock()
		}

		*notify = func(value []byte) error {
			err := characteristic.dbusProperties.Instance().Set(gatt.GattCharacteristic1Interface, "Value", dbus.MakeVariant(value))
			if err != nil {
				return errors.Errorf("unable to write: %v", err)
			}

			characteristic.properties.Lock()
			characteristic.properties.Value = value
			characteristic.properties.Unlock()

			return nil
		}

		return nil
	}
}

func WithDescriptorValue(value []byte) DescriptorOption {
	return func(descriptor *gattDescriptor) error {
		if !hasFlag(gatt.FlagDescriptorRead, descriptor.properties.Flags) {
			descriptor.properties.Lock()
			descriptor.properties.Flags = append(descriptor.properties.Flags, gatt.FlagDescriptorRead)
			descriptor.properties.Unlock()
		}

		descriptor.properties.Lock()
		descriptor.properties.Value = value
		descriptor.properties.Unlock()

		return nil
	}
}

func WithCharacteristicUserDescriptionDescriptor(description string) CharacteristicOption {
	return WithCharacteristicDescriptor(UserDescriptionDescriptorUuid, WithDescriptorValue([]byte(description)))
}

func WithCharacteristicPresentationFormatDescriptor() CharacteristicOption {
	type packet struct {
		Format      byte
		Exponent    int8
		Unit        uint16
		Namespace   byte
		Description [2]byte
	}

	value := packet{
		Format:      25,
		Exponent:    0,
		Unit:        0x2700,
		Namespace:   1,
		Description: [2]byte{0, 0},
	}

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, value)

	return WithCharacteristicDescriptor(PresentationFormatDescriptorUuid, WithDescriptorValue(buf.Bytes()))
}

func WithCharacteristicDescriptor(descriptorUuid string, opts ...DescriptorOption) CharacteristicOption {
	return func(characteristic *gattCharacteristic) error {
		var err error

		dbusProperties, err := api.NewDBusProperties(characteristic.dbusConn)
		if err != nil {
			logrus.Fatalf("unable to create dbus properties: %v", err)
		}

		descriptorObjectPath := dbus.ObjectPath(fmt.Sprintf("%s/%s", characteristic.characteristicObjectPath, replacer.Replace(descriptorUuid)))

		properties := &gatt.GattDescriptor1Properties{
			UUID:           descriptorUuid,
			Characteristic: characteristic.characteristicObjectPath,
			Flags:          []string{},
		}

		descriptor := &gattDescriptor{
			gattCharacteristic:   characteristic,
			dbusProperties:       dbusProperties,
			properties:           properties,
			descriptorObjectPath: descriptorObjectPath,
		}

		for _, opt := range opts {
			err := opt(descriptor)
			if err != nil {
				return errors.Errorf("option execution failed: %v", err)
			}
		}

		err = characteristic.dbusConn.Export(descriptor, descriptorObjectPath, gatt.GattDescriptor1Interface)
		if err != nil {
			logrus.Fatalf("unable to export descriptor: %v", err)
		}

		err = dbusProperties.AddProperties(gatt.GattDescriptor1Interface, properties)
		if err != nil {
			logrus.Fatalf("unable to create add descriptors: %v", err)
		}

		dbusProperties.Expose(descriptorObjectPath)

		err = characteristic.objManager.AddObject(descriptorObjectPath, map[string]bluez.Properties{
			gatt.GattDescriptor1Interface: properties,
		})
		if err != nil {
			logrus.Fatalf("unable to add char %v: %v", descriptorObjectPath, err)
		}

		// add newly created descriptor to parent characteristic
		characteristic.properties.Lock()
		characteristic.properties.Descriptors = append(characteristic.properties.Descriptors, descriptorObjectPath)
		characteristic.properties.Unlock()

		return nil
	}
}

func NewGattApp(adapterId string, appPath string, opts ...AppOption) *GattApp {
	adapterObjectPath := dbus.ObjectPath(bluez.OrgBluezPath + "/" + "hci0")
	appObjectPath := dbus.ObjectPath(appPath)

	app := &GattApp{
		options:           opts,
		adapterId:         adapterId,
		adapterObjectPath: adapterObjectPath,
		appObjectPath:     appObjectPath,
	}

	return app
}

func (a *GattApp) Start() error {
	var err error

	a.dbusConn, err = dbus.SystemBus()
	if err != nil {
		logrus.Fatalf("unable to get connection: %v", err)
	}

	a.objManager, err = api.NewDBusObjectManager(a.dbusConn)
	if err != nil {
		logrus.Fatalf("unable to create object manager: %v", err)
	}

	err = a.dbusConn.Export(a.objManager, a.appObjectPath, bluez.ObjectManagerInterface)
	if err != nil {
		logrus.Fatalf("unable to export object manager: %v", err)
	}

	for _, opt := range a.options {
		err := opt(a)
		if err != nil {
			return errors.Errorf("option execution failed: %v", err)
		}
	}

	a.gattManager, err = gatt.NewGattManager1(a.adapterObjectPath)
	if err != nil {
		logrus.Fatalf("unable to create gatt manager: %v", err)
	}

	err = a.gattManager.RegisterApplication(a.appObjectPath, map[string]interface{}{})
	if err != nil {
		logrus.Fatalf("unable to register app %v: %v", a.appObjectPath, err)
	}

	return nil
}

func (a *GattApp) Stop() error {
	var err error

	err = a.gattManager.UnregisterApplication(a.appObjectPath)
	if err != nil {
		logrus.Fatalf("unable to unregister app %v: %v", a.appObjectPath, err)
	}

	err = a.dbusConn.Close()
	if err != nil {
		return errors.Errorf("unable to close dbus connection: %v", err)
	}

	return nil
}
