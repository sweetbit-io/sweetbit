package dispenser

import (
	"github.com/go-errors/errors"
	"net"
	"net/http"
	"sync"
)

func (d *Dispenser) runApi(wg sync.WaitGroup) error {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		return errors.Errorf("unable to listen: %v", err)
	}

	// point the onion service to the listener
	d.apiOnionService.SetListener(listener)

	go func() {
		wg.Add(1)

		go func() {
			wg.Add(1)

			err := http.Serve(listener, d.apiHandler)
			if err != nil {
				d.log.Errorf("unable to serve: %v", err)
			}

			wg.Done()
		}()

		// subscribe to network updates
		networkClient := d.network.Subscribe()

		if d.network.Status().Connected() {
			d.apiOnionService.Start()
		}

		done := false

		d.log.Infof("start api")

		for !done {
			select {
			case update := <-networkClient.Updates:
				d.log.Infof("network changed to %v", update)

				if update.Connected {
					d.apiOnionService.Start()
				}
			case <-d.done:
				done = true
			}
		}

		if listener != nil {
			err = listener.Close()
			if err != nil {
				d.log.Errorf("could not close api listener: %v", err)
			}
		}

		networkClient.Cancel()

		d.apiOnionService.Stop()

		d.log.Infof("stopped api")

		wg.Done()
	}()

	return nil
}

func (d *Dispenser) GetApiOnionID() string {
	return d.apiOnionService.ID()
}
