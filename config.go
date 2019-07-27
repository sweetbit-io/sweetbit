package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"net"
)

const (
	defaultRPCPort = 9000
)

type raspberryConfig struct {
	TouchPin  string `long:"touchpin" description:"BCM number of the touch input pin."`
	MotorPin  string `long:"motorpin" description:"BCM number of the motor output pin."`
	BuzzerPin string `long:"buzzerpin" description:"BCM number of the buzzer output pin."`
}

type menderConfig struct {
	ConfigFile string `long:"config" description:"The file that holds mender configurations."`
	DataDir    string `long:"data" description:"The directory that stores mender data."`
}

type mockConfig struct {
	Listen string `long:"listen" description:"Add an interface/port to listen for mock touches."`
}

type lndConfig struct {
}

type torConfig struct {
	DataDir string `long:"data" description:"The directory that stores Tor data."`
}

type profilingConfig struct {
	Listen string `long:"listen" description:"Add an interface/port to listen for profiling data."`
}

type config struct {
	ShowVersion  bool     `short:"v" long:"version" description:"Display version information and exit."`
	Debug        bool     `long:"debug" description:"Start in debug mode."`
	RawListeners []string `long:"listen" description:"Add an interface/port/socket to listen for RPC connections"`
	Listeners    []net.Addr
	Machine      string           `long:"machine" description:"The machine controller to use." choice:"raspberry" choice:"mock"`
	Raspberry    *raspberryConfig `group:"Raspberry" namespace:"raspberry"`
	Mock         *mockConfig      `group:"Mock" namespace:"mock"`
	Lnd          *lndConfig       `group:"lnd" namespace:"lnd"`
	Net          string           `long:"net" description:"The networking system to use." choice:"dispenser" choice:"mock"`
	DataDir      string           `long:"datadir" description:"The directory to store sweetd's data within.'"`
	MemoPrefix   string           `long:"memoprefix" description:"Only react to invoices that have a memo starting with this prefix. (default empty, react to all invoices)'"`
	Updater      string           `long:"updater" description:"The updater to use." choice:"none" choice:"mender"`
	Mender       *menderConfig    `group:"Mender" namespace:"mender"`
	Tor          *torConfig       `group:"Tor" namespace:"tor"`
	Profiling    *profilingConfig `group:"Profiling" namespace:"profiling"`
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
		Lnd:        &lndConfig{},
		Net:        "dispenser",
		DataDir:    "./data",
		MemoPrefix: "",
		Updater:    "none",
		Mender: &menderConfig{
			ConfigFile: "/etc/mender/mender.conf",
			DataDir:    "/var/lib/mender",
		},
		Tor: &torConfig{
			DataDir: "./tor",
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
