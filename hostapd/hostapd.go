package hostapd

import (
	"bufio"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"text/template"
)

var configTemplString = `interface=uap0
ssid={{.Ssid}}
hw_mode=g
channel={{.Channel}}
wmm_enabled=0
macaddr_acl=0
auth_algs=1
ignore_broadcast_ssid=0
wpa=2
wpa_passphrase={{.Passphrase}}
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP
`

type configTmpl struct {
	Ssid       string
	Passphrase string
	Channel    int
}

type Config struct {
	Ssid       string
	Passphrase string
	Log        func(string) ()
}

type HostapdState int

const (
	ENABLED HostapdState = iota
	DISABLED
)

type Hostapd struct {
	channel int
	started int32 // atomic
	config  *Config
	cmd     *exec.Cmd
	states  chan HostapdState // Signals ENABLED OR DISABLED from the log
}

func New(config *Config) (*Hostapd, error) {
	_, err := exec.LookPath("hostapd")
	if err != nil {
		return nil, errors.New("hostapd is not installed or missing in $PATH")
	}

	passphraseLen := len(config.Passphrase)
	if passphraseLen < 8 || passphraseLen > 63 {
		return nil, errors.New("Invalid WPA passphrase length (expected 8..63)")
	}

	return &Hostapd{
		config: config,
		states: make(chan HostapdState),
	}, nil
}

func (h *Hostapd) Start(channel int) error {
	// Already running?
	if !atomic.CompareAndSwapInt32(&h.started, 0, 1) {
		return errors.New("hostapd already started")
	}

	// Create config template
	tmpl, err := template.New("hostapd").Parse(configTemplString)
	if err != nil {
		return err
	}

	h.cmd = exec.Command("hostapd", "-d", "/dev/stdin")

	hostapdPipe, _ := h.cmd.StdinPipe()
	cmdStdoutReader, err := h.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stdOutScanner := bufio.NewScanner(cmdStdoutReader)
	go func() {
		// Goroutine will finish on process end
		for stdOutScanner.Scan() {
			text := stdOutScanner.Text()

			if h.config.Log != nil {
				h.config.Log(text)
			}

			if strings.Contains(text, "AP-ENABLED") {
				h.states <- ENABLED
			} else if strings.Contains(text, "AP-DISABLED") {
				h.states <- DISABLED
			}
		}

		h.states <- DISABLED
	}()

	// Write config to stdin, where the process is set up to read it from
	err = tmpl.Execute(hostapdPipe, configTmpl{
		Ssid:       h.config.Ssid,
		Passphrase: h.config.Passphrase,
		Channel:    channel,
	})
	if err != nil {
		return err
	}

	err = h.cmd.Start()
	hostapdPipe.Close()

	if err != nil {
		return err
	}

	// Block until the process has started
	for {
		state := <-h.states
		if state == ENABLED {
			// Remember the channel hostapd was started with
			h.channel = channel
			return nil
		}
	}
}

func (h *Hostapd) ChangeChannel(channel int) error {
	if channel == h.channel {
		// don't restart if desired channel is already set
		return nil
	}

	err := h.Stop()
	if err != nil {
		return errors.Errorf("Could not stop hostapd: %v", err)
	}

	err = h.Start(channel)
	if err != nil {
		return errors.Errorf("Could not start hostapd: %v", err)
	}

	return nil
}

func (h *Hostapd) Stop() error {
	if !atomic.CompareAndSwapInt32(&h.started, 1, 0) {
		return errors.New("hostapd not started")
	}

	h.cmd.Process.Signal(os.Interrupt)

	// Block until the process has finished
	for {
		state := <-h.states
		if state == DISABLED {
			return nil
		}
	}
}
