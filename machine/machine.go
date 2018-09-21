package machine

type Machine interface {
	Start() error
	Stop()
	TouchEvents() <-chan bool
	ToggleMotor(on bool)
	ToggleBuzzer(on bool)
}