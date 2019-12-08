package network

type Connectivity struct {
	Connected bool
	Ip        string
	Ssid      string
}

type Client struct {
	Updates    chan *Connectivity
	Id         uint32
	cancelChan chan struct{}
	network    Network
}

func (c *Client) Cancel() {
	c.network.deleteClient(c.Id)
	close(c.cancelChan)
	c.cancelChan = nil
}
