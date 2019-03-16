package ap

import (
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
	"github.com/the-lightning-land/sweetd/dnsmasq"
	"github.com/the-lightning-land/sweetd/hostapd"
	"os/exec"
	"strings"
)

type DispenserApConfig struct {
	Hotspot *DispenserApHotspotConfig
}

type DispenserApHotspotConfig struct {
	Ip         string
	Interface  string
	Ssid       string
	Passphrase string
	DhcpRange  string
}

type DispenserAp struct {
	config          *DispenserApConfig
	hostapdInstance *hostapd.Hostapd
	dnsmasqInstance *dnsmasq.Dnsmasq
}

// Ensure we implement the Ap interface with this compile-time check
var _ Ap = (*DispenserAp)(nil)

func NewDispenserAp(config *DispenserApConfig) (*DispenserAp, error) {
	return &DispenserAp{
		config: config,
	}, nil
}

func (a *DispenserAp) Start() error {
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

	if err := a.startHotspot(); err != nil {
		return errors.Errorf("Could not start hotspot: %v", err)
	}

	return nil
}

func (a *DispenserAp) ConnectWifi(ssid string, psk string) error {
	return nil
}

func (a *DispenserAp) startHotspot() error {
	log.Info("Starting hostapd for access point management...")

	var err error

	a.hostapdInstance, err = hostapd.New(&hostapd.Config{
		Channel:    6,
		Ssid:       a.config.Hotspot.Ssid,
		Passphrase: a.config.Hotspot.Passphrase,
		Log: func(s string) {
			log.WithField("service", "hostapd").Debug(s)
		},
	})
	if err != nil {
		return errors.Errorf("Could not create hostapd service: %v", err)
	}

	if err := a.hostapdInstance.Start(); err != nil {
		return errors.Errorf("Could not start hostapd: %v", err)
	}

	log.Info("Started hostapd.")

	log.Info("Creating dnsmasq for DNS and DHCP management...")

	a.dnsmasqInstance, err = dnsmasq.New(&dnsmasq.Config{
		Address:   "/#/" + strings.Split(a.config.Hotspot.Ip, "/")[0],
		DhcpRange: a.config.Hotspot.DhcpRange,
		Log: func(s string) {
			log.WithField("service", "dnsmasq").Debug(s)
		},
	})

	if err != nil {
		return errors.Errorf("Could not create dnsmasq service: %v", err)
	}

	log.Info("Starting dnsmasq...")

	if err := a.dnsmasqInstance.Start(); err != nil {
		return errors.Errorf("Could not start dnsmasq: %v", err)
	}

	log.Info("Started dnsmasq.")

	return nil
}

func (a *DispenserAp) Stop() error {
	if err := downApInterface(a.config.Hotspot.Interface); err != nil {
		return errors.Errorf("Could not down AP interface: %v", err)
	}

	if err := removeApInterface(a.config.Hotspot.Interface); err != nil {
		return errors.Errorf("Could not remove AP interface: %v", err)
	}

	return nil
}

// removeApInterface removes the AP interface.
func removeApInterface(iface string) error {
	cmd := exec.Command("iw", "dev", iface, "del")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

// configureApInterface configured the AP interface.
func configureApInterface(ip string, iface string) error {
	cmd := exec.Command("ip", "addr", "add", ip, "dev", iface)

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

// upApInterface ups the AP Interface.
func upApInterface(iface string) error {
	cmd := exec.Command("ip", "link", "set", iface, "up")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

// downApInterface downs the AP Interface.
func downApInterface(iface string) error {
	cmd := exec.Command("ip", "link", "set", iface, "down")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}

// addApInterface adds the AP interface.
func addApInterface(iface string) error {
	cmd := exec.Command("iw", "phy", "phy0", "interface", "add", iface, "type", "__ap")

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return nil
	}

	return nil
}
