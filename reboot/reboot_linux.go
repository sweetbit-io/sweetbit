package reboot

import (
	"github.com/go-errors/errors"
	"golang.org/x/sys/unix"
)

func Reboot() error {
	unix.Sync()

	err := unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART)
	if err != nil {
		return errors.Errorf("Could not reboot: %v", err)
	}

	return nil
}