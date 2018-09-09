package machine

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"
	"gobot.io/x/gobot/drivers/gpio"
)

const (
	motorPin = "13"
	vibratorPin = "11"
	sensorPin = "7"
)

type Machine struct {
	// Public receiving channel for touch events
	TouchEvents <-chan bool
	// Internal robot instance
	robot *gobot.Robot
	// Internal raw touch events channel
	touchEvents chan bool
	// Internal motor events channel
	motorEvents chan bool
	// Internal vibrator events channel
	vibratorEvents chan bool
}

func NewMachine() *Machine {
	rpi := raspi.NewAdaptor()
	motorPin := gpio.NewDirectPinDriver(rpi, motorPin)
	vibratorPin := gpio.NewDirectPinDriver(rpi, vibratorPin)
	touchSensor := gpio.NewButtonDriver(rpi, sensorPin)

	touchEvents := make(chan bool)
	motorEvents := make(chan bool)
	vibratorEvents := make(chan bool)

	robot := gobot.NewRobot("dispenser",
		[]gobot.Connection{rpi},
		[]gobot.Device{motorPin, vibratorPin, touchSensor},
	)

	robot.Work = func() {
		defer touchSensor.DeleteEvent(gpio.ButtonPush)
		touchSensor.On(gpio.ButtonPush, func(data interface{}) {
			touchEvents <- true
		})

		defer touchSensor.DeleteEvent(gpio.ButtonRelease)
		touchSensor.On(gpio.ButtonRelease, func(data interface{}) {
			touchEvents <- false
		})

		for {
			select {
			case turnOn := <-motorEvents:
				if turnOn {
					motorPin.On()
				} else {
					motorPin.Off()
				}
			case turnOn := <-vibratorEvents:
				if turnOn {
					vibratorPin.On()
				} else {
					vibratorPin.Off()
				}
			}
		}
	}

	m := &Machine{
		TouchEvents: touchEvents,
		robot: robot,
		touchEvents: touchEvents,
		motorEvents: motorEvents,
		vibratorEvents: vibratorEvents,
	}

	return m
}

func (m *Machine) Start() {
	m.robot.Start()
}

func (m *Machine) Stop() {
	m.robot.Stop()
}

func (m *Machine) ToggleMotor(on bool) {
	m.motorEvents <- on
}

func (m *Machine) ToggleBuzzer(on bool) {
	m.vibratorEvents <- on
}