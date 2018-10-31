package reboot

import (
	"github.com/go-errors/errors"
)

func Reboot() error {
	return errors.New("Reboot not supported on this system")
}