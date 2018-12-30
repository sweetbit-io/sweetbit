package wpa

import (
	"github.com/go-errors/errors"
	"os/exec"
	"strings"
)

type ConfiguredNetwork struct {
	Id   NetworkId
	Ssid string
}

// Status returns the WPA wireless status.
func ListConfiguredNetworks(iface string) ([]*ConfiguredNetwork, error) {
	configuredNetworkListOut, err := exec.Command("wpa_cli", "-i", iface, "list_networks").Output()
	if err != nil {
		return nil, errors.Errorf("Command: %s", err.Error())
	}

	configuredNetworks, err := parseListNetworksOutput(&configuredNetworkListOut)
	if err != nil {
		return nil, errors.Errorf("Failed parsing network list: %v", err)
	}

	return configuredNetworks, nil
}

func parseListNetworksOutput(output *[]byte) ([]*ConfiguredNetwork, error) {
	configuredNetworks := make([]*ConfiguredNetwork, 0)

	networkListOutArr := strings.Split(string(*output), "\n")
	for _, netRecord := range networkListOutArr[1:] {
		fields := strings.Split(netRecord, "\t")

		if len(fields) >= 2 {
			configuredNetworks = append(configuredNetworks, &ConfiguredNetwork{
				Id:   NetworkId(fields[0]),
				Ssid: fields[1],
			})
		}
	}

	return configuredNetworks, nil
}
