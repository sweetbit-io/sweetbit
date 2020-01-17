package lightning

import (
	"net"
)

// A compile time check to ensure that noopListener fully implements the net.Listener interface
var _ net.Listener = (*noopListener)(nil)

type noopListener struct {
	addr net.Addr
}

func (n noopListener) Accept() (net.Conn, error) {
	return nil, nil
}

func (n noopListener) Close() error {
	return nil
}

func (n noopListener) Addr() net.Addr {
	return n.addr
}
