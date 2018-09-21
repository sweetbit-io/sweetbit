package dnsmasq

import (
	"os/exec"
	"bufio"
	"strings"
	"sync/atomic"
	"github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"
)

type Config struct {
	Address   string
	DhcpRange string
}

type Dnsmasq struct {
	started  int32 // atomic
	config   *Config
	cmd      *exec.Cmd
	messages chan string
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
	}

	return &Dnsmasq{
		config:   config,
		cmd:      exec.Command("dnsmasq", args...),
		messages: make(chan string),
	}
}

func (h *Dnsmasq) Start() error {
	// Already running?
	if !atomic.CompareAndSwapInt32(&h.started, 0, 1) {
		return errors.New("dnsmasq already started")
	}

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

	if err := h.cmd.Start(); err != nil {
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

func (h *Dnsmasq) Stop() error {
	if atomic.LoadInt32(&h.started) != 1 {
		return errors.New("dnsmasq not started yet")
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
