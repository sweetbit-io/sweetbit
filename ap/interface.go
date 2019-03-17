package ap

import (
	"os/exec"
	"strconv"
	"strings"
)

type interfaceInfo struct {
	channel int
}

// removeApInterface removes the AP interface.
func removeApInterface(iface string) error {
	cmd := exec.Command("iw", "dev", iface, "del")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

// configureApInterface configured the AP interface.
func configureApInterface(ip string, iface string) error {
	cmd := exec.Command("ip", "addr", "add", ip, "dev", iface)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

// upApInterface ups the AP Interface.
func upApInterface(iface string) error {
	cmd := exec.Command("ip", "link", "set", iface, "up")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

// downApInterface downs the AP Interface.
func downApInterface(iface string) error {
	cmd := exec.Command("ip", "link", "set", iface, "down")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

// addApInterface adds the AP interface.
func addApInterface(iface string) error {
	cmd := exec.Command("iw", "phy", "phy0", "interface", "add", iface, "type", "__ap")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

// getApInterfaceInfo returns interface information
func getApInterfaceInfo(iface string) (*interfaceInfo, error) {
	infoOut, err := exec.Command("iw", "dev", iface, "info").Output()
	if err != nil {
		return nil, err
	}

	info, err := parseApInterfaceInfoOutput(&infoOut)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func parseApInterfaceInfoOutput(output *[]byte) (*interfaceInfo, error) {
	interfaceInfo := &interfaceInfo{}

	infoOutArr := strings.Split(string(*output), "\n")

	// Skip first line which contains the interface name
	for _, infoRecord := range infoOutArr[1:] {
		// Skip all lines that don't start with a tab character
		if len(infoRecord) == 0 || infoRecord[0] != '\t' {
			continue
		}

		// Strip the leading tab character
		line := infoRecord[1:]
		fields := strings.Fields(line)

		if len(fields) > 0 && fields[0] == "channel" {
			if i, err := strconv.Atoi(fields[1]); err == nil {
				interfaceInfo.channel = i
			}
		}
	}

	return interfaceInfo, nil
}
