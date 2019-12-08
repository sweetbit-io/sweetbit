package dispenser

import (
	"github.com/the-lightning-land/sweetd/network"
	"github.com/the-lightning-land/sweetd/pairing"
)

// PairingAdapter adapts the dispenser to the pairing controller api
type PairingAdapter struct {
	Dispenser *Dispenser
	Pairing   pairing.Controller
}

func NewPairingAdapter() *PairingAdapter {
	return &PairingAdapter{}
}

func (a *PairingAdapter) ScanWifi() (*network.ScanClient, error) {
	return a.Dispenser.ScanWifi()
}

func (a *PairingAdapter) ConnectWifi(connection network.Connection) error {
	return a.Dispenser.ConnectToWifi(connection)
}

func (a *PairingAdapter) GetApiOnionID() string {
	return a.Dispenser.GetApiOnionID()
}

func (a *PairingAdapter) GetName() string {
	return a.Dispenser.GetName()
}
