package dispenser

import (
	"crypto/x509"
	"encoding/hex"
	"github.com/go-errors/errors"
	"github.com/lightningnetwork/lnd/lnrpc"
	log "github.com/sirupsen/logrus"
	"github.com/the-lightning-land/sweetd/ap"
	"github.com/the-lightning-land/sweetd/machine"
	"github.com/the-lightning-land/sweetd/sweetdb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"strings"
	"sync"
	"time"
)

type Dispenser struct {
	machine              machine.Machine
	AccessPoint          ap.Ap
	db                   *sweetdb.DB
	DispenseOnTouch      bool
	BuzzOnDispense       bool
	done                 chan struct{}
	payments             chan *lnrpc.Invoice
	grpcConn             *grpc.ClientConn
	LightningNodeUri     string
	dispenses            chan bool
	dispenseClients      map[uint32]*DispenseClient
	dispenseClientMtx    sync.Mutex
	nextDispenseClientID uint32
	memoPrefix           string
}

type DispenseClient struct {
	Dispenses  chan bool
	Id         uint32
	cancelChan chan struct{}
	dispenser  *Dispenser
}

type Config struct {
	Machine     machine.Machine
	AccessPoint ap.Ap
	DB          *sweetdb.DB
	MemoPrefix  string
}

func NewDispenser(config *Config) *Dispenser {
	return &Dispenser{
		machine:         config.Machine,
		AccessPoint:     config.AccessPoint,
		db:              config.DB,
		DispenseOnTouch: true,
		BuzzOnDispense:  false,
		done:            make(chan struct{}),
		payments:        make(chan *lnrpc.Invoice),
		dispenses:       make(chan bool),
		dispenseClients: make(map[uint32]*DispenseClient),
		memoPrefix:      config.MemoPrefix,
	}
}

func (d *Dispenser) Run() error {
	log.Info("Starting machine...")

	if err := d.machine.Start(); err != nil {
		return errors.Errorf("Could not start machine: %v", err)
	}

	// Signal successful startup with two short buzzer noises
	d.machine.DiagnosticNoise()

	node, err := d.db.GetLightningNode()
	if err != nil {
		return err
	}

	// connect to remote lightning node
	if node != nil {
		err := d.ConnectLndNode(node.Uri, node.Cert, node.Macaroon)
		if err != nil {
			log.Errorf("Could not connect to remote lightning node: %v", err)
		}
	}

	// Notify all subscribed dispense clients
	go func() {
		for {
			on := <-d.dispenses

			for _, client := range d.dispenseClients {
				client.Dispenses <- on
			}
		}
	}()

	for {
		select {
		case on := <-d.machine.TouchEvents():
			// react on direct touch events of the machine
			log.Infof("Touch event %v", on)

			if d.DispenseOnTouch && on {
				d.ToggleDispense(true)
			} else {
				d.ToggleDispense(false)
			}

		case <-d.payments:
			// react on incoming payments
			dispense := 1500 * time.Millisecond

			log.Debugf("Dispensing for a duration of %v", dispense)

			d.ToggleDispense(true)
			time.Sleep(dispense)
			d.ToggleDispense(false)

		case <-d.done:
			// finish loop when program is done
			return nil
		}
	}
}

func (d *Dispenser) ToggleDispense(on bool) {
	// Always make sure that buzzing stops
	if d.BuzzOnDispense || !on {
		d.machine.ToggleBuzzer(on)
	}

	d.machine.ToggleMotor(on)

	d.dispenses <- on
}

func (d *Dispenser) SaveLndNode(uri string, certBytes []byte, macaroonBytes []byte) error {
	err := d.db.SetLightningNode(&sweetdb.LightningNode{
		Uri:      uri,
		Cert:     certBytes,
		Macaroon: macaroonBytes,
	})

	if err != nil {
		return errors.Errorf("Couldn not save lnd node connection: %v", err)
	}

	return nil
}

func (d *Dispenser) DeleteLndNode() error {
	err := d.db.SetLightningNode(nil)

	if err != nil {
		return errors.Errorf("Couldn not delete lnd node connection: %v", err)
	}

	return nil
}

var (
	beginCertificateBlock = []byte("-----BEGIN CERTIFICATE-----\n")
	endCertificateBlock   = []byte("\n-----END CERTIFICATE-----")
)

func (d *Dispenser) ConnectLndNode(uri string, certBytes []byte, macaroonBytes []byte) error {
	log.Infof("Connecting to remote lightning node %s", uri)

	cert := x509.NewCertPool()

	fullCertBytes := append(beginCertificateBlock, certBytes...)
	fullCertBytes = append(fullCertBytes, endCertificateBlock...)

	if ok := cert.AppendCertsFromPEM(fullCertBytes); !ok {
		return errors.New("Could not parse tls cert.")
	}

	creds := credentials.NewClientTLSFromCert(cert, "")

	conn, err := grpc.Dial(uri, grpc.WithTransportCredentials(creds))
	if err != nil {
		return errors.Errorf("Could not connect to lightning node: %v", err)
	}

	client := lnrpc.NewLightningClient(conn)

	hexMacaroon := hex.EncodeToString(macaroonBytes)

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("macaroon", hexMacaroon))

	log.Info("Subscribing to invoices...")

	invoices, err := client.SubscribeInvoices(ctx, &lnrpc.InvoiceSubscription{})
	if err != nil {
		return errors.Errorf("Could not subscribe to invoices: %v", err)
	}

	log.Info("Connected to lightning node.")

	// close any previous connections
	if d.grpcConn != nil {
		d.grpcConn.Close()
	}

	// assign new connection
	d.grpcConn = conn

	// save currently connected node uri
	d.LightningNodeUri = uri

	go func() {
		log.Info("Listening to paid invoices...")

		for {
			invoice, err := invoices.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				errStatus, ok := status.FromError(err)
				if !ok {
					log.Errorf("Could not get status from err: %v", err)
				}

				if errStatus.Code() == 1 {
					log.Info("Stopping invoice listener")
					break
				} else if err != nil {
					log.WithError(err).Error("Failed receiving subscription items")
					break
				}
			}

			if invoice.Settled {
				if d.memoPrefix == "" ||
					(d.memoPrefix != "" && strings.HasPrefix(invoice.Memo, d.memoPrefix)) {
					log.Debugf("Received settled payment of %v sat", invoice.Value)
					d.payments <- invoice
				} else {
					log.Infof("Received payment with memo %s but memo prefix is %s.", invoice.Memo, d.memoPrefix)
				}
			} else {
				log.Debugf("Generated invoice of %v sat", invoice.Value)
			}
		}
	}()

	return nil
}

func (d *Dispenser) DisconnectLndNode() error {
	log.Infof("Disconnecting from remote lightning node")

	// close open connection
	if d.grpcConn != nil {
		d.grpcConn.Close()
	}

	// remove currently connected node uri
	d.LightningNodeUri = ""

	return nil
}

func (d *Dispenser) SetWifiConnection(connection *sweetdb.WifiConnection) error {
	log.Infof("Setting Wifi connection")

	err := d.db.SetWifiConnection(connection)
	if err != nil {
		return errors.Errorf("Failed setting Wifi connection: %v", err)
	}

	return nil
}

func (d *Dispenser) GetName() (string, error) {
	log.Infof("Getting name")

	name, err := d.db.GetName()
	if err != nil {
		return "", errors.Errorf("Failed getting name: %v", err)
	}

	return name, nil
}

func (d *Dispenser) SetName(name string) error {
	log.Infof("Setting name")

	err := d.db.SetName(name)
	if err != nil {
		return errors.Errorf("Failed setting name: %v", err)
	}

	return nil
}

func (d *Dispenser) SetDispenseOnTouch(dispenseOnTouch bool) error {
	log.Infof("Setting dispense on touch")

	d.DispenseOnTouch = dispenseOnTouch

	err := d.db.SetDispenseOnTouch(dispenseOnTouch)
	if err != nil {
		return errors.Errorf("Failed setting dispense on touch: %v", err)
	}

	return nil
}

func (d *Dispenser) SetBuzzOnDispense(buzzOnDispense bool) error {
	log.Infof("Setting buzz on dispense")

	d.BuzzOnDispense = buzzOnDispense

	err := d.db.SetBuzzOnDispense(buzzOnDispense)
	if err != nil {
		return errors.Errorf("Failed setting buzz on dispense: %v", err)
	}

	return nil
}

func (d *Dispenser) ConnectToWifi(ssid string, psk string) error {
	log.Infof("Connecting to wifi %v", ssid)

	err := d.AccessPoint.ConnectWifi(ssid, psk)
	if err != nil {
		log.Errorf("Could not get Wifi networks: %v", err)
		return errors.New("Could not get Wifi networks")
	}

	err = d.SetWifiConnection(&sweetdb.WifiConnection{
		Ssid: ssid,
		Psk:  psk,
	})
	if err != nil {
		log.Errorf("Could not save wifi connection: %v", err)
	}

	return nil
}

func (d *Dispenser) Shutdown() {
	d.machine.Stop()

	if d.grpcConn != nil {
		d.grpcConn.Close()
	}

	close(d.done)
}

func (d *Dispenser) SubscribeDispenses() *DispenseClient {
	client := &DispenseClient{
		Dispenses:  make(chan bool),
		cancelChan: make(chan struct{}),
		dispenser:  d,
	}

	d.dispenseClientMtx.Lock()
	client.Id = d.nextDispenseClientID
	d.nextDispenseClientID++
	d.dispenseClientMtx.Unlock()

	d.dispenseClients[client.Id] = client

	return client
}

func (c *DispenseClient) Cancel() {
	delete(c.dispenser.dispenseClients, c.Id)

	close(c.cancelChan)
}
