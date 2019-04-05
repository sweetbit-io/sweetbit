package ap

import "os/exec"

// renewUdhcp renews UDHCP lease
func renewUdhcp() error {
	cmd := exec.Command("killall", "-SIGUSR1", "udhcpc")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}
