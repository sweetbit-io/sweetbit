package ap

type ApUpdate struct {
	Connected bool
	Ip        string
	Ssid      string
}

type ApClient struct {
	Updates    chan *ApUpdate
	Id         uint32
	cancelChan chan struct{}
	ap         Ap
}

func (c *ApClient) Cancel() {
	c.ap.deleteApClient(c.Id)

	close(c.cancelChan)
}
