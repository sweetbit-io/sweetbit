package wpa

import (
	"github.com/go-errors/errors"
	"github.com/godbus/dbus/v5"
)

type Network struct {
	wpa *Wpa
	obj dbus.BusObject
}

func (n *Network) String() string {
	return string(n.obj.Path())
}

func (n *Network) Enable() error {
	err := n.obj.SetProperty("fi.w1.wpa_supplicant1.Network.Enabled", dbus.MakeVariant(true))
	if err != nil {
		return errors.Errorf("unable to set value: %v", err)
	}

	return nil
}
