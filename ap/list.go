package ap

import (
	"github.com/go-errors/errors"
	"os/exec"
	"strings"
)

type configuredNetwork struct {
	id   networkId
	ssid string
}

// Status returns the WPA wireless status.
func listConfiguredNetworks(iface string) ([]*configuredNetwork, error) {
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

func parseListNetworksOutput(output *[]byte) ([]*configuredNetwork, error) {
	configuredNetworks := make([]*configuredNetwork, 0)

	networkListOutArr := strings.Split(string(*output), "\n")
	for _, netRecord := range networkListOutArr[1:] {
		fields := strings.Split(netRecord, "\t")

		if len(fields) >= 2 {
			configuredNetworks = append(configuredNetworks, &configuredNetwork{
				id:   networkId(fields[0]),
				ssid: fields[1],
			})
		}
	}

	return configuredNetworks, nil
}
