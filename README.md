# `sweetd`

> 🔌 Daemon for pairing and control of the candy dispenser

* [x] Control of touch sensor, motor and buzzer
* [ ] Start pairing mode
* [ ] Bluetooth pairing
* [ ] WiFi pairing
* [ ] Setup and configuration of `lnd` node
* [ ] Bundle lnd in future
* [ ] Update procedure
* [ ] Automatically ensure that there's `denyinterfaces uap0` in `/etc/dhcpcd.conf`

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

## Generate grpc files

```text
protoc -I sweetrpc/ sweetrpc/rpc.proto --go_out=plugins=grpc:sweetrpc
```

## init.d

```
sudo ln -s $GOPATH/src/github.com/the-lightning-land/sweetd/contrib/init.d/sweetd /etc/init.d/sweetd
sudo update-rc.d sweetd defaults
sudo service start sweetd
```