package lightning

type Invoice struct {
	RHash          string
	PaymentRequest string
	Settled        bool
	MSat           int64
	Memo           string
}

type InvoiceRequest struct {
	MSat int64
	Memo string
}

type Node interface {
	Start() error
	Stop() error
	GetInvoice(rHash string) (*Invoice, error)
	AddInvoice(request *InvoiceRequest) (*Invoice, error)
	SubscribeInvoices() (*InvoicesClient, error)
	unsubscribeInvoices(client *InvoicesClient)
}
