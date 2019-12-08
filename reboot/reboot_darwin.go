package reboot

import (
	"github.com/go-errors/errors"
)

func Reboot() error {
	return errors.New("Reboot not supported on this system")
}

func ShutDown() error {
	return errors.New("Shut down not supported on this system")
}