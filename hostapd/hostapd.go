package hostapd

import (
	"text/template"
	"os/exec"
	"bufio"
	"strings"
	"sync/atomic"
	"github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"
)

var configTemplString = `interface=uap0
ssid={{.Ssid}}
hw_mode=g
channel=6
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
}

type Hostapd struct {
	started  int32 // atomic
	config   *Config
	cmd      *exec.Cmd
	messages chan string
}

func New(config *Config) *Hostapd {
	return &Hostapd{
		config:   config,
		cmd:      exec.Command("hostapd", "-d", "/dev/stdin"),
		messages: make(chan string),
	}
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
			log.Info(stdOutScanner.Text())
			h.messages <- stdOutScanner.Text()
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
		out := <-h.messages
		if strings.Contains(out, "AP-ENABLED") {
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
		out := <-h.messages
		if strings.Contains(out, "AP-DISABLED") {
			return nil
		}
	}
}
