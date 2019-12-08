# 🔌 `sweetd`

[![license](https://img.shields.io/github/license/the-lightning-land/sweetd.svg)](https://github.com/the-lightning-land/sweetd/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/the-lightning-land/sweetd.svg)](https://github.com/the-lightning-land/sweetd/releases)

> Sweet daemon for pairing and control of the Bitcoin-enabled candy dispenser

## Intro

`sweetd` is the daemon process running on the Bitcoin-enabled candy dispenser.
It manages pairing and control, which is used by the Candy Dispenser app:

* 📱 [Candy Dispenser iOS app](https://github.com/the-lightning-land/Dispenser-iOS)
* 📱 *Candy Dispenser Android app coming in the future*

The `sweetd` program offers the following features:

* [x] 🍬 Control of the motor for dispensing candy
* [x] 📳 Control of the buzzer for user feedback
* [x] ☝️ React on events from the touch sensor
* [x] 🔵 Pair through Bluetooth
* [x] 🌐 Set up Wi-Fi on candy dispenser
* [x] ⚡ Dispense candy on payments from remote `lnd` node
* [x] 💅 Customize name of your dispenser
* [x] 🔄 Update itself through app
* [x] ⚙️ Ensure all system configs are made

## Download

Download the pre-built binary for your system from the GitHub releases page.

* ⬇️ [Download `sweetd`](https://github.com/the-lightning-land/sweetd/releases)

Extract and open the downloaded archive, then run `sweetd`.

## Structure

The `sweetd` program's source code is split into small modules:

* 🔌 [`api`](api) - REST api for remote management of the dispenser
* ⚙️ [`app`](app) - website for managing the dispenser
* 🍬 [`dispenser`](dispenser) - orchestrator for everything the dispenser does
* ⚡️ [`lightning`](lightning) - controller for configured Lightning nodes, remote and local
* 🔩️ [`machine`](machine) - hardware controller for the touch sensor, motor and buzzer
* 📶 [`network`](network) - network subsystem that handles Wi-Fi discovery and connectivity
* 🤹‍ [`nodeman`](nodeman) - node manager
* 🧅 [`onion`](onion) - Tor onion service conveniences and .onion address generation
* 📲 [`pairing`](pairing) - pairing controller for BLE pairing
* 💵 [`pos`](pos) - point-of-sale website that creates invoices
* 🛑 [`reboot`](reboot) - methods for rebooting and shutting down the system 
* 📁 [`sweetdb`](sweetdb) - persistent database manager
* 📃 [`sweetlog`](sweetlog) - logging middleware for intercepting logs
* 🔖 [`sysid`](sysid) - methods for determining a system-specific id
* 🔄 [`updater`](updater) - update subsystem that controls system updates

## Configure data directory

By default, `sweetd` stores all data to `./data`.
You can easily override this location:

```sh
sweetd --datadir=/data/sweetd
``` 

## Configure machine access

Currently, the `sweetd` program is only tested and executed on a Raspberry Pi.
Running the executable with no options is the same as providing the following
options:

```sh
sweetd \
  --machine=raspberry \
  --raspberry.touchpin=25 \
  --raspberry.motorpin=23 \
  --raspberry.buzzerpin=24
```

You can also mock the underlying machine with the following option:

```sh
sweetd \
  --machine=mock \
  --mock.listen=localhost:5000
```

With this option, you can fake touches by sending simple
HTTP requests to the mock machine:

```
curl http://localhost:5000/touch/on
curl http://localhost:5000/touch/off
```

## Configure the `sweetd` API server

`sweetd` exposes a gRPC API. It can be used to configure the
Wi-Fi network that the candy dispenser connects to,
personalize it and change settings.

By default, the API server listens on `0.0.0.0:9000`. This can be changed
with the following option:

```sh
sweetd --listen=localhost:9000
```

It's also possible to specify multiple `--listen` options and
listen to multiple interfaces at once.

## Enable Wi-Fi hotspot pairing

At the moment, the only app pairing mechanism is through a Wi-Fi hotspot
that is created by the `sweetd` program.

This feature needs to be activated first:

```sh
sweetd --ap
```

Make sure that the following dependencies are installed when
running the access point mode:

```
hostapd wireless-tools wpasupplicant dnsmasq iw
```

The access point is configured with the below defaults. Any of these
can be changed to your needs.

```sh
sweetd \
  --ap \
  --ap.ip=192.168.27.1/24 \
  --ap.interface=uap0 \
  --ap.ssid=candy \
  --ap.passphrase=reckless \
  --ap.dhcprange=192.168.27.100,192.168.27.150,1h
```

This will create a Wi-Fi network called `candy` with the passphrase `reckless`.
An app will connect to that network for pairing and use
the gRPC api that is provided by the `sweetd` program.

## Development

`go get -d github.com/the-lightning-land/sweetd`

`cd $GOPATH/src/github.com/the-lightning-land/sweetd`

`go build`

`./sweetd`

## Releasing using [`goreleaser`](https://goreleaser.com)

The tool goreleaser can automatically sign the release and upload it to GitHub.

`git tag -a v0.1.0 -m "Release name"`

`git push origin v0.1.0`

`goreleaser --rm-dist`
