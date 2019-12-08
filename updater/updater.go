package updater

import "time"

type State = string

const StateStarted State = "started"
const StateCancelled State = "cancelled"
const StateFailed State = "failed"
const StateInstalled State = "installed"
const StateRejected State = "rejected"
const StateCompleted State = "completed"

type Update struct {
	Id           string
	Started      time.Time
	Url          string
	State        State
	Progress     uint8
	ShouldReboot bool
	ShouldCommit bool
}

type Updater interface {
	GetVersion() (string, error)
	StartUpdate(url string) (*Update, error)
	GetUpdate(id string) (*Update, error)
	GetCurrentUpdate() (*Update, error)
	CancelUpdate(id string) (*Update, error)
	SubscribeUpdate(id string) (*UpdateClient, error)
	unsubscribeUpdate(client *UpdateClient)
	CommitUpdate(id string) (*Update, error)
	RejectUpdate(id string) (*Update, error)
}
