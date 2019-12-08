package dispenser

import (
	"github.com/go-errors/errors"
	"github.com/the-lightning-land/sweetd/lightning"
	"github.com/the-lightning-land/sweetd/nodeman"
	"sync"
)

// runLightningNodes
func (d *Dispenser) runLightningNodes(wg sync.WaitGroup) {
	wg.Add(1)

	d.nodeman.Load()

	d.log.Infof("restored %d lightning nodes from database", len(d.nodeman.GetNodes()))

	// subscribe to network updates
	networkClient := d.network.Subscribe()

	if d.network.Status().Connected() {
		d.startLightningNodes()
	}

	done := false

	d.log.Infof("start running lightning nodes")

	for !done {
		select {
		case update := <-networkClient.Updates:
			d.log.Infof("Network changed to %v", update)

			if update.Connected {
				d.startLightningNodes()
			}
		case <-d.done:
			done = true
		}
	}

	networkClient.Cancel()

	for _, node := range d.nodeman.GetNodes() {
		err := node.Stop()
		if err != nil {
			d.log.Errorf("could not stop node %v", node)
		}
	}

	d.log.Infof("stopped running lightning nodes")

	wg.Done()
}

func (d *Dispenser) startLightningNodes() {
	for _, node := range d.nodeman.GetNodes() {
		if node.Enabled() {
			err := node.Start()
			if err != nil {
				d.log.Errorf("could not start node %v", node)
				continue
			}

			client, err := node.SubscribeInvoices()
			if err != nil {
				d.log.Errorf("could not subscribe to invoices: %v", err)
			}

			go d.handleLightningNodeInvoices(client)
		}
	}
}

func (d *Dispenser) handleLightningNodeInvoices(client *lightning.InvoicesClient) {
	d.log.Infof("start handling lightning invoices")

	for {
		invoice, ok := <-client.Invoices
		if !ok {
			d.log.Infof("stopped handling lightning invoices")
			break
		}

		if invoice.Settled {
			d.payments <- invoice
		}
	}
}

func (d *Dispenser) GetNodes() []nodeman.LightningNode {
	return d.nodeman.GetNodes()
}

func (d *Dispenser) GetNode(id string) nodeman.LightningNode {
	return d.nodeman.GetNode(id)
}

func (d *Dispenser) AddNode(config nodeman.NodeConfig) (nodeman.LightningNode, error) {
	return d.nodeman.AddNode(config)
}

func (d *Dispenser) RemoveNode(id string) error {
	return d.nodeman.RemoveNode(id)
}

func (d *Dispenser) EnableNode(id string) error {
	node := d.nodeman.GetNode(id)

	if node == nil {
		return errors.Errorf("unable to find node %s", id)
	}

	err := node.Start()
	if err != nil {
		return errors.Errorf("unable to start node %v", err)
	}

	client, err := node.SubscribeInvoices()
	if err != nil {
		return errors.Errorf("unable to subscribe to invoices: %v", err)
	}

	go d.handleLightningNodeInvoices(client)

	return d.nodeman.EnableNode(id)
}

func (d *Dispenser) DisableNode(id string) error {
	node := d.nodeman.GetNode(id)

	if node == nil {
		return errors.Errorf("unable to find node %s", id)
	}

	err := node.Stop()
	if err != nil {
		return errors.Errorf("unable to stop node %v", node)
	}

	return d.nodeman.DisableNode(id)
}

func (d *Dispenser) RenameNode(id string, name string) error {
	return d.nodeman.RenameNode(id, name)
}
