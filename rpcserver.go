package main

import (
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
	"github.com/the-lightning-land/sweetd/reboot"
	"github.com/the-lightning-land/sweetd/sweetdb"
	"github.com/the-lightning-land/sweetd/sweetrpc"
	"github.com/the-lightning-land/sweetd/sysid"
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
		// TODO: Name the dispenser individually by default
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

	log.Info("Requested wifi connection info")

	status, err := s.dispenser.accessPoint.GetConnectionStatus()
	if err != nil {
		log.Errorf("Could not get Wifi connection status: %v", err)
		return nil, errors.New("Could not get Wifi connection status")
	}

	log.Info("Getting wpa connection info...")

	return &sweetrpc.GetWpaConnectionInfoResponse{
		Ssid:  status.Ssid,
		State: status.State,
		Ip:    status.Ip,
	}, nil
}

func (s *rpcServer) ConnectWpaNetwork(ctx context.Context,
	req *sweetrpc.ConnectWpaNetworkRequest) (*sweetrpc.ConnectWpaNetworkResponse, error) {

	log.Infof("Requested wifi connection to %v", req.Ssid)

	err := s.dispenser.accessPoint.ConnectWifi(req.Ssid, req.Psk)
	if err != nil {
		log.Errorf("Could not get Wifi networks: %v", err)
		return nil, errors.New("Could not get Wifi networks")
	}

	tries := 1

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		status, err := s.dispenser.accessPoint.GetConnectionStatus()
		if err != nil {
			log.Errorf("Getting WPA status failed: %s", err.Error())
			return nil, errors.New("Getting Wifi connection status failed")
		}

		log.Infof("Got status %v for ssid %v.", status.State, status.Ssid)

		if status.Ssid == req.Ssid && (status.State == "ASSOCIATED" || status.State == "COMPLETED") {
			err := s.dispenser.setWifiConnection(&sweetdb.WifiConnection{
				Ssid: req.Ssid,
				Psk:  req.Psk,
			})
			if err != nil {
				log.Errorf("Could not save wifi connection: %v", err)
			}

			return &sweetrpc.ConnectWpaNetworkResponse{
				Status: sweetrpc.ConnectWpaNetworkResponse_CONNECTED,
			}, nil
		}

		if tries > 30 {
			break
		}

		tries++
	}

	log.Errorf("Could not connect to network %v in time", req.Ssid)

	return &sweetrpc.ConnectWpaNetworkResponse{
		Status: sweetrpc.ConnectWpaNetworkResponse_FAILED,
	}, nil
}

func (s *rpcServer) GetWpaNetworks(ctx context.Context,
	req *sweetrpc.GetWpaNetworksRequest) (*sweetrpc.GetWpaNetworksResponse, error) {

	log.Info("Requested wifi networks")

	networks, err := s.dispenser.accessPoint.ListWifiNetworks()
	if err != nil {
		log.Errorf("Getting Wifi networks failed: %v", err)
		return nil, errors.New("Getting Wifi networks failed")
	}

	log.Infof("Found %v networks", len(networks))

	results := make([]*sweetrpc.WpaNetwork, len(networks))

	for i, result := range networks {
		results[i] = &sweetrpc.WpaNetwork{
			Ssid:        result.Ssid,
			Bssid:       result.Bssid,
			Flags:       result.Flags,
			Frequency:   result.Frequency,
			SignalLevel: result.SignalLevel,
		}
	}

	return &sweetrpc.GetWpaNetworksResponse{
		Networks: results,
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

	client := r.dispenser.subscribeDispenses()

	defer client.cancel()

	for {
		on := <-client.Dispenses

		log.Debugf("Sending dispense event %v to client %v", on, client.id)

		dispense := &sweetrpc.Dispense{
			Dispense: on,
		}

		if err := updateStream.Send(dispense); err != nil {
			log.Infof("Client %v failed receiving: %v", client.id, err)
			return err
		}
	}

	return nil
}
