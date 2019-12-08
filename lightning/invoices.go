package lightning

type InvoicesClient struct {
	Invoices   chan *Invoice
	Id         uint32
	cancelChan chan struct{}
	node       Node
}

func (c *InvoicesClient) Cancel() {
	c.node.unsubscribeInvoices(c)
}
