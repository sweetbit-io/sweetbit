package ap

type Ap interface {
	Start() error
	ConnectWifi(ssid string, psk string) error
	Stop() error
}