package machine

import (
	"net/http"
)

type MockMachine struct {
	// Listeners
	listen string
	// Internal sending channel for touch events
	touchEvents chan bool
}

func NewMockMachine(listen string) *MockMachine {
	touchEvents := make(chan bool)

	return &MockMachine{
		listen:      listen,
		touchEvents: touchEvents,
	}
}

func (m *MockMachine) Start() error {
	http.HandleFunc("/touch/on", func(w http.ResponseWriter, r *http.Request) {
		m.touchEvents <- true
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/touch/off", func(w http.ResponseWriter, r *http.Request) {
		m.touchEvents <- false
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("OK"))
	})

	go http.ListenAndServe(m.listen, nil)

	return nil
}

func (m *MockMachine) Stop() {
	// nothing
}

func (m *MockMachine) TouchEvents() <-chan bool {
	return m.touchEvents
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