package updater

type Client struct {
	Progress   chan bool
	Id         uint32
	cancelChan chan struct{}
	updater    Updater
}

func (c *Client) Cancel() error {
	return c.updater.UnsubscribeUpdate(c)
}
