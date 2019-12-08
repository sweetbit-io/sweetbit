package main

import (
	"github.com/jessevdk/go-flags"
)

type raspberryConfig struct {
	TouchPin  string `long:"touchpin" description:"BCM number of the touch input pin."`
	MotorPin  string `long:"motorpin" description:"BCM number of the motor output pin."`
	BuzzerPin string `long:"buzzerpin" description:"BCM number of the buzzer output pin."`
}

type mockConfig struct {
	Listen string `long:"listen" description:"Add an interface/port to listen for mock touches."`
}

type torConfig struct {
	Path string `long:"path" description:"The path to the Tor binary."`
}

type profilingConfig struct {
	Listen string `long:"listen" description:"Add an interface/port to listen for profiling data."`
}

type networkConfig struct {
	Interface string `long:"ifname" description:"Wifi interface name."`
}

type pairingConfig struct {
	Interface string `long:"ifname" description:"Bluetooth interface name."`
}

type config struct {
	ShowVersion bool             `short:"v" long:"version" description:"Display version information and exit."`
	Debug       bool             `long:"debug" description:"Start in debug mode."`
	Machine     string           `long:"machine" description:"The machine controller to use." choice:"raspberry" choice:"mock"`
	Raspberry   *raspberryConfig `group:"Raspberry" namespace:"raspberry"`
	Mock        *mockConfig      `group:"Mock" namespace:"mock"`
	Net         *networkConfig   `group:"Network" namespace:"network"`
	Pairing     *pairingConfig   `group:"Pairing" namespace:"pairing"`
	DataDir     string           `long:"datadir" description:"The directory to store sweetd's data within.'"`
	Updater     string           `long:"updater" description:"The updater to use." choice:"none" choice:"mender"`
	Tor         *torConfig       `group:"Tor" namespace:"tor"`
	Profiling   *profilingConfig `group:"Profiling" namespace:"profiling"`
}

func loadConfig() (*config, error) {
	defaultCfg := config{
		Machine: "raspberry",
		Debug:   false,
		Raspberry: &raspberryConfig{
			TouchPin:  "25",
			MotorPin:  "23",
			BuzzerPin: "24",
		},
		Net:     nil,
		Pairing: nil,
		DataDir: "./data",
		Updater: "none",
		Tor: &torConfig{
			Path: "",
		},
	}

	preCfg := defaultCfg

	if _, err := flags.Parse(&preCfg); err != nil {
		return nil, err
	}

	cfg := preCfg

	return &cfg, nil
}
