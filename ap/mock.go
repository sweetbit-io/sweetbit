package ap

type MockAp struct {
}

// Ensure we implement the Ap interface with this compile-time check
var _ Ap = (*MockAp)(nil)

func NewMockAp() *MockAp {
	return &MockAp{}
}

func (a *MockAp) Start() error {
	return nil
}

func (a *MockAp) StartHotspot() error {
	return nil
}

func (a *MockAp) ListWifiNetworks() ([]*Network, error) {
	empty := make([]*Network, 0)
	return empty, nil
}

func (a *MockAp) ConnectWifi(ssid string, psk string) error {
	return nil
}

func (a *MockAp) GetConnectionStatus() (*ConnectionStatus, error) {
	return &ConnectionStatus{
		Ssid:  "mock",
		Ip:    "192.168.1.42",
		State: "mocked",
	}, nil
}

func (a *MockAp) Stop() error {
	return nil
}
