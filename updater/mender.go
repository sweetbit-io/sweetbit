package updater

import (
	"bufio"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/the-lightning-land/sweetd/sweetdb"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var sizeRegexp = regexp.MustCompile("of size (\\d+)...")
var progressRegex = regexp.MustCompile("(\\d+)% (\\d+) KiB")

type MenderUpdaterConfig struct {
	Logger Logger
	DB     *sweetdb.DB
}

type nextClient struct {
	sync.Mutex
	id uint32
}

type MenderUpdater struct {
	log          Logger
	db           *sweetdb.DB
	artifactName string
	updating     int32 // atomic
	update       *Update
	updateCmd    *exec.Cmd
	progress     uint8
	shouldReboot bool
	shouldCommit bool
	clients      map[uint32]*UpdateClient
	nextClient   nextClient
}

// Compile time check for protocol compatibility
var _ Updater = (*MenderUpdater)(nil)

func NewMenderUpdater(config *MenderUpdaterConfig) (*MenderUpdater, error) {
	_, err := exec.LookPath("mender")
	if err != nil {
		return nil, errors.New("mender is not installed or missing in $PATH")
	}

	artifact, err := exec.Command("mender", "-show-artifact").Output()
	if err != nil {
		return nil, errors.Errorf("Could not retrieve artifact name: %v", err)
	}

	updater := &MenderUpdater{
		log:          config.Logger,
		db:           config.DB,
		artifactName: string(artifact),
		updating:     0,
		progress:     0,
		shouldReboot: false,
		shouldCommit: false,
		clients:      make(map[uint32]*UpdateClient),
		nextClient:   nextClient{id: 0},
		update:       nil,
	}

	if updater.log == nil {
		updater.log = noopLogger{}
	}

	return updater, nil
}

func (m *MenderUpdater) GetVersion() (string, error) {
	return m.artifactName, nil
}

func (m *MenderUpdater) GetUpdate(id string) (*Update, error) {
	update, err := m.db.GetUpdate(id)
	if err != nil {
		return nil, errors.Errorf("Could not get update with id %s: %v", id, err)
	}

	if update == nil {
		return nil, nil
	}

	result := &Update{
		Id:      update.Id,
		Started: update.Started,
		Url:     update.Url,
		State:   update.State,
	}

	// add progress information if it's the currently running update
	if update.Id == m.update.Id {
		result.Progress = m.progress
		result.ShouldReboot = m.shouldReboot
		result.ShouldCommit = m.shouldCommit
	}

	return result, nil
}

func (m *MenderUpdater) GetCurrentUpdate() (*Update, error) {
	return m.update, nil
}

func (m *MenderUpdater) StartUpdate(url string) (*Update, error) {
	if !atomic.CompareAndSwapInt32(&m.updating, 0, 1) {
		return nil, errors.New("update already running")
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Errorf("Could not generate id: %v", err)
	}

	m.update = &Update{
		Id:      id.String(),
		Url:     url,
		Started: time.Now(),
		State:   StateStarted,
	}

	err = m.db.SaveUpdate(&sweetdb.Update{
		Id:      m.update.Id,
		Started: m.update.Started,
		Url:     m.update.Url,
		State:   m.update.State,
	})
	if err != nil {
		return nil, errors.Errorf("Could not add update entry: %v", err)
	}

	m.updateCmd = exec.Command("mender", "-install", url)

	stdoutReader, err := m.updateCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	go func() {
		stdoutScanner := bufio.NewScanner(stdoutReader)
		for stdoutScanner.Scan() {
			text := stdoutScanner.Text()
			m.log.Debugf("%s", text)

			if matches := sizeRegexp.FindStringSubmatch(text); len(matches) == 2 {
				m.log.Infof("Total size of update is %s bytes", matches[1])
			}

			if matches := progressRegex.FindStringSubmatch(text); len(matches) == 3 {
				m.log.Infof("Update progressed %s%%", matches[1])
				progress, err := strconv.ParseUint(matches[1], 10, 8)
				if err != nil {
					m.log.Errorf("Could not parse progress: %v", err)
					continue
				}

				m.progress = uint8(progress)
				m.notifyUpdateClients()
			}
		}
	}()

	stderrReader, err := m.updateCmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	go func() {
		stderrScanner := bufio.NewScanner(stderrReader)
		for stderrScanner.Scan() {
			text := stderrScanner.Text()
			m.log.Errorf("%s", text)
		}
	}()

	if err := m.updateCmd.Start(); err != nil {
		m.update.State = StateFailed
		err = m.db.SaveUpdate(&sweetdb.Update{
			Id:      m.update.Id,
			Started: m.update.Started,
			Url:     m.update.Url,
			State:   m.update.State,
		})
		if err != nil {
			m.log.Errorf("could not save update state: %v", err)
		}
		return nil, err
	}

	go func() {
		err := m.updateCmd.Wait()
		if err != nil {
			m.update.State = StateFailed

			// TODO(davidknezic): detect cancelled state
			// m.update.State = StateCancelled
		} else {
			m.update.State = StateInstalled
			m.shouldReboot = true
		}

		err = m.db.SaveUpdate(&sweetdb.Update{
			Id:      m.update.Id,
			Started: m.update.Started,
			Url:     m.update.Url,
			State:   m.update.State,
		})
		if err != nil {
			m.log.Errorf("could not save update state: %v", err)
		}

		m.notifyUpdateClients()

		atomic.StoreInt32(&m.updating, 0)
	}()

	return m.update, nil
}

func (m *MenderUpdater) CancelUpdate(id string) (*Update, error) {
	if atomic.LoadInt32(&m.updating) != 1 {
		return nil, errors.New("no update in progress")
	}

	err := m.updateCmd.Process.Kill()
	if err != nil {
		return nil, errors.Errorf("could not kill update: %v", err)
	}

	m.update.State = StateCancelled

	return m.update, nil
}

func (m *MenderUpdater) SubscribeUpdate(id string) (*UpdateClient, error) {
	if m.update == nil || m.update.Id != id {
		return nil, nil
	}

	client := &UpdateClient{
		Update:     make(chan *Update),
		cancelChan: make(chan struct{}),
		updateId:   id,
		updater:    m,
	}

	m.nextClient.Lock()
	client.Id = m.nextClient.id
	m.nextClient.id++
	m.nextClient.Unlock()

	m.clients[client.Id] = client

	return client, nil
}

func (m *MenderUpdater) notifyUpdateClients() {
	for _, client := range m.clients {
		if m.update.Id == client.updateId {
			client.Update <- m.update
		}
	}
}

func (m *MenderUpdater) unsubscribeUpdate(client *UpdateClient) {
	delete(m.clients, client.Id)
	close(client.cancelChan)
}

func (m *MenderUpdater) CommitUpdate(id string) (*Update, error) {
	_, err := exec.Command("mender", "-commit").Output()
	if err != nil {
		return nil, errors.Errorf("Could not retrieve artifact name: %v", err)
	}

	m.update.State = StateCompleted

	return m.update, nil
}

func (m *MenderUpdater) RejectUpdate(id string) (*Update, error) {
	_, err := exec.Command("mender", "-rollback").Output()
	if err != nil {
		return nil, errors.Errorf("Could not retrieve artifact name: %v", err)
	}

	m.update.State = StateRejected

	return m.update, nil
}
