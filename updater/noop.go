package updater

import "errors"

type NoopUpdater struct {
}

// Compile time check for protocol compatibility
var _ Updater = (*NoopUpdater)(nil)

func NewNoopUpdater() *NoopUpdater {
	return &NoopUpdater{}
}

func (n *NoopUpdater) GetArtifactName() (string, error) {
	return "", errors.New("no updater available")
}

func (n *NoopUpdater) StartUpdate(url string) error {
	return errors.New("no updater available")
}

func (n *NoopUpdater) CancelUpdate() error {
	return errors.New("no updater available")
}

func (n *NoopUpdater) SubscribeUpdate() (*Client, error) {
	return nil, errors.New("no updater available")
}

func (n *NoopUpdater) UnsubscribeUpdate(client *Client) error {
	return errors.New("no updater available")
}
