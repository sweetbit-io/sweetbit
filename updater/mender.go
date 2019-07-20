package updater

import (
	"bufio"
	"github.com/go-errors/errors"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
)

type MenderUpdaterConfig struct {
	Logger      Logger
	configPath  string
	dataDirPath string
}

type MenderUpdater struct {
	log          Logger
	configPath   string
	dataDirPath  string
	updating     int32 // atomic
	updateCmd    *exec.Cmd
	clients      map[uint32]*Client
	clientMtx    sync.Mutex
	nextClientID uint32
}

// Compile time check for protocol compatibility
var _ Updater = (*MenderUpdater)(nil)

func NewMenderUpdater(config *MenderUpdaterConfig) (*MenderUpdater, error) {
	_, err := exec.LookPath("mender")
	if err != nil {
		return nil, errors.New("mender is not installed or missing in $PATH")
	}

	updater := &MenderUpdater{}
	updater.configPath = config.configPath
	updater.dataDirPath = config.dataDirPath

	if config.Logger != nil {
		updater.log = config.Logger
	} else {
		updater.log = noopLogger{}
	}

	updater.clients = make(map[uint32]*Client)

	return updater, nil
}

func (m *MenderUpdater) GetArtifactName() (string, error) {
	args := []string{
		"-config " + m.configPath,
		"-data " + m.dataDirPath,
		"-show-artifact",
	}

	result, err := exec.Command("mender", args...).Output()
	if err != nil {
		return "", errors.Errorf("Command: %s", err.Error())
	}

	return string(result), nil
}

func (m *MenderUpdater) StartUpdate(url string) error {
	// Already running?
	if !atomic.CompareAndSwapInt32(&m.updating, 0, 1) {
		return errors.New("update already running")
	}

	args := []string{
		"-config " + m.configPath,
		"-data " + m.dataDirPath,
		"-install " + url,
	}

	m.updateCmd = exec.Command("mender", args...)

	cmdStderrReader, err := m.updateCmd.StderrPipe()
	if err != nil {
		return err
	}

	stdErrScanner := bufio.NewScanner(cmdStderrReader)
	go func() {
		// Goroutine will finish on process end
		for stdErrScanner.Scan() {
			text := stdErrScanner.Text()

			if m.log != nil {
				m.log.Infof(text)
			}

			if strings.Contains(text, "started, version") {
				// Notify all subscribed clients
				for _, client := range m.clients {
					client.Progress <- true
				}
			}
		}
	}()

	if err := m.updateCmd.Start(); err != nil {
		return err
	}

	return nil
}

func (m *MenderUpdater) CancelUpdate() error {
	if atomic.LoadInt32(&m.updating) != 1 {
		return errors.New("no update in progress")
	}

	err := m.updateCmd.Process.Kill()
	if err != nil {
		return errors.Errorf("could not kill update: %v", err)
	}

	return nil
}

func (m *MenderUpdater) SubscribeUpdate() (*Client, error) {
	client := &Client{
		Progress:   make(chan bool),
		cancelChan: make(chan struct{}),
		updater:    m,
	}

	m.clientMtx.Lock()
	client.Id = m.nextClientID
	m.nextClientID++
	m.clientMtx.Unlock()

	m.clients[client.Id] = client

	return client, nil
}

func (m *MenderUpdater) UnsubscribeUpdate(client *Client) error {
	delete(m.clients, client.Id)
	close(client.cancelChan)
	return nil
}
