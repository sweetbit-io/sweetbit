package onion

import (
	"context"
	"crypto"
	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil"
	"net"
	"sync"
	"time"
)

type ServiceState int

const (
	ServiceStateStarting ServiceState = iota
	ServiceStateStarted
	ServiceStateStopping
	ServiceStateStopped
)

type Service struct {
	tor      *tor.Tor
	log      Logger
	service  *tor.OnionService
	key      crypto.PrivateKey
	port     int
	listener net.Listener
	context  context.Context
	cancel   context.CancelFunc
}

type ServiceConfig struct {
	Tor    *tor.Tor
	Logger Logger
	Port   int
	Key    crypto.PrivateKey
}

func NewService(config *ServiceConfig) *Service {
	service := &Service{
		tor: config.Tor,
		key: config.Key,
	}

	if config.Logger != nil {
		service.log = config.Logger
	} else {
		service.log = noopLogger{}
	}

	if config.Port > 0 {
		service.port = config.Port
	} else {
		service.port = 80
	}

	return service
}

func (s *Service) ID() string {
	return torutil.OnionServiceIDFromPrivateKey(s.key)
}

var mutex = sync.Mutex{}

func (s *Service) Start() {
	go func() {
		mutex.Lock()

		s.context, s.cancel = context.WithTimeout(context.Background(), 120*time.Second)

		s.log.Infof("starting onion %s.onion", s.ID())

		var err error
		s.service, err = s.tor.Listen(s.context, &tor.ListenConf{
			LocalListener: s.listener,
			Key:           s.key,
			RemotePorts:   []int{s.port},
		})
		if err != nil {
			s.log.Errorf("unable to start %s.onion: %v", s.ID(), err)
		} else {
			s.log.Infof("started %s.onion", s.ID())
		}

		mutex.Unlock()
	}()
}

func (s *Service) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Service) SetListener(listener net.Listener) {
	s.listener = listener
}

func (s *Service) SetPrivateKey(key crypto.PrivateKey) {
	s.key = key
}
