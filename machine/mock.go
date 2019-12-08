package machine

import (
	"net/http"
)

type MockMachine struct {
	listen            string
	touchesClients    map[uint32]*TouchesClient
	nextTouchesClient nextTouchesClient
}

// Compile time check for protocol compatibility
var _ Machine = (*MockMachine)(nil)

func NewMockMachine(listen string) *MockMachine {
	return &MockMachine{
		listen:            listen,
		touchesClients:    make(map[uint32]*TouchesClient),
		nextTouchesClient: nextTouchesClient{id: 0},
	}
}

func (m *MockMachine) Start() error {
	http.HandleFunc("/touch/on", func(w http.ResponseWriter, r *http.Request) {
		m.notifyTouchesClients(true)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/touch/off", func(w http.ResponseWriter, r *http.Request) {
		m.notifyTouchesClients(false)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("OK"))
	})

	go http.ListenAndServe(m.listen, nil)

	return nil
}

func (m *MockMachine) Stop() error {
	// nothing
	return nil
}

func (m *MockMachine) ToggleMotor(on bool) {
	// nothing
}

func (m *MockMachine) ToggleBuzzer(on bool) {
	// nothing
}

func (m *MockMachine) DiagnosticNoise() {
	// nothing
}

func (m *MockMachine) SubscribeTouches() *TouchesClient {
	client := &TouchesClient{
		Touches:    make(chan bool),
		cancelChan: make(chan struct{}),
		machine:    m,
	}

	m.nextTouchesClient.Lock()
	client.Id = m.nextTouchesClient.id
	m.nextTouchesClient.id++
	m.nextTouchesClient.Unlock()

	m.touchesClients[client.Id] = client

	return client
}

func (m *MockMachine) notifyTouchesClients(touch bool) {
	for _, client := range m.touchesClients {
		client.Touches <- touch
	}
}

func (m *MockMachine) unsubscribeTouches(client *TouchesClient) {
	delete(m.touchesClients, client.Id)
	close(client.cancelChan)
}
