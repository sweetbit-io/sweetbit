package main

import (
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
	"github.com/the-lightning-land/sweetd/sweetrpc"
	"github.com/the-lightning-land/sweetd/sysid"
	"github.com/the-lightning-land/sweetd/wpa"
	"golang.org/x/net/context"
	"time"
)

type rpcServerConfig struct {
	version   string
	commit    string
	dispenser *dispenser
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

func (s *rpcServer) GetInfo(ctx context.Context,
	req *sweetrpc.GetInfoRequest) (*sweetrpc.GetInfoResponse, error) {

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

func (s *rpcServer) GetWpaConnectionInfo(ctx context.Context,
	req *sweetrpc.GetWpaConnectionInfoRequest) (*sweetrpc.GetWpaConnectionInfoResponse, error) {

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

func (s *rpcServer) ConnectWpaNetwork(ctx context.Context,
	req *sweetrpc.ConnectWpaNetworkRequest) (*sweetrpc.ConnectWpaNetworkResponse, error) {

	log.Info("Adding new network...")

	net, err := wpa.AddNetwork("wlan0")
	if err != nil {
		log.Errorf("Adding network failed: %s", err.Error())
		return nil, errors.New("Connection failed")
	}

	log.Infof("Setting ssid %v for network %v...", req.Ssid, net)

	err = wpa.SetNetwork("wlan0", net, wpa.Ssid, req.Ssid)
	if err != nil {
		log.Errorf("Setting ssid failed: %s", err.Error())
		return nil, errors.New("Connection failed")
	}

	log.Infof("Setting psk for network %v...", net)

	err = wpa.SetNetwork("wlan0", net, wpa.Psk, req.Psk)
	if err != nil {
		log.Errorf("Setting psk failed: %s", err.Error())
		return nil, errors.New("Connection failed")
	}

	log.Infof("Enabling network %v...", net)

	err = wpa.EnableNetwork("wlan0", net)
	if err != nil {
		log.Errorf("Enabling network failed: %s", err.Error())
		return nil, errors.New("Connection failed")
	}

	tries := 1

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		log.Info("Fetching connection status...")

		status, err := wpa.GetStatus("wlan0")
		if err != nil {
			log.Errorf("Getting WPA status failed: %s", err.Error())
			return nil, errors.New("Getting WPA status failed")
		}

		log.Infof("Got status %v for ssid %v.", status.State, status.Ssid)

		if status.Ssid == req.Ssid && status.State == "COMPLETED" {
			log.Infof("Saving network %v", net)

			err := wpa.Save("wlan0")
			if err != nil {
				log.Errorf("Saving config failed: %s", err.Error())
			}

			log.Info("Confirming successful connection")

			return &sweetrpc.ConnectWpaNetworkResponse{
				Status: sweetrpc.ConnectWpaNetworkResponse_CONNECTED,
			}, nil
		}

		if tries > 5 {
			break
		}

		tries++
	}

	log.Errorf("Could not connect to network %v", req.Ssid)

	return &sweetrpc.ConnectWpaNetworkResponse{
		Status: sweetrpc.ConnectWpaNetworkResponse_FAILED,
	}, nil
}

func (s *rpcServer) GetWpaNetworks(ctx context.Context,
	req *sweetrpc.GetWpaNetworksRequest) (*sweetrpc.GetWpaNetworksResponse, error) {

	err := wpa.Scan("wlan0")
	if err != nil {
		log.Errorf("Scan failed: %s", err.Error())
		return nil, errors.New("Scan failed")
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	<-ticker.C

	results, err := wpa.Results("wlan0")
	if err != nil {
		log.Errorf("Scan failed: %s", err.Error())
		return nil, errors.New("Scan failed")
	}

	log.Infof("Found %v networks", len(results))

	// Map results []*wpa.Network to networks []*sweetrpc.WpaNetwork
	networks := make([]*sweetrpc.WpaNetwork, len(results))
	for i, result := range results {
		networks[i] = &sweetrpc.WpaNetwork{
			Ssid:        result.Ssid,
			Bssid:       result.Bssid,
			Flags:       result.Flags,
			Frequency:   result.Frequency,
			SignalLevel: result.SignalLevel,
		}
	}

	return &sweetrpc.GetWpaNetworksResponse{
		Networks: networks,
	}, nil
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

func (s *rpcServer) ConnectToRemoteNode(ctx context.Context,
	req *sweetrpc.ConnectToRemoteNodeRequest) (*sweetrpc.ConnectToRemoteNodeResponse, error) {
	log.Infof("Connecting to lightning node %s", req.Uri)

	err := s.config.dispenser.connectLndNode(req.Uri, req.Cert, req.Macaroon)
	if err != nil {
		log.Errorf("Connection failed: %s", err.Error())
		return nil, errors.New("Connection failed")
	}

	err = s.config.dispenser.saveLndNode(req.Uri, req.Cert, req.Macaroon)
	if err != nil {
		log.Errorf("Could not save remote lightning connection: %s", err)
	}

	return &sweetrpc.ConnectToRemoteNodeResponse{}, nil
}
