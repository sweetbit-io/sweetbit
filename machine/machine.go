package machine

type Machine interface {
	Start() error
	Stop() error
	ToggleMotor(on bool)
	ToggleBuzzer(on bool)
	DiagnosticNoise()
	SubscribeTouches() *TouchesClient
	unsubscribeTouches(client *TouchesClient)
}