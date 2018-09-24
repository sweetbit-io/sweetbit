package main

import (
	"github.com/jessevdk/go-flags"
	"net"
	"fmt"
)

const (
	defaultRPCPort = 9000
)

type apConfig struct {
	Ip         string `long:"ip" description:"IP address of device on access point interface."`
	Interface  string `long:"interface" description:"Name of the access point interface."`
	Ssid       string `long:"ssid" description:"Name of the access point."`
	Passphrase string `long:"passphrase" description:"WPA Passphrase (expected 8..63)."`
	DhcpRange  string `long:"dhcprange" description:"IP range of the DHCP service."`
}

type raspberryConfig struct {
	TouchPin  string `long:"touchpin" description:"The touch input pin."`
	MotorPin  string `long:"motorpin" description:"The motor output pin."`
	BuzzerPin string `long:"buzzerpin" description:"The buzzer output pin."`
}

type mockConfig struct {
	Listen string `long:"listen" description:"Add an interface/port to listen for mock touches."`
}

type config struct {
	ShowVersion  bool             `short:"v" long:"version" description:"Display version information and exit."`
	RawListeners []string         `long:"listen" description:"Add an interface/port/socket to listen for RPC connections"`
	Listeners    []net.Addr
	Machine      string           `long:"machine" description:"The machine controller to use." choice:"raspberry" choice:"mock"`
	Raspberry    *raspberryConfig `group:"Raspberry" namespace:"raspberry"`
	Mock         *mockConfig      `group:"Mock" namespace:"mock"`
	RunAp        bool             `long:"ap" description:"Run the access point service."`
	Ap           *apConfig        `group:"Access point" namespace:"ap"`
}

func loadConfig() (*config, error) {
	defaultCfg := config{
		Machine: "raspberry",
		Raspberry: &raspberryConfig{
			TouchPin:  "4",
			MotorPin:  "27",
			BuzzerPin: "17",
		},
		RunAp: false,
		Ap: &apConfig{
			Ip:         "192.168.27.1/24",
			Interface:  "uap0",
			Ssid:       "candy",
			Passphrase: "reckless",
			DhcpRange:  "192.168.27.100,192.168.27.150,1h",
		},
	}

	preCfg := defaultCfg

	if _, err := flags.Parse(&preCfg); err != nil {
		return nil, err
	}

	cfg := preCfg

	// Listen on the default interface/port if no listeners were specified.
	// An empty address string means default interface/address, which on
	// most unix systems is the same as 0.0.0.0.
	if len(cfg.RawListeners) == 0 {
		addr := fmt.Sprintf(":%d", defaultRPCPort)
		cfg.RawListeners = append(cfg.RawListeners, addr)
	}

	cfg.Listeners = make([]net.Addr, 0, len(cfg.RawListeners))
	for _, addr := range cfg.RawListeners {
		parsedAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			return nil, err
		}

		cfg.Listeners = append(cfg.Listeners, parsedAddr)
	}

	return &cfg, nil
}
