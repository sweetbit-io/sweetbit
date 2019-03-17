package ap

import (
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
	"github.com/the-lightning-land/sweetd/dnsmasq"
	"github.com/the-lightning-land/sweetd/hostapd"
	"strings"
	"time"
)

// TODO Better handle the following wpa_cli states
// DISCONNECTED
// INACTIVE
// ASSOCIATING
// ASSOCIATED

type DispenserApConfig struct {
	Hotspot   *DispenserApHotspotConfig
	Interface string
}

type DispenserApHotspotConfig struct {
	Ip         string
	Interface  string
	Ssid       string
	Passphrase string
	DhcpRange  string
}

type DispenserAp struct {
	config           *DispenserApConfig
	hostapdInstance  *hostapd.Hostapd
	dnsmasqInstance  *dnsmasq.Dnsmasq
	networks         []*Network
	connectionStatus *ConnectionStatus
	done             chan bool
}

// Ensure we implement the Ap interface with this compile-time check
var _ Ap = (*DispenserAp)(nil)

func NewDispenserAp(config *DispenserApConfig) (*DispenserAp, error) {
	return &DispenserAp{
		config: config,
	}, nil
}

func (a *DispenserAp) Start() error {
	a.done = make(chan bool)

	log.Infof("Setting up %s interface as access point...", a.config.Hotspot.Interface)

	if err := removeApInterface(a.config.Hotspot.Interface); err != nil {
		return errors.Errorf("Could not remove AP interface: %v", err)
	}

	if err := addApInterface(a.config.Hotspot.Interface); err != nil {
		return errors.Errorf("Could not add AP interface: %v", err)
	}

	if err := upApInterface(a.config.Hotspot.Interface); err != nil {
		return errors.Errorf("Could not up AP interface: %v", err)
	}

	if err := configureApInterface(a.config.Hotspot.Ip, a.config.Hotspot.Interface); err != nil {
		return errors.Errorf("Could not configure AP interface: %v", err)
	}

	a.syncConnectionStatus()

	go a.pollNetworksAndConnectionStatus()

	return nil
}

func (a *DispenserAp) pollNetworksAndConnectionStatus() {
	statusTicker := time.NewTicker(1 * time.Second)
	defer statusTicker.Stop()

	for {
		select {
		case <-statusTicker.C:
			a.syncConnectionStatus()
			break
		}
	}
}

func (a *DispenserAp) syncConnectionStatus() {
	status, err := getStatus(a.config.Interface)
	if err != nil {
		log.Errorf("Getting Wifi connection status failed: %v", err)
		return
	} else if status.State == "SCANNING" || status.State == "ASSOCIATING" || status.State == "ASSOCIATED" {
		// Skip processing state SCANNING and try in next cycle
		return
	}

	info, err := getApInterfaceInfo(a.config.Interface)
	if err != nil {
		log.Errorf("Getting interface info failed: %v", err)
		return
	}

	if a.connectionStatus == nil {
		log.Infof("Set channel to %v", info.channel)
	} else if a.connectionStatus.Channel != info.channel {
		log.Infof("Changed channel from %v to %v", a.connectionStatus.Channel, info.channel)
	}

	if a.connectionStatus == nil {
		log.Infof("Set Wifi to %v", status.Ssid)
	} else if a.connectionStatus.Ssid != status.Ssid {
		log.Infof("Changed Wifi from %v to %v", a.connectionStatus.Ssid, status.Ssid)
	}

	if a.connectionStatus == nil {
		log.Infof("Set state to %v", status.State)
	} else if a.connectionStatus.State != status.State {
		log.Infof("Changed state from %v to %v", a.connectionStatus.State, status.State)
	}

	if a.connectionStatus == nil {
		log.Infof("Set ip to %v", status.Ip)
	} else if a.connectionStatus.Ip != status.Ip {
		log.Infof("Changed ip from %v to %v", a.connectionStatus.Ip, status.Ip)
	}

	a.connectionStatus = status
	a.connectionStatus.Channel = info.channel
}

func (a *DispenserAp) ListWifiNetworks() ([]*Network, error) {
	statusTicker := time.NewTicker(1 * time.Second)
	defer statusTicker.Stop()

	err := scan(a.config.Interface)
	if err != nil {
		return nil, errors.Errorf("Scan failed: %v", err)
	}

	for {
		<-statusTicker.C

		status, err := getStatus(a.config.Interface)
		if err != nil {
			return nil, errors.Errorf("Getting Wifi connection status failed: %v", err)
		}

		if status.State == "SCANNING" {
			// Skip processing state SCANNING and try in next cycle
			continue
		}

		results, err := results(a.config.Interface)
		if err != nil {
			log.Errorf("Result failed: %v", err)
		}

		return results, nil
	}
}

func (a *DispenserAp) ConnectWifi(ssid string, psk string) error {
	shouldRestartHostapd := false

	if a.hostapdInstance != nil {
		shouldRestartHostapd = true

		log.Info("Stopping hostapd...")

		err := a.hostapdInstance.Stop()
		if err != nil {
			return errors.Errorf("Could not stop hostapd: %v", err)
		}

		log.Info("Stopped hostapd")
	}

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

	if shouldRestartHostapd {
		log.Info("Starting hostapd...")

		a.hostapdInstance, err = a.startHostapd()
		if err != nil {
			return errors.Errorf("Could not start hostapd: %v", err)
		}

		log.Info("Started hostapd")
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

func (a *DispenserAp) StartHotspot() error {
	var err error

	a.hostapdInstance, err = a.startHostapd()
	if err != nil {
		return errors.Errorf("Could not start hostapd: %v", err)
	}

	a.dnsmasqInstance, err = a.startDnsmasq()
	if err != nil {
		return errors.Errorf("Could not start dnsmasq: %v", err)
	}

	return nil
}

func (a *DispenserAp) startHostapd() (*hostapd.Hostapd, error) {
	log.Info("Starting hostapd for access point management...")

	hostapdInstance, err := hostapd.New(&hostapd.Config{
		Ssid:       a.config.Hotspot.Ssid,
		Passphrase: a.config.Hotspot.Passphrase,
		// TODO: just allow passing an object with a Logger interface
		Log: func(s string) {
			log.WithField("service", "hostapd").Debug(s)
		},
	})
	if err != nil {
		return nil, errors.Errorf("Could not create: %v", err)
	}

	if err := hostapdInstance.Start(a.connectionStatus.Channel); err != nil {
		return nil, errors.Errorf("Could not start: %v", err)
	}

	log.Info("Started hostapd.")

	return hostapdInstance, nil
}

func (a *DispenserAp) startDnsmasq() (*dnsmasq.Dnsmasq, error) {
	log.Info("Creating dnsmasq for DNS and DHCP management...")

	dnsmasqInstance, err := dnsmasq.New(&dnsmasq.Config{
		Address:   "/#/" + strings.Split(a.config.Hotspot.Ip, "/")[0],
		DhcpRange: a.config.Hotspot.DhcpRange,
		// TODO: just allow passing an object with a Logger interface
		Log: func(s string) {
			log.WithField("service", "dnsmasq").Debug(s)
		},
	})

	if err != nil {
		return nil, errors.Errorf("Could not create: %v", err)
	}

	log.Info("Starting dnsmasq...")

	if err := dnsmasqInstance.Start(); err != nil {
		return nil, errors.Errorf("Could not start: %v", err)
	}

	log.Info("Started dnsmasq.")

	return dnsmasqInstance, nil
}

func (a *DispenserAp) Stop() error {
	close(a.done)

	if a.dnsmasqInstance != nil {
		a.dnsmasqInstance.Stop()
	}

	if a.hostapdInstance != nil {
		a.hostapdInstance.Stop()
	}

	err := a.removeAllConfiguredNetworksExcept([]networkId{})
	if err != nil {
		return errors.Errorf("Could not remove networks: %v", err)
	}

	if err := downApInterface(a.config.Hotspot.Interface); err != nil {
		return errors.Errorf("Could not down AP interface: %v", err)
	}

	if err := removeApInterface(a.config.Hotspot.Interface); err != nil {
		return errors.Errorf("Could not remove AP interface: %v", err)
	}

	return nil
}
