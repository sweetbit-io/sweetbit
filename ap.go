package main

import "os/exec"

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