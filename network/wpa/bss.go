package wpa

import (
	"encoding/hex"
	"github.com/go-errors/errors"
	"github.com/godbus/dbus/v5"
)

type BSS struct {
	obj dbus.BusObject
}

func (b *BSS) String() string {
	return string(b.obj.Path())
}

type WpaType uint8

const (
	WpaPersonal WpaType = iota
	WpaEnterprise
	WpaNone
)

type Bss struct {
	Ssid    string
	Bssid   string
	WpaType WpaType
}

func (b *BSS) GetAll() (*Bss, error) {
	call := b.obj.Call("org.freedesktop.DBus.Properties.GetAll", 0, "fi.w1.wpa_supplicant1.BSS")
	if call.Err != nil {
		return nil, errors.Errorf("could not get all properties: %v", call.Err)
	}

	props, ok := call.Body[0].(map[string]dbus.Variant)
	if !ok {
		return nil, errors.Errorf("could convert output")
	}

	bss := Bss{}

	if val, ok := props["SSID"]; ok {
		if ssid, ok := val.Value().([]byte); ok {
			bss.Ssid = string(ssid)
		} else {
			return nil, errors.Errorf("could not convert SSID to string: %v", val)
		}
	} else {
		return nil, errors.Errorf("mandatory property SSID was missing")
	}

	if val, ok := props["BSSID"]; ok {
		if bssid, ok := val.Value().([]byte); ok {
			bss.Bssid = hex.EncodeToString(bssid)
		} else {
			return nil, errors.Errorf("could not convert BSSID to string: %v", val)
		}
	} else {
		return nil, errors.Errorf("mandatory property BSSID was missing")
	}

	bss.WpaType = WpaNone

	if val, ok := props["RSN"]; ok {
		if wpa, ok := val.Value().(map[string]dbus.Variant); ok {
			if val, ok := wpa["KeyMgmt"]; ok {
				if keyMgmts, ok := val.Value().([]string); ok {
					for _, keyMgmt := range keyMgmts {
						if keyMgmt == "wpa-psk" {
							bss.WpaType = WpaPersonal
						} else if keyMgmt == "wpa-eap" {
							bss.WpaType = WpaEnterprise
						} else if keyMgmt == "wpa-none" {
							bss.WpaType = WpaNone
						} else {
							return nil, errors.Errorf("unknown wpa key management type: %s", keyMgmt)
						}
					}
				}
			}
		} else {
			return nil, errors.Errorf("could not convert WPA to dictionary: %v", val)
		}
	} else {
		return nil, errors.Errorf("mandatory property WPA was missing")
	}

	return &bss, nil
}
