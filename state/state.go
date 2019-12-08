package state

type State string

const (
	StateStarting State = "starting"
	StateStarted        = "running"
	StateStopping       = "stopping"
	StateStopped        = "stopped"
)

func String(state State) string {
	return string(state)
}
