package ap

type MockAp struct {
}

// Ensure we implement the Ap interface with this compile-time check
var _ Ap = (*MockAp)(nil)

func NewMockAp() *MockAp {
	return &MockAp{}
}

func (ap *MockAp) Start() error {
	// Do nothing
	return nil
}

func (ap *MockAp) ConnectWifi(ssid string, psk string) error {
	// Do nothing
	return nil
}

func (ap *MockAp) Stop() error {
	// Do nothing
	return nil
}
