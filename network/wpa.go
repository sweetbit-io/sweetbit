package network

import (
	"github.com/go-errors/errors"
	"github.com/the-lightning-land/sweetd/network/wpa"
	"sync"
)

// check WpaNetworks compliance to its interface during compile time
var _ Network = (*WpaNetwork)(nil)

type Config struct {
	Interface string
	Logger    Logger
}

type nextClient struct {
	sync.Mutex
	id uint32
}

type WpaNetwork struct {
	log        Logger
	wpa        *wpa.Wpa
	ifname     string
	iface      *wpa.Interface
	clients    map[uint32]*Client
	nextClient nextClient
	done       chan struct{}
}

func NewWpaNetwork(config *Config) *WpaNetwork {
	net := &WpaNetwork{
		ifname:  config.Interface,
		wpa:     wpa.New(),
		clients: make(map[uint32]*Client),
	}

	if config.Logger != nil {
		net.log = config.Logger
	} else {
		net.log = noopLogger{}
	}

	return net
}

func (n *WpaNetwork) Start() error {
	err := n.wpa.Start()
	if err != nil {
		return errors.Errorf("could not start wpa: %v", err)
	}

	iface, err := n.wpa.GetInterface(n.ifname)
	if err != nil {
		_ = n.Stop()
		return errors.Errorf("could not find interface %v: %v", n.ifname, err)
	}

	n.iface = iface

	client, err := iface.State()
	if err != nil {
		return errors.Errorf("unable to subscribe to states: %v", err)
	}

	// initialize a new done channel to be closed to stop
	n.done = make(chan struct{})

	go func() {
		done := false

		for !done {
			select {
			case on := <-client.State:
				for _, client := range n.clients {
					client.Updates <- &Connectivity{
						Connected: on,
						Ip:        "",
						Ssid:      "",
					}
				}
			case <-n.done:
				// finish loop when program is done
				done = true
			}
		}

		client.Cancel()
	}()

	return nil
}

func (n *WpaNetwork) Stop() error {
	// signal the dispenser run loop to stop
	close(n.done)

	err := n.wpa.Stop()
	if err != nil {
		return errors.Errorf("could not stop wpa: %v", err)
	}

	return nil
}

func (n *WpaNetwork) Status() *Status {
	return &Status{
	}
}

func (n *WpaNetwork) Connect(connection Connection) error {
	var err error
	var net *wpa.Network

	err = n.iface.RemoveAllNetworks()
	if err != nil {
		return errors.Errorf("unable to remove all networks: %v", err)
	}

	switch conn := connection.(type) {
	case *WpaConnection:
		net, err = n.iface.AddWpaNetwork(conn.Ssid)
	case *WpaPersonalConnection:
		net, err = n.iface.AddWpaPersonalNetwork(conn.Ssid, conn.Psk)
	case *WpaEnterpriseConnection:
		net, err = n.iface.AddWpaEnterpriseNetwork(conn.Ssid, conn.Identity, conn.Password)
	default:
		return errors.Errorf("unknown connection type %T provided", connection)
	}

	if err != nil {
		return errors.Errorf("unable to add network: %v", err)
	}

	err = net.Enable()
	if err != nil {
		return errors.Errorf("unable to enable network: %v", err)
	}

	return nil
}

func (n *WpaNetwork) Scan() (*ScanClient, error) {
	err := n.iface.Scan()
	if err != nil {
		return nil, errors.Errorf("unable to scan: %v", err)
	}

	wifisChan := make(chan *Wifi)

	client, err := n.iface.BSSAdded()
	if err != nil {
		return nil, errors.Errorf("unable to listen for added wifis: %v", err)
	}

	doneClient, err := n.iface.ScanDone()
	if err != nil {
		return nil, errors.Errorf("unable to listen to scan completion: %v", err)
	}

	bsss, err := n.iface.BSSs()
	if err != nil {
		return nil, errors.Errorf("unable to get BSSs: %v", err)
	}

	go func() {
		for _, bss := range bsss {
			b, err := bss.GetAll()
			if err != nil {
				continue
			}

			if b.Ssid == "" {
				// skip stations with empty ssid
				continue
			}

			wifi := &Wifi{
				Ssid: b.Ssid,
			}

			if b.WpaType == wpa.WpaPersonal {
				wifi.Encryption = EncryptionPersonal
			} else if b.WpaType == wpa.WpaEnterprise {
				wifi.Encryption = EncryptionEnterprise
			} else {
				wifi.Encryption = EncryptionNone
			}

			wifisChan <- wifi
		}

		for {
			select {
			case bss, ok := <-client.BSSAdded:
				if !ok {
					close(wifisChan)
					return
				}

				b, err := bss.GetAll()
				if err != nil {
					continue
				}

				if b.Ssid == "" {
					// skip stations with empty ssid
					continue
				}

				wifi := &Wifi{
					Ssid: b.Ssid,
				}

				if b.WpaType == wpa.WpaPersonal {
					wifi.Encryption = EncryptionPersonal
				} else if b.WpaType == wpa.WpaEnterprise {
					wifi.Encryption = EncryptionEnterprise
				} else {
					wifi.Encryption = EncryptionNone
				}

				wifisChan <- wifi
			case done, ok := <-doneClient.ScanDone:
				if !ok {
					close(wifisChan)
					return
				}

				if done {
					client.Cancel()
					doneClient.Cancel()
				}
			}
		}
	}()

	return &ScanClient{
		Wifis: wifisChan,
		Cancel: func() {
			client.Cancel()
			doneClient.Cancel()
		},
	}, nil
}

func (n *WpaNetwork) Subscribe() *Client {
	client := &Client{
		Updates:    make(chan *Connectivity),
		cancelChan: make(chan struct{}),
		network:    n,
	}

	n.nextClient.Lock()
	client.Id = n.nextClient.id
	n.nextClient.id++
	n.nextClient.Unlock()

	n.clients[client.Id] = client

	return client
}

func (n *WpaNetwork) deleteClient(id uint32) {
	delete(n.clients, id)
}
