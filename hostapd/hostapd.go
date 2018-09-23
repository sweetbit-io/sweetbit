package hostapd

import (
	"text/template"
	"os/exec"
	"bufio"
	"strings"
	"sync/atomic"
	"github.com/pkg/errors"
)

var configTemplString = `interface=uap0
ssid={{.Ssid}}
hw_mode=g
channel=11
macaddr_acl=0
auth_algs=1
ignore_broadcast_ssid=0
wpa=2
wpa_passphrase={{.Passphrase}}
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP
`

type Config struct {
	Ssid       string
	Passphrase string
	Log        func(string) ()
}

type HostapdState int

const (
	ENABLED  HostapdState = iota
	DISABLED
)

type Hostapd struct {
	started int32 // atomic
	config  *Config
	cmd     *exec.Cmd
	states  chan HostapdState // Signals ENABLED OR DISABLED from the log
}

func New(config *Config) (*Hostapd, error) {
	passphraseLen := len(config.Passphrase)
	if passphraseLen < 8 || passphraseLen > 63 {
		return nil, errors.New("Invalid WPA passphrase length (expected 8..63)")
	}

	return &Hostapd{
		config: config,
		cmd:    exec.Command("hostapd", "-d", "/dev/stdin"),
		states: make(chan HostapdState),
	}, nil
}

func (h *Hostapd) Start() error {
	// Already running?
	if !atomic.CompareAndSwapInt32(&h.started, 0, 1) {
		return errors.New("hostapd already started")
	}

	// Create config template
	tmpl, err := template.New("hostapd").Parse(configTemplString)
	if err != nil {
		return err
	}

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
	}()

	// Write config to stdin, where the process is set up to read it from
	err = tmpl.Execute(hostapdPipe, h.config)
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
			return nil
		}
	}
}

func (h *Hostapd) Stop() error {
	if atomic.LoadInt32(&h.started) != 1 {
		return errors.New("hostapd not started yet")
	}

	h.cmd.Process.Kill()

	// Block until the process has finished
	for {
		state := <-h.states
		if state == DISABLED {
			return nil
		}
	}
}
