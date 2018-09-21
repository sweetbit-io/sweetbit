package main

import (
	"golang.org/x/net/context"
	"github.com/the-lightning-land/sweetd/sweetrpc"
)

type rpcServer struct {
	tmp <-chan bool
}

// A compile time check to ensure that rpcServer fully implements the SweetServer gRPC service.
var _ sweetrpc.SweetServer = (*rpcServer)(nil)

func newRPCServer(tmp <-chan bool) *rpcServer {
	return &rpcServer{
		tmp: tmp,
	}
}

func (s *rpcServer) GetFeature(ctx context.Context, req *sweetrpc.Test) (*sweetrpc.Test, error) {
	return &sweetrpc.Test{
		One: 233,
	}, nil
}

func (s *rpcServer) SubscribeWpaNetworks(req *sweetrpc.SubscribeWpaNetworksRequest,
	updateStream sweetrpc.Sweet_SubscribeWpaNetworksServer) error {

	for {
		<-s.tmp

		network := sweetrpc.WpaNetwork{
			Ssid: "a network",
		}

		if err := updateStream.Send(&network); err != nil {
			return err
		}
	}

	return nil
}
