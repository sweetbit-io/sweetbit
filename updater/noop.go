package updater

import "errors"

type NoopUpdater struct {
}

// Compile time check for protocol compatibility
var _ Updater = (*NoopUpdater)(nil)

func NewNoopUpdater() *NoopUpdater {
	return &NoopUpdater{}
}

func (n *NoopUpdater) GetVersion() (string, error) {
	return "", errors.New("no updater available")
}

func (n *NoopUpdater) StartUpdate(url string) (*Update, error) {
	return nil, errors.New("no updater available")
}

func (n *NoopUpdater) GetUpdate(url string) (*Update, error) {
	return nil, errors.New("no updater available")
}

func (n *NoopUpdater) GetCurrentUpdate() (*Update, error) {
	return nil, nil
}

func (n *NoopUpdater) CancelUpdate(id string) (*Update, error) {
	return nil, errors.New("no updater available")
}

func (n *NoopUpdater) SubscribeUpdate(id string) (*UpdateClient, error) {
	return nil, errors.New("no updater available")
}

func (n *NoopUpdater) unsubscribeUpdate(client *UpdateClient) {
}

func (m *NoopUpdater) CommitUpdate(id string) (*Update, error) {
	return nil, errors.New("no updater available")
}

func (m *NoopUpdater) RejectUpdate(id string) (*Update, error) {
	return nil, errors.New("no updater available")
}
