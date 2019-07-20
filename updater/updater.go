package updater

type Updater interface {
	GetArtifactName() (string, error)
	StartUpdate(url string) error
	CancelUpdate() error
	SubscribeUpdate() (*Client, error)
	UnsubscribeUpdate(client *Client) error
}