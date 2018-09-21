package main

import "os/exec"

// removeApInterface removes the AP interface.
func removeApInterface() error {
	cmd := exec.Command("iw", "dev", "uap0", "del")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

// configureApInterface configured the AP interface.
func configureApInterface(ip string) error {
	cmd := exec.Command("ifconfig", "uap0", ip)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

// upApInterface ups the AP Interface.
func upApInterface() error {
	cmd := exec.Command("ifconfig", "uap0", "up")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

// addApInterface adds the AP interface.
func addApInterface() error {
	cmd := exec.Command("iw", "phy", "phy0", "interface", "add", "uap0", "type", "__ap")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}