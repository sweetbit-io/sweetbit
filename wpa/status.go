package wpa

import (
	"os/exec"
	"bytes"
	"github.com/go-errors/errors"
)

type Status struct {
	Ssid  string
	State string
	Ip    string
}

// Status returns the WPA wireless status.
func GetStatus(iface string) (*Status, error) {
	result, err := exec.Command("wpa_cli", "-i", iface, "status").Output()
	if err != nil {
		return nil, errors.Errorf("Command: %s", err.Error())
	}

	lines := bytes.Split(result, []byte("\n"))
	status := Status{}

	for _, line := range lines {
		key, value, err := kvp(line)
		if err != nil {
			continue
		}

		switch key {
		case "ssid":
			status.Ssid = value
		case "wpa_state":
			status.State = value
		case "ip_address":
			status.Ip = value
		}
	}

	return &status, nil
}

func kvp(in []byte) (string, string, error) {
	s := bytes.Split(in, []byte("="))
	if len(s) != 2 {
		return "", "", errors.New("Not in format key=value")
	}
	return string(s[0]), string(s[1]), nil
}
