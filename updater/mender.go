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

var sizeRegexp = regexp.MustCompile(`of size (\d+)\.\.\.`)
var progressRegex = regexp.MustCompile(`(\d+)% (\d+) KiB`)
var artifactNameRegex = regexp.MustCompile(`(?s)(.*No match between boot and root partitions.*\s)?(\S+)\s$`)

type artifactNameOutput struct {
	artifactName      string
	partitionMismatch bool
}

func parseArtifactNameOutput(output string) *artifactNameOutput {
	match := artifactNameRegex.FindStringSubmatch(output)
	if len(match) != 3 {
		return nil
	}

	return &artifactNameOutput{
		artifactName:      match[2],
		partitionMismatch: match[1] != "",
	}
}

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

func NewMenderUpdater(config *MenderUpdaterConfig) *MenderUpdater {
	updater := &MenderUpdater{
		log:          config.Logger,
		db:           config.DB,
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

	return updater
}

func (m *MenderUpdater) Setup() error {
	_, err := exec.LookPath("mender")
	if err != nil {
		return errors.New("unable to find mender in $PATH")
	}

	output, err := exec.Command("mender", "-show-artifact").Output()
	if err != nil {
		return errors.Errorf("unable to read mender artifact name: %v", err)
	}

	artifactNameOutput := parseArtifactNameOutput(string(output))
	if artifactNameOutput == nil {
		return errors.Errorf("unable to parse artifact name output: %v", err)
	}

	m.artifactName = artifactNameOutput.artifactName

	currentUpdate, err := m.db.GetCurrentUpdate()
	if err != nil {
		return errors.Errorf("unable to get current update: %v", err)
	}

	if currentUpdate != nil && currentUpdate.State == StateInstalled {
		// resume handling of update when it's installed, for example after reboot

		m.update = &Update{
			Id:      currentUpdate.Id,
			Started: currentUpdate.Started,
			Url:     currentUpdate.Url,
			State:   currentUpdate.State,
		}

		upgradeAvailable, err := checkUpgradeAvailable()
		if err != nil {
			m.log.Errorf("unable to check available upgrade: %v", err)
		}

		if artifactNameOutput.partitionMismatch {
			// we haven't rebooted yet and can detect this through an error
			// message returned when obtaining the artifact name
			m.shouldReboot = true
		} else if upgradeAvailable {
			// signal that the installed update was booted into and needs confirmation
			m.shouldCommit = true
		} else {
			// assume the update was rejected and update its state accordingly
			// don't notify subscribers as they're not subscribed yet during the setup phase
			m.update.State = StateRejected

			m.saveUpdate(m.update)
		}
	} else if currentUpdate != nil {
		// clear old update as it is considered stale if it has any
		// other state than installed

		err := m.db.ClearCurrentUpdate()
		if err != nil {
			m.log.Errorf("unable to clear old update: %v", err)
		}
	}

	return nil
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
	if m.update != nil && update.Id == m.update.Id {
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

	update := &Update{
		Id:      id.String(),
		Url:     url,
		Started: time.Now(),
		State:   StateStarted,
	}

	err = m.db.SaveUpdate(&sweetdb.Update{
		Id:      update.Id,
		Started: update.Started,
		Url:     update.Url,
		State:   update.State,
	})
	if err != nil {
		return nil, errors.Errorf("unable to save update: %v", err)
	}

	err = m.db.SetCurrentUpdate(update.Id)
	if err != nil {
		return nil, errors.Errorf("unable to set current update: %v", err)
	}

	m.updateCmd = exec.Command("mender", "-install", url)

	stdoutReader, err := m.updateCmd.StdoutPipe()
	if err != nil {
		return nil, errors.Errorf("unable to pipe stdout: %v", err)
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
				progress, err := strconv.ParseUint(matches[1], 10, 8)
				if err != nil {
					m.log.Errorf("Could not parse progress: %v", err)
					continue
				}

				m.log.Infof("Update progressed %d%%", progress)

				m.progress = uint8(progress)
				m.notifyUpdateClients(m.update)
			}
		}
	}()

	stderrReader, err := m.updateCmd.StderrPipe()
	if err != nil {
		return nil, errors.Errorf("unable to pipe stderr: %v", err)
	}

	go func() {
		stderrScanner := bufio.NewScanner(stderrReader)
		for stderrScanner.Scan() {
			text := stderrScanner.Text()
			m.log.Errorf("%s", text)
		}
	}()

	m.update = update

	if err := m.updateCmd.Start(); err != nil {
		m.update.State = StateFailed

		m.saveUpdate(m.update)
		m.notifyUpdateClients(m.update)

		return nil, err
	}

	go func() {
		err := m.updateCmd.Wait()
		if err != nil {
			if m.update.State != StateCancelled {
				// only set failure state if the process didn't stop
				// because of a user initiated cancellation
				m.update.State = StateFailed
			}
		} else {
			m.update.State = StateInstalled
			m.shouldReboot = true
		}

		m.saveUpdate(m.update)
		m.notifyUpdateClients(m.update)

		atomic.StoreInt32(&m.updating, 0)
	}()

	return m.update, nil
}

func (m *MenderUpdater) CancelUpdate(id string) (*Update, error) {
	if atomic.LoadInt32(&m.updating) != 1 {
		return nil, errors.New("no update in progress")
	}

	// set this state before killing the process so that the
	// handler can tell cancellations and failures apart
	m.update.State = StateCancelled

	err := m.updateCmd.Process.Kill()
	if err != nil {
		return nil, errors.Errorf("could not kill update: %v", err)
	}

	// reset progress to zero
	m.progress = 0

	m.saveUpdate(m.update)
	m.notifyUpdateClients(m.update)

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

func (m *MenderUpdater) notifyUpdateClients(update *Update) {
	if update == nil {
		return
	}

	for _, client := range m.clients {
		if update.Id == client.updateId {
			// add progress information if it's the currently running update
			if m.update != nil && update.Id == m.update.Id {
				update.Progress = m.progress
				update.ShouldReboot = m.shouldReboot
				update.ShouldCommit = m.shouldCommit
			}

			client.Update <- update
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

	m.saveUpdate(m.update)
	m.notifyUpdateClients(m.update)

	return m.update, nil
}

func (m *MenderUpdater) RejectUpdate(id string) (*Update, error) {
	_, err := exec.Command("mender", "-rollback").Output()
	if err != nil {
		return nil, errors.Errorf("Could not retrieve artifact name: %v", err)
	}

	m.update.State = StateRejected
	m.update.ShouldReboot = true

	m.saveUpdate(m.update)
	m.notifyUpdateClients(m.update)

	return m.update, nil
}

func (m *MenderUpdater) saveUpdate(update *Update) {
	err := m.db.SaveUpdate(&sweetdb.Update{
		Id:      update.Id,
		Started: update.Started,
		Url:     update.Url,
		State:   update.State,
	})
	if err != nil {
		m.log.Errorf("could not save update state: %v", err)
	}
}
