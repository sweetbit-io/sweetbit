package ap

type Network struct {
	Bssid       string
	Frequency   string
	SignalLevel string
	Flags       string
	Ssid        string
}

type ConnectionStatus struct {
	Ssid    string
	State   string
	Ip      string
	Channel int
}

type Ap interface {
	Start() error
	ScanWifiNetworks() error
	ListWifiNetworks() ([]*Network, error)
	ConnectWifi(ssid string, psk string) error
	GetConnectionStatus() (*ConnectionStatus, error)
	Stop() error
	SubscribeUpdates() *ApClient
	deleteApClient(id uint32)
}
