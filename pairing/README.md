# Bluetooth Low Energy (BLE) pairing

```
ca001001-75dd-4a0e-b688-66b7df342cc6
\/\-/\-/\--------------------------/
⎟  ⎟  ⎟  ⎟
⎟  ⎟  ⎟   \- Candy Dispenser UUID Prefix
⎟  ⎟  ⎟   
⎟  ⎟   \- Descriptor identifier 
⎟  ⎟
⎟   \- Characteristic identifier
⎟
 \- Service identifier
```

## `ca000000` pairing service

### `ca001000` `status` characteristic
  
* read `offline`
* ~~write~~
* notify `running`

### `ca002000` `scanWifi` characteristic

* read `true`
* write `true`
* notify `false`

### `ca003000` `discoveredWifi` characteristic

* ~~read~~
* ~~write~~
* notify `{ "ssid": "onion", "secure": "true" }`

### `ca004000` `connectWifi` characteristic

* ~~read~~
* write `{ "ssid": "onion", "psk": "" }`
* ~~notify~~

### `ca005000` `onionApi` characteristic

* read `kgwozt2fbhdhruhi.onion`
* ~~write~~
* ~~notify~~
