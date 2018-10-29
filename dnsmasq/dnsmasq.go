package dnsmasq

import (
	"os/exec"
	"bufio"
	"sync/atomic"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Config struct {
	Address   string
	DhcpRange string
	Log       func(string) ()
}

type DnsmasqState int

const (
	STARTED DnsmasqState = iota
)

type Dnsmasq struct {
	started int32 // atomic
	config  *Config
	cmd     *exec.Cmd
	states  chan DnsmasqState
}

func New(config *Config) *Dnsmasq {
	args := []string{
		"--no-hosts", // Don't read the hostnames in /etc/hosts.
		"--keep-in-foreground",
		"--log-queries",
		"--no-resolv",
		"--address=" + config.Address,
		"--dhcp-range=" + config.DhcpRange,
		"--dhcp-vendorclass=set:device,IoT",
		"--dhcp-authoritative",
		"--log-facility=-",
		"--interface=uap0",
	}

	return &Dnsmasq{
		config: config,
		cmd:    exec.Command("dnsmasq", args...),
		states: make(chan DnsmasqState),
	}
}

func (d *Dnsmasq) Start() error {
	// Already running?
	if !atomic.CompareAndSwapInt32(&d.started, 0, 1) {
		return errors.New("dnsmasq already started")
	}

	cmdStderrReader, err := d.cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdErrScanner := bufio.NewScanner(cmdStderrReader)
	go func() {
		// Goroutine will finish on process end
		for stdErrScanner.Scan() {
			text := stdErrScanner.Text()

			if d.config.Log != nil {
				d.config.Log(text)
			}

			if strings.Contains(text, "started, version") {
				d.states <- STARTED
			}
		}
	}()

	if err := d.cmd.Start(); err != nil {
		return err
	}

	timer := time.NewTimer(60 * time.Second)

	// Block until the process has started
	for {
		<-timer.C
		return errors.New("Timed out")
		state := <-d.states
		if state == STARTED {
			return nil
		}
	}
}

func (d *Dnsmasq) Stop() error {
	if atomic.LoadInt32(&d.started) != 1 {
		return errors.New("dnsmasq not started yet")
	}

	d.cmd.Process.Kill()

	return nil
}
