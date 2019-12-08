package reboot

import (
	"github.com/go-errors/errors"
	"golang.org/x/sys/unix"
)

func Reboot() error {
	err := unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART)
	if err != nil {
		return errors.Errorf("Could not reboot: %v", err)
	}

	return nil
}

func ShutDown() error {
	err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF)
	if err != nil {
		return errors.Errorf("Could not shut down: %v", err)
	}

	return nil
}
