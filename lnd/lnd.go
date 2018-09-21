package lnd

import (
	"os/exec"
	log "github.com/sirupsen/logrus"
)

type Lnd struct {
	cmd    *exec.Cmd
	config LndConfig
}

type LndConfig struct {
	TLSExtraIP           string
	TLSExtraDomain       string
	NeutrinoConnectPeers []string
	Alias                string
	Color                string
	NAT                  bool
	ExternalIPs          []string
	RPCListeners         []string
	RESTListeners        []string
	Listeners            []string
}

func NewLnd(config LndConfig) *Lnd {
	l := &Lnd{
		config: config,
	}
	return l
}

func (l *Lnd) Start() error {
	l.cmd = exec.Command("lnd")

	err := l.cmd.Run()
	if err != nil {
		log.WithError(err).Error("Error with the browser process")
	}

	return err
}

func (l *Lnd) Stop() error {
	err := l.cmd.Process.Release()
	if err != nil {
		log.WithError(err).Fatal("Error shutting down lnd")
	}

	return err
}

func (l *Lnd) Kill() error {
	err := l.cmd.Process.Kill()
	if err != nil {
		log.WithError(err).Fatal("Error killing lnd")
	}

	return err
}
