package node

type InvoicesClient struct {
	Invoices   chan *Invoice
	Id         uint32
	cancelChan chan struct{}
	node       Node
}

func (c *InvoicesClient) Cancel() error {
	return c.node.unsubscribeInvoices(c)
}
