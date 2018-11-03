# 🔌 `sweetd`

> Daemon for pairing and control of the Bitcoin-enabled candy dispenser

[![license](https://img.shields.io/github/license/the-lightning-land/sweetd.svg)](https://github.com/the-lightning-land/sweetd/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/the-lightning-land/sweetd.svg)](https://github.com/the-lightning-land/sweetd/releases)

## Intro

`sweetd` is the daemon process running on the Bitcoin-enabled candy dispenser.
It manages pairing and control, which is used by the Candy Dispenser app:

* 📱 [Candy Dispenser iOS app](https://github.com/the-lightning-land/Dispenser-iOS)
* 📱 *Candy Dispenser Android app coming in the future*

The `sweetd` program offers the following features:

* [x] 🍬 Control of the motor for dispensing candy
* [x] 📳 Control of the buzzer for user feedback
* [x] ☝️ React on events from the touch sensor
* [x] 📶 Pair through Wi-Fi hotspot
* [ ] 🔵 Pair through Bluetooth
* [x] 🌐 Set up Wi-Fi on candy dispenser
* [x] ⚡ Dispense candy on payments from remote `lnd` node
* [x] 💅 Customize name of your dispenser
* [x] 🔄 Update itself through app
* [ ] ⚙️ Ensure all system configs are made

## Usage

Download the pre-built binary for your system from the GitHub releases page.

* ⬇️ [Download `sweetd`](https://github.com/the-lightning-land/sweetd/releases)

Make sure that [all necessary dependencies](#dependencies) are installed.

Extract and open the downloaded archive, then run `sweetd`.
The following options are available. 

```
Usage:
  sweetd [OPTIONS]

Application Options:
  -v, --version                  Display version information and exit.
      --debug                    Start in debug mode.
      --listen=                  Add an interface/port/socket to listen for RPC
                                 connections
      --machine=[raspberry|mock] The machine controller to use. (default:
                                 raspberry)
      --ap                       Run the access point service.
      --datadir=                 The directory to store sweetd's data within.'
                                 (default: ./data)

Raspberry:
      --raspberry.touchpin=      The touch input pin. (default: 4)
      --raspberry.motorpin=      The motor output pin. (default: 27)
      --raspberry.buzzerpin=     The buzzer output pin. (default: 17)

Mock:
      --mock.listen=             Add an interface/port to listen for mock
                                 touches.

Access point:
      --ap.ip=                   IP address of device on access point
                                 interface. (default: 192.168.27.1/24)
      --ap.interface=            Name of the access point interface. (default:
                                 uap0)
      --ap.ssid=                 Name of the access point. (default: candy)
      --ap.passphrase=           WPA Passphrase (expected 8..63). (default:
                                 reckless)
      --ap.dhcprange=            IP range of the DHCP service. (default:
                                 192.168.27.100,192.168.27.150,1h)

Help Options:
  -h, --help                     Show this help message
```

## Dependencies

If you want to enable Wi-Fi pairing through the `--ap` option,
you'll need to install the following packages on your system: 

```
hostapd wireless-tools wpasupplicant dnsmasq iw
```

## Installation

`curl -LO https://github.com/the-lightning-land/sweetd/releases/download/v0.1.0/sweetd_0.1.0_linux_armv6.tar.gz`

`tar xfvz sweetd_0.1.0_linux_armv6.tar.gz`

`rm sweetd_0.1.0_linux_armv6.tar.gz`

`sudo mv sweetd_0.1.0_linux_armv6/sweetd /usr/local/bin/`

`sudo chown root:staff /usr/local/bin/sweetd`

`sudo echo "denyinterfaces uap0" >> /etc/dhcpcd.conf`

`sudo curl -L -o /etc/init.d/sweetd https://raw.githubusercontent.com/the-lightning-land/sweetd/master/contrib/init.d/sweetd`

`sudo chmod a+x /etc/init.d/sweetd`

`sudo update-rc.d sweetd defaults`

`sudo apt-get install hostapd wireless-tools wpasupplicant dnsmasq iw`

`sudo systemctl mask hostapd`

`sudo systemctl mask dnsmasq`

`sudo service sweetd start`

## Development

`go get -d github.com/the-lightning-land/sweetd`

`cd $GOPATH/src/github.com/the-lightning-land/sweetd`

`go build`

`./sweetd`

### Update lnd grpc client

`protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis lnrpc/rpc.proto --go_out=plugins=grpc:.`

## Releasing using [`goreleaser`](https://goreleaser.com)

`git tag -a v0.1.0 -m "Release name"`

`git push origin v0.1.0`

`goreleaser --rm-dist`

## Regenerate grpc files

`sweetd` exposes a gRPC API. When the API definition in `sweetrpc/rpc.proto` changes,
the following command will regenerate the source files:

```text
protoc -I sweetrpc/ sweetrpc/rpc.proto --go_out=plugins=grpc:sweetrpc
```

## init.d

```
sudo ln -s $GOPATH/src/github.com/the-lightning-land/sweetd/contrib/init.d/sweetd /etc/init.d/sweetd
sudo update-rc.d sweetd defaults
sudo service sweetd start
```