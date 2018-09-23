package sysid

import (
	"github.com/pkg/errors"
	"regexp"
	"os/exec"
)

func GetId() (string, error) {
	rSerial := regexp.MustCompile(`(?m)Serial\s*:\s*(\w+)`)

	res, err := exec.Command("cat", "/proc/cpuinfo").Output()

	if err != nil {
		return "", err
	}

	matches := rSerial.FindSubmatch(res)

	if len(matches) < 1 {
		return "", errors.New("Could not identify serial number.")
	}

	serial := string(matches[1])

	return serial, nil
}
