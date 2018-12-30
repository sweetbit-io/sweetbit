package wpa

import (
	"os/exec"
	"strings"
	"github.com/pkg/errors"
)

type NetworkId string

type AddNetworkKey string

const (
	Ssid AddNetworkKey = "ssid"
	Psk  AddNetworkKey = "psk"
)

func AddNetwork(iface string) (NetworkId, error) {
	result, err := exec.Command("wpa_cli", "-i", iface, "add_network").Output()
	if err != nil {
		return "", errors.Errorf("Command: %s", err.Error())
	}

	net := strings.TrimSpace(string(result))

	return NetworkId(net), nil
}

func RemoveNetwork(iface string, net NetworkId) (error) {
	_, err := exec.Command("wpa_cli", "-i", iface, "remove_network", string(net)).Output()
	if err != nil {
		return errors.Errorf("Command: %s", err.Error())
	}

	return nil
}

func SetNetwork(iface string, net NetworkId, key AddNetworkKey, value string) error {
	result, err := exec.Command("wpa_cli", "-i", iface, "set_network", string(net), string(key), "\""+value+"\"").Output()
	if err != nil {
		return errors.Errorf("Command: %s", err.Error())
	}

	status := strings.TrimSpace(string(result))

	if status != "OK" {
		return errors.Errorf("Got %s", status)
	}

	return nil
}

func EnableNetwork(iface string, net NetworkId) error {
	result, err := exec.Command("wpa_cli", "-i", iface, "enable_network", string(net)).Output()
	if err != nil {
		return errors.Errorf("Command: %s", err.Error())
	}

	status := strings.TrimSpace(string(result))

	if status != "OK" {
		return errors.Errorf("Got %s", status)
	}

	return nil
}
