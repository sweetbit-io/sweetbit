package ap

import (
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// TODO Better handle the following wpa_cli states
// DISCONNECTED
// INACTIVE
// ASSOCIATING
// ASSOCIATED
// COMPLETED

type DispenserApConfig struct {
	Interface string
}

type DispenserAp struct {
	config           *DispenserApConfig
	networks         []*Network
	connectionStatus *ConnectionStatus
	done             chan bool

	// Ap clients
	apClients      map[uint32]*ApClient
	apClientMtx    sync.Mutex
	nextApClientID uint32
}

// Ensure we implement the Ap interface with this compile-time check
var _ Ap = (*DispenserAp)(nil)

func NewDispenserAp(config *DispenserApConfig) (*DispenserAp, error) {
	return &DispenserAp{
		config:    config,
		apClients: make(map[uint32]*ApClient),
	}, nil
}

func (a *DispenserAp) Start() error {
	a.done = make(chan bool)

	a.syncConnectionStatus()

	go func() {
		statusTicker := time.NewTicker(5 * time.Second)
		defer statusTicker.Stop()

		for {
			select {
			case <-statusTicker.C:
				a.syncConnectionStatus()
				break
			case <-a.done:
				log.Infof("Stopping connection status sync...")
				return
			}
		}
	}()

	return nil
}

func (a *DispenserAp) syncConnectionStatus() {
	status, err := getStatus(a.config.Interface)
	if err != nil {
		log.Errorf("Getting Wifi connection status failed: %v", err)
		return
	} else if status.State == "SCANNING" || status.State == "ASSOCIATING" || status.State == "ASSOCIATED" || status.State == "4WAY_HANDSHAKE" {
		// log.Infof("Skipping state: %v", status.State)
		return
	}

	info, err := getApInterfaceInfo(a.config.Interface)
	if err != nil {
		log.Errorf("Getting interface info failed: %v", err)
		return
	}

	notify := false

	if a.connectionStatus == nil {
		log.Infof("Set channel to %v", info.channel)
	} else if a.connectionStatus.Channel != info.channel {
		log.Infof("Changed channel from %v to %v", a.connectionStatus.Channel, info.channel)
	}

	if a.connectionStatus == nil {
		notify = true
		log.Infof("Set Wifi to %v", status.Ssid)
	} else if a.connectionStatus.Ssid != status.Ssid {
		notify = true
		log.Infof("Changed Wifi from %v to %v", a.connectionStatus.Ssid, status.Ssid)
	}

	if a.connectionStatus == nil {
		notify = true
		log.Infof("Set state to %v", status.State)
	} else if a.connectionStatus.State != status.State {
		notify = true
		log.Infof("Changed state from %v to %v", a.connectionStatus.State, status.State)
	}

	if a.connectionStatus == nil {
		notify = true
		log.Infof("Set ip to %v", status.Ip)
	} else if a.connectionStatus.Ip != status.Ip {
		notify = true
		log.Infof("Changed ip from %v to %v", a.connectionStatus.Ip, status.Ip)
	}

	// TODO: this might be too early
	if a.connectionStatus != nil && a.connectionStatus.Ssid != status.Ssid {
		log.Infof("Renewing IP address")

		err := renewUdhcp()
		if err != nil {
			log.Errorf("Could not renew address")
		}
	}

	a.connectionStatus = status
	a.connectionStatus.Channel = info.channel

	if notify {
		log.Infof("Notifying of network update...")

		for _, client := range a.apClients {
			log.Infof("Notifying client %v...", client.Id)

			client.Updates <- &ApUpdate{
				Ip:        status.Ip,
				Ssid:      status.Ssid,
				Connected: status.State == "COMPLETED",
			}
		}
	}
}

func (a *DispenserAp) ScanWifiNetworks() error {
	err := scan(a.config.Interface)
	if err != nil {
		return errors.Errorf("Scan failed: %v", err)
	}

	return nil
}

func (a *DispenserAp) ListWifiNetworks() ([]*Network, error) {
	results, err := results(a.config.Interface)
	if err != nil {
		log.Errorf("Result failed: %v", err)
	}

	return results, nil
}

func (a *DispenserAp) ConnectWifi(ssid string, psk string) error {
	net, err := addNetwork(a.config.Interface)
	if err != nil {
		return errors.Errorf("Adding network failed: %v", err)
	}

	err = setNetwork(a.config.Interface, net, ssidKey, ssid)
	if err != nil {
		return errors.Errorf("Setting ssid failed: %v", err)
	}

	err = setNetwork(a.config.Interface, net, pskKey, psk)
	if err != nil {
		return errors.Errorf("Setting psk failed: %v", err)
	}

	err = enableNetwork(a.config.Interface, net)
	if err != nil {
		return errors.Errorf("Enabling network failed: %v", err)
	}

	err = a.removeAllConfiguredNetworksExcept([]networkId{net})
	if err != nil {
		return errors.Errorf("Could not remove old networks: %v", err)
	}

	return nil
}

func (a *DispenserAp) removeAllConfiguredNetworksExcept(except []networkId) error {
	configuredNetworks, err := listConfiguredNetworks(a.config.Interface)
	if err != nil {
		return errors.Errorf("Listing configured networks failed: %v", err)
	} else {
		for _, configuredNetwork := range configuredNetworks {
			exclude := false

			for i := 0; i < len(except); i++ {
				if except[i] == configuredNetwork.id {
					exclude = true
					break
				}
			}

			if exclude {
				continue
			}

			err := removeNetwork(a.config.Interface, configuredNetwork.id)
			if err != nil {
				return errors.Errorf("Unable to remove configured network %v", configuredNetwork.id)
			}
		}
	}

	return nil
}

func (a *DispenserAp) GetConnectionStatus() (*ConnectionStatus, error) {
	return a.connectionStatus, nil
}

func (a *DispenserAp) Stop() error {
	close(a.done)

	err := a.removeAllConfiguredNetworksExcept([]networkId{})
	if err != nil {
		return errors.Errorf("Could not remove networks: %v", err)
	}

	return nil
}

func (a *DispenserAp) SubscribeUpdates() *ApClient {
	client := &ApClient{
		Updates:    make(chan *ApUpdate),
		cancelChan: make(chan struct{}),
		ap:         a,
	}

	a.apClientMtx.Lock()
	client.Id = a.nextApClientID
	a.nextApClientID++
	a.apClientMtx.Unlock()

	a.apClients[client.Id] = client

	return client
}

func (a *DispenserAp) deleteApClient(id uint32) {
	delete(a.apClients, id)
}
