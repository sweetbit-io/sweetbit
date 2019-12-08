package pairing

type Controller interface {
	Start() error
	Stop() error
	Advertise() error
}
