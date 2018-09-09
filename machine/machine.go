package machine

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"periph.io/x/periph/host"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/gpio"
	"time"
)

const (
	touchPin  = "4"  // BCM 4 = #7
	motorPin  = "27" // BCM 27 = #13
	buzzerPin = "17" // BCM 17 = #11
)

type Machine struct {
	// Public receiving channel for touch events
	TouchEvents <-chan bool
	// Internal sending channel for touch events
	touchEvents chan<- bool
	// Internal motor events channel
	motorEvents chan bool
	// Internal buzzer events channel
	buzzerEvents chan bool
	// Internal done channel
	done chan bool
	// Internal goroutine WaitGroup
	waitGroup sync.WaitGroup
}

func NewMachine() *Machine {
	touchEvents := make(chan bool)
	motorEvents := make(chan bool)
	buzzerEvents := make(chan bool)

	m := &Machine{
		TouchEvents:  touchEvents,
		touchEvents:  touchEvents,
		motorEvents:  motorEvents,
		buzzerEvents: buzzerEvents,
	}

	return m
}

func (m *Machine) Start() {
	log.Info("Starting machine")

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	m.done = make(chan bool)

	var waitGroup sync.WaitGroup
	m.waitGroup = waitGroup

	go m.handleTouch()
	go m.driveMotor()
	go m.driveBuzzer()
}

func (m *Machine) Stop() {
	log.Info("Stopping machine")

	close(m.done)

	// Blocking until all goroutines finished executing
	m.waitGroup.Wait()

	log.Info("Machine stopped")
}

func (m *Machine) ToggleMotor(on bool) {
	log.Info("Toggling motor {}", on)
	m.motorEvents <- on
}

func (m *Machine) ToggleBuzzer(on bool) {
	log.Info("Toggling buzzer {}", on)
	m.buzzerEvents <- on
}

func (m *Machine) handleTouch() {
	log.Info("Starting to handle touch events")

	m.waitGroup.Add(1)
	defer m.waitGroup.Done()

	p := gpioreg.ByName(touchPin)

	// set as input, with an internal pull down resistor
	if err := p.In(gpio.PullDown, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}

	// Turn blocking WaitForEdge() func into channel
	edges := make(chan bool)
	go func() {
		// m.waitGroup.Add(1)
		// defer m.waitGroup.Done()

		// TODO: Stop this goroutine on done signal

		var notifyAfterThrottledTime = time.Time{}
		var hasSentHigh = false

		for {
			p.WaitForEdge(-1)

			if p.Read() == gpio.High {
				if notifyAfterThrottledTime.IsZero() {
					// just save time for throttling
					notifyAfterThrottledTime = time.Now().Add(20 * time.Millisecond)
				} else if !hasSentHigh && time.Now().After(notifyAfterThrottledTime) {
					// send throttled touch start
					edges <- true

					// make sure next high signal isn't sent anymore
					hasSentHigh = true
				}
			} else {
				// reset time for throttling
				notifyAfterThrottledTime = time.Time{}

				// send touch stop
				edges <- false

				// make sure the next high is sent
				hasSentHigh = false
			}
		}
	}()

	for {
		select {
		case touch := <-edges:
			log.Info("Received touch event", touch)
			m.touchEvents <- touch
		case <-m.done:
			log.Info("Got done event in handleTouch")

			return
		}
	}

	log.Info("Leaving handleTouch goroutine")
}

func (m *Machine) driveMotor() {
	log.Info("Starting to handle motor events")

	m.waitGroup.Add(1)
	defer m.waitGroup.Done()

	p := gpioreg.ByName(motorPin)

	for {
		select {
		case on := <-m.motorEvents:
			log.Info("Driving motor", on)

			if on {
				p.Out(gpio.High)
			} else {
				p.Out(gpio.Low)
			}
		case <-m.done:
			log.Info("Got done event in driveMotor")

			p.Out(gpio.Low)
			return
		}
	}

	log.Info("Leaving driveMotor goroutine")
}

func (m *Machine) driveBuzzer() {
	log.Info("Starting to handle buzzer events")

	m.waitGroup.Add(1)
	defer m.waitGroup.Done()

	p := gpioreg.ByName(buzzerPin)

	for {
		select {
		case on := <-m.buzzerEvents:
			log.Info("Driving buzzer", on)
			if on {
				p.Out(gpio.High)
			} else {
				p.Out(gpio.Low)
			}
		case <-m.done:
			log.Info("Got done event in driveBuzzer")

			p.Out(gpio.Low)
			return
		}
	}

	log.Info("Leaving driveBuzzer goroutine")
}
