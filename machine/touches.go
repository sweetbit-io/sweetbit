package machine

import "sync"

type nextTouchesClient struct {
	sync.Mutex
	id uint32
}

type TouchesClient struct {
	Touches    chan bool
	Id         uint32
	cancelChan chan struct{}
	machine    Machine
}

func (c *TouchesClient) Cancel() {
	c.machine.unsubscribeTouches(c)
}
