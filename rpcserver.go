package main

import (
	"golang.org/x/net/context"
	"github.com/the-lightning-land/sweetd/sweetrpc"
	"github.com/the-lightning-land/sweetd/sysid"
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
	return &sweetrpc.GetWpaConnectionInfoResponse{
		Ssid:    "",
		State:   "",
		Ip:      "",
		Message: "",
	}, nil
}

func (s *rpcServer) ConnectWpaNetwork(req *sweetrpc.ConnectWpaNetworkRequest,
	updateStream sweetrpc.Sweet_ConnectWpaNetworkServer) error {

	update := sweetrpc.WpaConnectionUpdate{
		Status: sweetrpc.WpaConnectionUpdate_CONNECTED,
	}

	if err := updateStream.Send(&update); err != nil {
		return err
	}

	return nil
}

func (s *rpcServer) SubscribeWpaNetworkScanUpdates(req *sweetrpc.SubscribeWpaNetworkScanUpdatesRequest,
	updateStream sweetrpc.Sweet_SubscribeWpaNetworkScanUpdatesServer) error {

	network := &sweetrpc.WpaNetwork{
		Ssid: "a network",
	}

	update := &sweetrpc.WpaNetworkScanUpdate{
		Update: &sweetrpc.WpaNetworkScanUpdate_Appeared{
			Appeared: network,
		},
	}

	if err := updateStream.Send(update); err != nil {
		return err
	}

	return nil
}
