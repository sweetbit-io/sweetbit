package dispenser

import "github.com/the-lightning-land/sweetd/updater"

type Update struct {
}

func (d *Dispenser) StartUpdate(url string) (*updater.Update, error) {
	return d.updater.StartUpdate(url)
}

func (d *Dispenser) GetUpdate(id string) (*updater.Update, error) {
	return d.updater.GetUpdate(id)
}

func (d *Dispenser) GetCurrentUpdate() (*updater.Update, error) {
	return d.updater.GetCurrentUpdate()
}

func (d *Dispenser) CancelUpdate(id string) (*updater.Update, error) {
	return d.updater.CancelUpdate(id)
}

func (d *Dispenser) SubscribeUpdate(id string) (*updater.UpdateClient, error) {
	return d.updater.SubscribeUpdate(id)
}

func (d *Dispenser) CommitUpdate(id string) (*updater.Update, error) {
	return d.updater.CommitUpdate(id)
}

func (d *Dispenser) RejectUpdate(id string) (*updater.Update, error) {
	return d.updater.CommitUpdate(id)
}
