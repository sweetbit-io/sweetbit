package main

import (
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
	"github.com/the-lightning-land/sweetd/reboot"
	"github.com/the-lightning-land/sweetd/sweetrpc"
	"github.com/the-lightning-land/sweetd/sysid"
	"github.com/the-lightning-land/sweetd/wpa"
	"golang.org/x/net/context"
	"time"
)

type rpcServerConfig struct {
	version string
	commit  string
}

type rpcServer struct {
	dispenser *dispenser
	config    *rpcServerConfig
}

// A compile time check to ensure that rpcServer fully implements the SweetServer gRPC service.
var _ sweetrpc.SweetServer = (*rpcServer)(nil)

func newRPCServer(dispenser *dispenser, config *rpcServerConfig) *rpcServer {
	return &rpcServer{
		dispenser: dispenser,
		config:    config,
	}
}

func (s *rpcServer) GetInfo(ctx context.Context,
	req *sweetrpc.GetInfoRequest) (*sweetrpc.GetInfoResponse, error) {
	log.Info("Retrieving info...")

	id, err := sysid.GetId()
	if err != nil {
		return nil, err
	}

	var remoteNode *sweetrpc.RemoteNode = nil

	if s.dispenser.lightningNodeUri != "" {
		remoteNode = &sweetrpc.RemoteNode{
			Uri: s.dispenser.lightningNodeUri,
		}
	}

	name, err := s.dispenser.getName()
	if err != nil {
		log.Errorf("Failed getting info: %v", err)
		return nil, errors.New("Failed getting info")
	}

	if name == "" {
		name = "Candy Dispenser"
		// name = fmt.Sprintf("Candy %v", id)
	}

	return &sweetrpc.GetInfoResponse{
		Serial:          id,
		Version:         s.config.version,
		Commit:          s.config.commit,
		RemoteNode:      remoteNode,
		Name:            name,
		DispenseOnTouch: s.dispenser.dispenseOnTouch,
		BuzzOnDispense:  s.dispenser.buzzOnDispense,
	}, nil
}

func (s *rpcServer) SetName(ctx context.Context, req *sweetrpc.SetNameRequest) (*sweetrpc.SetNameResponse, error) {
	log.Infof("Setting name to '%v'...", req.Name)

	err := s.dispenser.setName(req.Name)
	if err != nil {
		log.Errorf("Failed setting name: %v", err)
		return nil, errors.New("Failed setting name")
	}

	return &sweetrpc.SetNameResponse{}, nil
}

func (s *rpcServer) SetDispenseOnTouch(ctx context.Context, req *sweetrpc.SetDispenseOnTouchRequest) (*sweetrpc.SetDispenseOnTouchResponse, error) {
	log.Infof("Setting dispense on touch to '%v'...", req.DispenseOnTouch)

	err := s.dispenser.setDispenseOnTouch(req.DispenseOnTouch)
	if err != nil {
		log.Errorf("Failed setting dispense on touch: %v", err)
		return nil, errors.New("Failed setting dispense on touch")
	}

	return &sweetrpc.SetDispenseOnTouchResponse{}, nil
}

func (s *rpcServer) SetBuzzOnDispense(ctx context.Context, req *sweetrpc.SetBuzzOnDispenseRequest) (*sweetrpc.SetBuzzOnDispenseResponse, error) {
	log.Infof("Setting buzz on dispense to '%v'...", req.BuzzOnDispense)

	err := s.dispenser.setBuzzOnDispense(req.BuzzOnDispense)
	if err != nil {
		log.Errorf("Failed setting buzz on dispense: %v", err)
		return nil, errors.New("Failed setting buzz on dispense")
	}

	return &sweetrpc.SetBuzzOnDispenseResponse{}, nil
}

func (s *rpcServer) GetWpaConnectionInfo(ctx context.Context,
	req *sweetrpc.GetWpaConnectionInfoRequest) (*sweetrpc.GetWpaConnectionInfoResponse, error) {
	log.Info("Getting wpa connection info...")

	status, err := wpa.GetStatus("wlan0")
	if err != nil {
		log.Errorf("Getting WPA status failed: %v", err)
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
		log.Errorf("Enabling network failed: %v", err)
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
				log.Errorf("Saving config failed: %v", err)
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
	log.Info("Getting wpa networks...")

	err := wpa.Scan("wlan0")
	if err != nil {
		log.Errorf("Scan failed: %v", err)
		return nil, errors.New("Scan failed")
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	<-ticker.C

	results, err := wpa.Results("wlan0")
	if err != nil {
		log.Errorf("Scan failed: %v", err)
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
		log.Errorf("Update failed: %v", err)
		return nil, errors.New("Update failed")
	}

	return &sweetrpc.UpdateResponse{}, nil
}

func (s *rpcServer) ConnectToRemoteNode(ctx context.Context,
	req *sweetrpc.ConnectToRemoteNodeRequest) (*sweetrpc.ConnectToRemoteNodeResponse, error) {
	log.Infof("Connecting to lightning node %s", req.Uri)

	err := s.dispenser.connectLndNode(req.Uri, req.Cert, req.Macaroon)
	if err != nil {
		log.Errorf("Connection failed: %v", err)
		return nil, errors.New("Connection failed")
	}

	err = s.dispenser.saveLndNode(req.Uri, req.Cert, req.Macaroon)
	if err != nil {
		log.Errorf("Could not save remote lightning connection: %v", err)
	}

	return &sweetrpc.ConnectToRemoteNodeResponse{}, nil
}

func (s *rpcServer) DisconnectFromRemoteNode(ctx context.Context,
	req *sweetrpc.DisconnectFromRemoteNodeRequest) (*sweetrpc.DisconnectFromRemoteNodeResponse, error) {
	log.Info("Disconnecting from lightning node")

	err := s.dispenser.disconnectLndNode()
	if err != nil {
		log.Errorf("Disconnect failed: %v", err)
		return nil, errors.New("Disconnect failed")
	}

	err = s.dispenser.deleteLndNode()
	if err != nil {
		log.Errorf("Could not delete remote lightning connection: %v", err)
	}

	return &sweetrpc.DisconnectFromRemoteNodeResponse{}, nil
}

func (s *rpcServer) Reboot(ctx context.Context,
	req *sweetrpc.RebootRequest) (*sweetrpc.RebootResponse, error) {
	log.Info("Rebooting dispenser...")

	err := reboot.Reboot()
	if err != nil {
		log.Errorf("Reboot failed: %v", err)
		return nil, errors.New("Reboot failed")
	}

	return &sweetrpc.RebootResponse{}, nil
}

func (s *rpcServer) ToggleDispenser(ctx context.Context,
	req *sweetrpc.ToggleDispenserRequest) (*sweetrpc.ToggleDispenserResponse, error) {
	log.Infof("Toggling dispenser %v", req.Dispense)

	s.dispenser.toggleDispense(req.Dispense)

	return &sweetrpc.ToggleDispenserResponse{}, nil
}

func (r *rpcServer) SubscribeDispenses(req *sweetrpc.SubscribeDispensesRequest,
	updateStream sweetrpc.Sweet_SubscribeDispensesServer) error {
	log.Info("Subscribing to dispenses")

	return nil
}