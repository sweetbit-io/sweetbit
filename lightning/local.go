package lightning

import (
	"bufio"
	"github.com/go-errors/errors"
	"os/exec"
)

//var sizeRegexp = regexp.MustCompile("BTCN: Processed 33870 blocks in the last 10.43s (height 267301, 2013-11-01 15:49:51 +0100 CET)")
//var sizeRegexp = regexp.MustCompile("BTCN: Verified 17000 filter headers in the last 10.13s (height 284001, 2014-02-03 20:51:30 +0100 CET)")

type LocalNodeConfig struct {
	DataDir string
	Logger  Logger
}

type LocalNode struct {
	*LndNode
	dataDir string
	log     Logger
	version string
	cmd     *exec.Cmd
}

func NewLocalNode(config *LocalNodeConfig) (*LocalNode, error) {
	var log Logger
	if config.Logger != nil {
		log = config.Logger
	} else {
		log = noopLogger{}
	}

	_, err := exec.LookPath("lnd")
	if err != nil {
		return nil, errors.New("lnd is not installed or missing in $PATH")
	}

	versionBytes, err := exec.Command("lnd", "--version").Output()
	if err != nil {
		return nil, errors.Errorf("unable to get version: %v", err)
	}

	version := string(versionBytes)

	log.Infof("using lnd version %s", version)

	lndNode, err := NewLndNode(&LndNodeConfig{
		Logger: log,
	})
	if err != nil {
		return nil, errors.Errorf("unable to create lnd node: %v", err)
	}

	return &LocalNode{
		LndNode: lndNode,
		dataDir: config.DataDir,
		log:     log,
		version: version,
	}, nil
}

func (n *LocalNode) Start() error {
	args := []string{}

	args = append(args, "--lnddir", n.dataDir)
	args = append(args, "--bitcoin.active")
	args = append(args, "--bitcoin.mainnet")
	args = append(args, "--bitcoin.node", "neutrino")
	args = append(args, "--neutrino.connect", "mainnet1-btcd.zaphq.io")

	n.cmd = exec.Command("lnd", args...)

	stdoutReader, err := n.cmd.StdoutPipe()
	if err != nil {
		return errors.Errorf("unable to get stdout reader: %v", err)
	}

	go func() {
		stdoutScanner := bufio.NewScanner(stdoutReader)
		for stdoutScanner.Scan() {
			text := stdoutScanner.Text()
			n.log.Debugf("%s", text)
		}

		n.log.Debugf("left stdin reader")
	}()

	stderrReader, err := n.cmd.StderrPipe()
	if err != nil {
		return errors.Errorf("unable to get stderr reader: %v", err)
	}

	go func() {
		stderrScanner := bufio.NewScanner(stderrReader)
		for stderrScanner.Scan() {
			text := stderrScanner.Text()
			n.log.Errorf("%s", text)
		}

		n.log.Debugf("left stderr reader")
	}()

	err = n.cmd.Start()
	if err != nil {
		n.log.Errorf("could not start: %v", err)
	}

	go func() {
		err := n.cmd.Wait()
		if err != nil {
			n.log.Errorf("exited with error: %v", err)
		} else {
			n.log.Infof("exited successfully")
		}
	}()

	return nil
}

func (n *LocalNode) Stop() error {
	return nil
}
