# `sweetd`

> 🔌 Daemon for pairing and control of the candy dispenser

* [x] Control of touch sensor, motor and buzzer
* [ ] Start pairing mode
* [ ] Bluetooth pairing
* [ ] WiFi pairing
* [ ] Setup and configuration of `lnd` node
* [ ] Bundle lnd in future
* [ ] Update procedure

## Dependencies

```
bridge hostapd wireless-tools wpa_supplicant dnsmasq iw
```

## Usage

`go get -d github.com/the-lightning-land/sweetd`

`go build`

`./sweetd`

## Releasing using [`goreleaser`](https://goreleaser.com)

`git tag -a v0.1.0 -m "Release name"`

`git push origin v0.1.0`

`goreleaser --rm-dist`
