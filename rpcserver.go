package main

import (
	"golang.org/x/net/context"
	"github.com/the-lightning-land/sweetd/sweetrpc"
	"github.com/the-lightning-land/sweetd/sysid"
	log "github.com/sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/the-lightning-land/sweetd/wpa"
	"time"
)

type rpcServerConfig struct {
	version string
	commit  string
}

type rpcServer struct {
	config *rpcServerConfig
}

// A compile time check to ensure that rpcServer fully implements the SweetServer gRPC service.
var _ sweetrpc.SweetServer = (*rpcServer)(nil)

func newRPCServer(config *rpcServerConfig) *rpcServer {
	return &rpcServer{
		config: config,
	}
}

func (s *rpcServer) GetInfo(ctx context.Context, req *sweetrpc.GetInfoRequest) (*sweetrpc.GetInfoResponse, error) {
	id, err := sysid.GetId()
	if err != nil {
		return nil, err
	}

	return &sweetrpc.GetInfoResponse{
		Serial:  id,
		Version: s.config.version,
		Commit:  s.config.commit,
	}, nil
}

func (s *rpcServer) GetWpaConnectionInfo(ctx context.Context, req *sweetrpc.GetWpaConnectionInfoRequest) (*sweetrpc.GetWpaConnectionInfoResponse, error) {
	status, err := wpa.GetStatus("wlan0")
	if err != nil {
		log.Errorf("Getting WPA status failed: %s", err.Error())
		return nil, errors.New("Getting WPA status failed")
	}

	return &sweetrpc.GetWpaConnectionInfoResponse{
		Ssid:  status.Ssid,
		State: status.State,
		Ip:    status.Ip,
	}, nil
}

func (s *rpcServer) ConnectWpaNetwork(req *sweetrpc.ConnectWpaNetworkRequest,
	updateStream sweetrpc.Sweet_ConnectWpaNetworkServer) error {

	net, err := wpa.AddNetwork("wlan0")
	if err != nil {
		log.Errorf("Adding network failed: %s", err.Error())
		return errors.New("Connection failed")
	}

	err = wpa.SetNetwork("wlan0", net, wpa.Ssid, req.Ssid)
	if err != nil {
		log.Errorf("Setting ssid failed: %s", err.Error())
		return errors.New("Connection failed")
	}

	err = wpa.SetNetwork("wlan0", net, wpa.Psk, req.Psk)
	if err != nil {
		log.Errorf("Setting psk failed: %s", err.Error())
		return errors.New("Connection failed")
	}

	err = wpa.EnableNetwork("wlan0", net)
	if err != nil {
		log.Errorf("Enabling network failed: %s", err.Error())
		return errors.New("Connection failed")
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		status, err := wpa.GetStatus("wlan0")
		if err != nil {
			log.Errorf("Getting WPA status failed: %s", err.Error())
			return errors.New("Getting WPA status failed")
		}

		if status.Ssid == req.Ssid && status.State == "COMPLETED" {
			err := wpa.Save("wlan0")
			if err != nil {
				log.Errorf("Saving config failed: %s", err.Error())
			}

			err = updateStream.Send(&sweetrpc.WpaConnectionUpdate{
				Status: sweetrpc.WpaConnectionUpdate_CONNECTED,
			})

			if err != nil {
				log.Errorf("Sending update failed: %s", err.Error())
				return errors.New("Scan failed")
			}

			return nil
		}
	}

	return nil
}

func (s *rpcServer) SubscribeWpaNetworkScanUpdates(req *sweetrpc.SubscribeWpaNetworkScanUpdatesRequest,
	updateStream sweetrpc.Sweet_SubscribeWpaNetworkScanUpdatesServer) error {

	var previous []wpa.Network

	err := wpa.Scan("wlan0")
	if err != nil {
		log.Errorf("Scan failed: %s", err.Error())
		return errors.New("Scan failed")
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		results, err := wpa.Results("wlan0")
		if err != nil {
			log.Errorf("Scan failed: %s", err.Error())
			return errors.New("Scan failed")
		}

		log.Infof("Found %v networks", len(results))

		for _, result := range results {
			alreadyAnnounced := false

			for _, prev := range previous {
				if result.Ssid == prev.Ssid {
					alreadyAnnounced = true
				}
			}

			if !alreadyAnnounced {
				err = updateStream.Send(&sweetrpc.WpaNetworkScanUpdate{
					Update: &sweetrpc.WpaNetworkScanUpdate_Appeared{
						Appeared: &sweetrpc.WpaNetwork{
							Ssid: result.Ssid,
						},
					},
				})

				if err != nil {
					log.Errorf("Scan failed: %s", err.Error())
					return errors.New("Scan failed")
				}
			}
		}

		for _, prev := range previous {
			stillHere := false

			for _, result := range results {
				if result.Ssid == prev.Ssid {
					stillHere = true
				}
			}

			if !stillHere {
				err = updateStream.Send(&sweetrpc.WpaNetworkScanUpdate{
					Update: &sweetrpc.WpaNetworkScanUpdate_Gone{
						Gone: &sweetrpc.WpaNetwork{
							Ssid: prev.Ssid,
						},
					},
				})

				if err != nil {
					log.Errorf("Scan failed: %s", err.Error())
					return errors.New("Scan failed")
				}
			}
		}

		previous = results
	}

	return nil
}

func (s *rpcServer) Update(ctx context.Context, req *sweetrpc.UpdateRequest) (*sweetrpc.UpdateResponse, error) {
	log.Infof("Go update request with %s", req.Url)

	err := doUpdate(req.Url)
	if err != nil {
		log.Errorf("Update failed: %s", err.Error())
		return nil, errors.New("Update failed")
	}

	return &sweetrpc.UpdateResponse{}, nil
}
