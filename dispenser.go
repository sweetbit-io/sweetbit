package main

import (
	"crypto/x509"
	"encoding/hex"
	"github.com/go-errors/errors"
	"github.com/lightningnetwork/lnd/lnrpc"
	log "github.com/sirupsen/logrus"
	"github.com/the-lightning-land/sweetd/machine"
	"github.com/the-lightning-land/sweetd/sweetdb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

type dispenser struct {
	machine          machine.Machine
	db               *sweetdb.DB
	dispenseOnTouch  bool
	buzzOnDispense   bool
	done             chan struct{}
	payments         chan *lnrpc.Invoice
	grpcConn         *grpc.ClientConn
	lightningNodeUri string
	// add bluetooth pairing
}

func newDispenser(machine machine.Machine, db *sweetdb.DB) *dispenser {
	return &dispenser{
		machine:         machine,
		db:              db,
		dispenseOnTouch: true,
		buzzOnDispense:  false,
		done:            make(chan struct{}),
		payments:        make(chan *lnrpc.Invoice),
	}
}

func (d *dispenser) run() error {
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
		err := d.connectLndNode(node.Uri, node.Cert, node.Macaroon)
		if err != nil {
			log.Errorf("Could not connect to remote lightning node: %v", err)
		}
	}

	for {
		select {
		case on := <-d.machine.TouchEvents():
			// react on direct touch events of the machine
			log.Info("Touch event {}", on)

			if d.dispenseOnTouch {
				d.machine.ToggleMotor(on)

				if d.buzzOnDispense {
					d.machine.ToggleBuzzer(on)
				}
			} else if !on {
				d.machine.ToggleBuzzer(false)
				d.machine.ToggleMotor(false)
			}

		case payment := <-d.payments:
			// react on incoming payments
			dispense := time.Duration(payment.Value/2) * time.Millisecond

			log.Debugf("Dispensing for a duration of %v", dispense)

			d.machine.ToggleMotor(true)

			if d.buzzOnDispense {
				d.machine.ToggleBuzzer(true)
			}

			time.Sleep(dispense)
			d.machine.ToggleMotor(false)

			if d.buzzOnDispense {
				d.machine.ToggleBuzzer(false)
			}

		case <-d.done:
			// finish loop when program is done
			return nil
		}
	}
}

func (d *dispenser) saveLndNode(uri string, certBytes []byte, macaroonBytes []byte) error {
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

func (d *dispenser) deleteLndNode() error {
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

func (d *dispenser) connectLndNode(uri string, certBytes []byte, macaroonBytes []byte) error {
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
	d.lightningNodeUri = uri

	go func() {
		log.Info("Listening to paid invoices...")

		for {
			invoice, err := invoices.Recv()
			if err == io.EOF {
				break
			}

			if err != nil && status.Code(err) == 1 {
				log.Info("Stopping invoice listener")
				break
			} else if err != nil {
				log.WithError(err).Error("Failed receiving subscription items")
				break
			}

			if invoice.Settled {
				log.Debugf("Received settled payment of %v sat", invoice.Value)
				d.payments <- invoice
			} else {
				log.Debugf("Generated invoice of %v sat", invoice.Value)
			}
		}
	}()

	return nil
}

func (d *dispenser) disconnectLndNode() error {
	log.Infof("Disconnecting from remote lightning node")

	// close open connection
	if d.grpcConn != nil {
		d.grpcConn.Close()
	}

	// remove currently connected node uri
	d.lightningNodeUri = ""

	return nil
}

func (d *dispenser) shutdown() {
	d.machine.Stop()

	if d.grpcConn != nil {
		d.grpcConn.Close()
	}

	close(d.done)
}
