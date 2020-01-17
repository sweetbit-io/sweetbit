package lightning

type Invoice struct {
	RHash          string
	PaymentRequest string
	Settled        bool
	MSat           int64
	Memo           string
}

type InvoicesClient struct {
	Invoices   chan *Invoice
	Id         uint32
	cancelChan chan struct{}
	node       Node
}

func (c *InvoicesClient) Cancel() {
	c.node.unsubscribeInvoices(c)
}

type StatusClient struct {
	Status     chan Status
	Id         uint32
	cancelChan chan struct{}
	node       Node
}

func (c *StatusClient) Cancel() {
	c.node.unsubscribeStatus(c)
}

type InvoiceRequest struct {
	MSat int64
	Memo string
}

type Status int

const (
	StatusStopped Status = iota
	StatusUninitialized
	StatusLocked
	StatusStarted
	StatusFailed
)

type Node interface {
	Start() error
	Stop() error
	GetInvoice(rHash string) (*Invoice, error)
	AddInvoice(request *InvoiceRequest) (*Invoice, error)
	SubscribeInvoices() (*InvoicesClient, error)
	SubscribeStatus() *StatusClient
	unsubscribeInvoices(client *InvoicesClient)
	unsubscribeStatus(client *StatusClient)
	Status() Status
	Unlock(password string) error
	Init(password string, mnemonics []string) error
	GenerateSeed() ([]string, error)
}
