package ap

import (
	"os/exec"
	"strings"
	"github.com/pkg/errors"
)

type networkId string

type addNetworkKey string

const (
	ssidKey addNetworkKey = "ssid"
	pskKey  addNetworkKey = "psk"
)

func addNetwork(iface string) (networkId, error) {
	result, err := exec.Command("wpa_cli", "-i", iface, "add_network").Output()
	if err != nil {
		return "", errors.Errorf("Command: %s", err.Error())
	}

	net := strings.TrimSpace(string(result))

	return networkId(net), nil
}

func removeNetwork(iface string, net networkId) (error) {
	_, err := exec.Command("wpa_cli", "-i", iface, "remove_network", string(net)).Output()
	if err != nil {
		return errors.Errorf("Command: %s", err.Error())
	}

	return nil
}

func setNetwork(iface string, net networkId, key addNetworkKey, value string) error {
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

func enableNetwork(iface string, net networkId) error {
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
