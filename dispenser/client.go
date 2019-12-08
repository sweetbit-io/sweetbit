package dispenser

type DispenseClient struct {
	Dispenser  chan *Dispenser
	Dispenses  chan DispenseState
	Id         uint32
	cancelChan chan struct{}
	dispenser  *Dispenser
}

func (c *DispenseClient) Cancel() {
	delete(c.dispenser.dispenseClients, c.Id)
	close(c.cancelChan)
}
