package network

// check MockNetworks compliance to its interface during compile time
var _ Network = (*MockNetwork)(nil)

type MockNetwork struct {
}

func NewMockNetwork() *MockNetwork {
	return &MockNetwork{}
}

func (m *MockNetwork) Start() error {
	return nil
}

func (m *MockNetwork) Status() *Status {
	return &Status{
		connected: false,
	}
}

func (m *MockNetwork) Stop() error {
	return nil
}

func (m *MockNetwork) Connect(connection Connection) error {
	return nil
}

func (m *MockNetwork) Scan() (*ScanClient, error) {
	return nil, nil
}

func (m *MockNetwork) Subscribe() *Client {
	return &Client{
		Updates:    make(chan *Connectivity),
		Id:         0,
		cancelChan: make(chan struct{}),
		network:    m,
	}
}

func (m *MockNetwork) deleteClient(id uint32) {
}
