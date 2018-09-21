package wpa

import (
	"os/exec"
	"strings"
	"regexp"
	"time"
	"bytes"
	log "github.com/sirupsen/logrus"
)

type Wpa struct {
}

func NewWpa() *Wpa {
	return &Wpa{}
}

// WpaNetwork defines a wifi network to connect to.
type WpaNetwork struct {
	Bssid       string
	Frequency   string
	SignalLevel string
	Flags       string
	Ssid        string
}

// WpaCredentials defines wifi network credentials.
type WpaCredentials struct {
	Ssid string
	Psk  string
}

// WpaConnection defines a WPA connection.
type WpaConnection struct {
	Ssid    string
	State   string
	Ip      string
	Message string
}

// ConfiguredNetworks returns a list of configured wifi networks.
func (wpa *Wpa) ConfiguredNetworks() string {
	netOut, err := exec.Command("wpa_cli", "-i", "wlan0", "scan").Output()
	if err != nil {
		log.Error(err)
	}

	return string(netOut)
}

// ConnectNetwork connects to a wifi network
func (wpa *Wpa) ConnectNetwork(creds WpaCredentials) (WpaConnection, error) {
	connection := WpaConnection{}

	// 1. Add a network
	addNetOut, err := exec.Command("wpa_cli", "-i", "wlan0", "add_network").Output()
	if err != nil {
		log.Error(err)
		return connection, err
	}
	net := strings.TrimSpace(string(addNetOut))
	log.Info("WPA add network got: %s", net)

	// 2. Set the ssid for the new network
	addSsidOut, err := exec.Command("wpa_cli", "-i", "wlan0", "set_network", net, "ssid", "\""+creds.Ssid+"\"").Output()
	if err != nil {
		log.Error(err)
		return connection, err
	}
	ssidStatus := strings.TrimSpace(string(addSsidOut))
	log.Info("WPA add ssid got: %s", ssidStatus)

	// 3. Set the psk for the new network
	addPskOut, err := exec.Command("wpa_cli", "-i", "wlan0", "set_network", net, "psk", "\""+creds.Psk+"\"").Output()
	if err != nil {
		log.Error(err.Error())
		return connection, err
	}
	pskStatus := strings.TrimSpace(string(addPskOut))
	log.Info("WPA psk got: %s", pskStatus)

	// 4. Enable the new network
	enableOut, err := exec.Command("wpa_cli", "-i", "wlan0", "enable_network", net).Output()
	if err != nil {
		log.Error(err.Error())
		return connection, err
	}
	enableStatus := strings.TrimSpace(string(enableOut))
	log.Info("WPA enable got: %s", enableStatus)

	// regex for state
	rState := regexp.MustCompile("(?m)wpa_state=(.*)\n")

	// loop for status every second
	for i := 0; i < 5; i++ {
		log.Info("WPA Checking wifi state")

		stateOut, err := exec.Command("wpa_cli", "-i", "wlan0", "status").Output()
		if err != nil {
			log.Error("Got error checking state: %s", err.Error())
			return connection, err
		}
		ms := rState.FindSubmatch(stateOut)

		if len(ms) > 0 {
			state := string(ms[1])
			log.Info("WPA Enable state: %s", state)
			// see https://developer.android.com/reference/android/net/wifi/SupplicantState.html
			if state == "COMPLETED" {
				// save the config
				saveOut, err := exec.Command("wpa_cli", "-i", "wlan0", "save_config").Output()
				if err != nil {
					log.Error(err.Error())
					return connection, err
				}
				saveStatus := strings.TrimSpace(string(saveOut))
				log.Info("WPA save got: %s", saveStatus)

				connection.Ssid = creds.Ssid
				connection.State = state

				return connection, nil
			}
		}

		time.Sleep(3 * time.Second)
	}

	connection.State = "FAIL"
	connection.Message = "Unable to connection to " + creds.Ssid
	return connection, nil
}

// Status returns the WPA wireless status.
func (wpa *Wpa) Status() (map[string]string, error) {
	cfgMap := make(map[string]string, 0)

	stateOut, err := exec.Command("wpa_cli", "-i", "wlan0", "status").Output()
	if err != nil {
		log.Error("Got error checking state: %s", err.Error())
		return cfgMap, err
	}

	cfgMap = cfgMapper(stateOut)

	return cfgMap, nil
}

// cfgMapper takes a byte array and splits by \n and then by = and puts it all in a map.
func cfgMapper(data []byte) map[string]string {
	cfgMap := make(map[string]string, 0)

	lines := bytes.Split(data, []byte("\n"))

	for _, line := range lines {
		kv := bytes.Split(line, []byte("="))
		if len(kv) > 1 {
			cfgMap[string(kv[0])] = string(kv[1])
		}
	}

	return cfgMap
}

// ScanNetworks returns a map of WpaNetwork data structures.
func (wpa *Wpa) ScanNetworks() (map[string]WpaNetwork, error) {
	wpaNetworks := make(map[string]WpaNetwork, 0)

	scanOut, err := exec.Command("wpa_cli", "-i", "wlan0", "scan").Output()
	if err != nil {
		log.Error(err)
		return wpaNetworks, err
	}
	scanOutClean := strings.TrimSpace(string(scanOut))

	// wait one second for results
	time.Sleep(1 * time.Second)

	if scanOutClean == "OK" {
		networkListOut, err := exec.Command("wpa_cli", "-i", "wlan0", "scan_results").Output()
		if err != nil {
			log.Error(err)
			return wpaNetworks, err
		}

		networkListOutArr := strings.Split(string(networkListOut), "\n")
		for _, netRecord := range networkListOutArr[1:] {
			if strings.Contains(netRecord, "[P2P]") {
				continue
			}

			fields := strings.Fields(netRecord)

			if len(fields) > 4 {
				ssid := strings.Join(fields[4:], " ")
				wpaNetworks[ssid] = WpaNetwork{
					Bssid:       fields[0],
					Frequency:   fields[1],
					SignalLevel: fields[2],
					Flags:       fields[3],
					Ssid:        ssid,
				}
			}
		}

	}

	return wpaNetworks, nil
}
