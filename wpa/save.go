package wpa

import (
	"os/exec"
	"github.com/pkg/errors"
	"strings"
)

func Save(iface string) error {
	result, err := exec.Command("wpa_cli", "-i", iface, "save_config").Output()
	if err != nil {
		return errors.Errorf("Command: %s", err.Error())
	}

	status := strings.TrimSpace(string(result))

	return errors.Errorf("Got %s", status)
}