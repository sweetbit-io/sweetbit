package dispenser

import (
	"github.com/go-errors/errors"
	"net"
	"net/http"
	"sync"
)

// runPos
func (d *Dispenser) runPos(wg sync.WaitGroup) error {
	listener, err := net.Listen("tcp", "127.0.0.1:9001")
	if err != nil {
		return errors.Errorf("unable to listen: %v", err)
	}

	// point the onion service to the listener
	d.posOnionService.SetListener(listener)

	go func() {
		wg.Add(1)

		go func() {
			wg.Add(1)

			err := http.Serve(listener, d.posHandler)
			if err != nil {
				d.log.Errorf("unable to serve: %v", err)
			}

			wg.Done()
		}()

		d.log.Infof("starting point of sales")

		// subscribe to network updates
		networkClient := d.network.Subscribe()

		if d.network.Status().Connected() {
			d.posOnionService.Start()
		}

		done := false

		d.log.Infof("start point of sales")

		for !done {
			select {
			case update := <-networkClient.Updates:
				d.log.Infof("point of sales network changed to %v", update)

				if update.Connected {
					d.posOnionService.Start()
				}
			case <-d.done:
				done = true
			}
		}

		networkClient.Cancel()

		d.posOnionService.Stop()

		d.log.Infof("stopped point of sales")

		wg.Done()
	}()

	return nil
}

func (d *Dispenser) GetPosOnionID() string {
	return d.posOnionService.ID()
}
