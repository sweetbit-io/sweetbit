package wpa

import (
	"github.com/go-errors/errors"
	"github.com/godbus/dbus/v5"
)

type Wpa struct {
	conn *dbus.Conn
	obj  dbus.BusObject
}

func New() *Wpa {
	wpa := &Wpa{
	}

	return wpa
}

func (w *Wpa) Start() error {
	conn, err := dbus.SystemBusPrivate()
	if err != nil {
		return errors.Errorf("could not create system bus: %v", err)
	}

	err = conn.Auth(nil)
	if err != nil {
		_ = conn.Close()
		return errors.Errorf("could not authenticate bus: %v", err)
	}

	err = conn.Hello()
	if err != nil {
		_ = conn.Close()
		return errors.Errorf("could not send hello: %v", err)
	}

	obj := conn.Object("fi.w1.wpa_supplicant1", "/fi/w1/wpa_supplicant1")
	if obj == nil {
		_ = conn.Close()
		return errors.Errorf("could not create wpa_supplicant object: %v", err)
	}

	call := conn.BusObject().AddMatchSignal("fi.w1.wpa_supplicant1", "PropertiesChanged", dbus.WithMatchObjectPath(obj.Path()))
	if call.Err != nil {
		_ = conn.Close()
		return errors.Errorf("could not add signal: %v", call.Err)
	}

	w.conn = conn
	w.obj = obj

	return nil
}

func (w *Wpa) GetInterface(name string) (*Interface, error) {
	call := w.obj.Call("fi.w1.wpa_supplicant1.GetInterface", 0, name)
	if call.Err != nil {
		return nil, errors.Errorf("could not find interface %v: %v", name, call.Err)
	}

	objectPath := call.Body[0].(dbus.ObjectPath)

	ifaceObj := w.conn.Object("fi.w1.wpa_supplicant1", objectPath)

	return &Interface{
		wpa: w,
		obj: ifaceObj,
	}, nil
}

func (w *Wpa) Stop() error {
	err := w.conn.Close()
	if err != nil {
		w.conn = nil
		return errors.Errorf("could not close connection: %v", err)
	}

	w.conn = nil

	return nil
}
