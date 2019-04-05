package pairing

import (
	"github.com/the-lightning-land/sweetd/ap"
	"github.com/the-lightning-land/sweetd/dispenser"
)

type Config struct {
	Logger      Logger
	AdapterId   string
	AccessPoint ap.Ap
	Dispenser   *dispenser.Dispenser
}
