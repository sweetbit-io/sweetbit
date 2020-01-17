package lightning

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/the-lightning-land/sweetd/onion"
	"io/ioutil"
	"net"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

//var sizeRegexp = regexp.MustCompile("BTCN: Processed 33870 blocks in the last 10.43s (height 267301, 2013-11-01 15:49:51 +0100 CET)")
//var sizeRegexp = regexp.MustCompile("BTCN: Verified 17000 filter headers in the last 10.13s (height 284001, 2014-02-03 20:51:30 +0100 CET)")
//var sizeRegexp = regexp.MustCompile("BTCN: Verified 17000 filter headers in the last 10.13s (height 284001, 2014-02-03 20:51:30 +0100 CET)")
// BTCN: Fully caught up with cfheaders at height 610208, waiting at tip for new blocks
var grpcPortRegexp = regexp.MustCompile(`RPCS: password RPC server listening on 127\.0\.0\.1:(\d+)`)
var rpcPortRegexp = regexp.MustCompile(`RPCS: password gRPC proxy started at 127\.0\.0\.1:(\d+)`)
var walletOpenedRegexp = regexp.MustCompile(`LNWL: Opened wallet`)

//var sizeRegexp = regexp.MustCompile("LTND: Waiting for wallet encryption password")
//var sizeRegexp = regexp.MustCompile("RPCS: Done generating TLS certificates")

type LocalNodeConfig struct {
	DataDir  string
	Logger   Logger
	OnionSvc *onion.Service
}

type LocalNode struct {
	*LndNode
	dataDir       string
	log           Logger
	version       string
	cmd           *exec.Cmd
	grpcPort      int
	rpcPort       int
	onionSvc      *onion.Service
	cert          string
	adminMacaroon string
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
		LndNode:  lndNode,
		dataDir:  config.DataDir,
		log:      log,
		version:  version,
		onionSvc: config.OnionSvc,
	}, nil
}

func (n *LocalNode) Start() error {
	var args []string

	args = append(args, "--lnddir", n.dataDir)
	args = append(args, "--bitcoin.active")
	args = append(args, "--bitcoin.mainnet")
	args = append(args, "--bitcoin.node", "neutrino")
	args = append(args, "--neutrino.connect", "btcd-mainnet.lightning.computer:8333")
	args = append(args, "--tlsextradomain", n.Uri())
	//args = append(args, "--listen", "127.0.0.1:0")
	//args = append(args, "--rpclisten", "127.0.0.1:0")
	//args = append(args, "--restlisten", "0.0.0.0:8080")
	// TODO(davidknezic) add tor support

	n.cmd = exec.Command("lnd", args...)

	stdoutReader, err := n.cmd.StdoutPipe()
	if err != nil {
		return errors.Errorf("unable to get stdout reader: %v", err)
	}

	listenChan := make(chan struct{})

	go func() {
		stdoutScanner := bufio.NewScanner(stdoutReader)
		for stdoutScanner.Scan() {
			text := stdoutScanner.Text()
			n.log.Debugf("%s", text)

			if matches := grpcPortRegexp.FindStringSubmatch(text); len(matches) == 2 {
				port, err := strconv.ParseInt(matches[1], 10, 16)
				if err != nil {
					n.log.Errorf("Could not parse port: %v", err)
					continue
				}

				n.grpcPort = int(port)

				n.log.Infof("grpc listens on port %d", n.grpcPort)

				if n.grpcPort > 0 && n.rpcPort > 0 {
					close(listenChan)
				}
			}

			if matches := rpcPortRegexp.FindStringSubmatch(text); len(matches) == 2 {
				port, err := strconv.ParseInt(matches[1], 10, 16)
				if err != nil {
					n.log.Errorf("Could not parse port: %v", err)
					continue
				}

				n.rpcPort = int(port)

				n.log.Infof("rpc listens on port %d", n.rpcPort)

				if n.grpcPort > 0 && n.rpcPort > 0 {
					close(listenChan)
				}
			}

			if matches := walletOpenedRegexp.FindStringSubmatch(text); len(matches) == 1 {
				adminMacaroonBytes, err := ioutil.ReadFile(filepath.Join(n.dataDir, "data/chain/bitcoin/mainnet/admin.macaroon"))
				if err != nil {
					n.log.Errorf("unable to read macaroon: %v", err)
				} else {
					n.setMacaroon(adminMacaroonBytes)
					n.updateStatus(StatusStarted)
				}
			}
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
		return errors.Errorf("unable to start: %v", err)
	}

	go func() {
		err := n.cmd.Wait()
		if err != nil {
			n.log.Errorf("exited with error: %v", err)
		} else {
			n.log.Infof("exited successfully")
		}
	}()

	// only continue when server is listening and ports are determined
	<-listenChan

	certBytes, err := ioutil.ReadFile(filepath.Join(n.dataDir, "tls.cert"))
	err = n.setTlsCredentials(certBytes, true)
	if err != nil {
		return errors.Errorf("unable to set certificate: %v", err)
	}

	var inlineCert string
	certLines := strings.Split(string(certBytes), "\n")
	for _, line := range certLines {
		if line != "" && !strings.Contains(line, "-") {
			inlineCert += line
		}
	}

	n.cert = inlineCert

	n.setUri(fmt.Sprintf("localhost:%d", n.grpcPort))

	adminMacaroonBytes, err := ioutil.ReadFile(filepath.Join(n.dataDir, "data/chain/bitcoin/mainnet/admin.macaroon"))
	n.setMacaroon(adminMacaroonBytes)
	n.adminMacaroon = base64.StdEncoding.EncodeToString(adminMacaroonBytes)

	err = n.LndNode.Start()
	if err != nil {
		return errors.Errorf("unable to start lnd node: %v", err)
	}

	n.onionSvc.SetListener(noopListener{addr: &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8080,
	}})
	n.onionSvc.Start()

	return nil
}

func (n *LocalNode) Stop() error {
	n.onionSvc.Stop()

	err := n.LndNode.Stop()
	if err != nil {
		return errors.Errorf("unable to stop: %v", err)
	}

	if n.cmd.Process != nil {
		err = n.cmd.Process.Kill()
		if err != nil {
			return errors.Errorf("unable to kill process: %v", err)
		}
	}

	return nil
}

func (n *LocalNode) Init(password string, mnemonic []string) error {
	err := n.LndNode.Init(password, mnemonic)
	if err != nil {
		return errors.Errorf("unable to create wallet: %v", err)
	}

	return nil
}

func (n *LocalNode) Unlock(walletPassword string) error {
	err := n.LndNode.Unlock(walletPassword)
	if err != nil {
		return errors.Errorf("unable to unlock wallet: %v", err)
	}

	return nil
}

func (n *LocalNode) Uri() string {
	return fmt.Sprintf("%s.onion:8080", n.onionSvc.ID())
}

func (n *LocalNode) Cert() string {
	return n.cert
}

func (n *LocalNode) AdminMacaroon() string {
	return n.adminMacaroon
}
