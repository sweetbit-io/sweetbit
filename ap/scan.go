package ap

import (
	"os/exec"
	"strings"
	"github.com/go-errors/errors"
)

var errBusy = errors.New("Device is busy")

func scan(iface string) error {
	result, err := exec.Command("wpa_cli", "-i", iface, "scan").Output()
	if err != nil {
		return errors.Errorf("Command: %v", err)
	}

	resultClean := strings.TrimSpace(string(result))
	if resultClean == "FAIL-BUSY" {
		return errBusy
	} else if resultClean != "OK" {
		return errors.Errorf("Got %s instead of OK", resultClean)
	}

	return nil
}

// Results returns an array of WpaNetwork data structures.
func results(iface string) ([]*Network, error) {
	wpaNetworks := make([]*Network, 0)

	networkListOut, err := exec.Command("wpa_cli", "-i", iface, "scan_results").Output()
	if err != nil {
		return wpaNetworks, errors.Errorf("Command: %v", err)
	}

	networkListOutArr := strings.Split(string(networkListOut), "\n")
	for _, netRecord := range networkListOutArr[1:] {
		if strings.Contains(netRecord, "[P2P]") {
			continue
		}

		fields := strings.Fields(netRecord)

		if len(fields) > 4 {
			ssid := strings.Join(fields[4:], " ")
			wpaNetworks = append(wpaNetworks, &Network{
				Bssid:       fields[0],
				Frequency:   fields[1],
				SignalLevel: fields[2],
				Flags:       fields[3],
				Ssid:        ssid,
			})
		}
	}

	return wpaNetworks, nil
}
