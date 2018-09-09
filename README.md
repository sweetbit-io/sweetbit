# `sweetd`

> 🔌 Daemon for pairing and control of the candy dispenser

* [x] Control of touch sensor, motor and buzzer
* [ ] Bluetooth pairing
* [ ] WiFi pairing
* [ ] Setup and configuration of `lnd` node

## Usage

`go get -d github.com/the-lightning-land/sweetd`

`go build`

`./sweetd`

## Releasing

`git tag -a v0.1.0 -m "Release name"`

`git push origin v0.1.0`

`goreleaser --rm-dist`
