package nodeman

import "github.com/the-lightning-land/sweetd/lightning"

type NodeConfig interface{}

type RemoteLndNodeConfig struct {
	Name     string
	Uri      string
	Cert     []byte
	Macaroon []byte
}

type LocalNodeConfig struct {
	Name string
}

type LightningNode interface {
	lightning.Node
	ID() string
	Name() string
	Enabled() bool
	setEnabled(enabled bool)
	setName(name string)
}

type RemoteLndNode struct {
	*lightning.LndNode
	id      string
	name    string
	enabled bool
	Uri     string
}

func (n *RemoteLndNode) ID() string              { return n.id }
func (n *RemoteLndNode) Name() string            { return n.name }
func (n *RemoteLndNode) setName(name string)     { n.name = name }
func (n *RemoteLndNode) Enabled() bool           { return n.enabled }
func (n *RemoteLndNode) setEnabled(enabled bool) { n.enabled = enabled }

type LocalNode struct {
	*lightning.LocalNode
	id      string
	name    string
	enabled bool
}

func (n *LocalNode) ID() string              { return n.id }
func (n *LocalNode) Name() string            { return n.name }
func (n *LocalNode) setName(name string)     { n.name = name }
func (n *LocalNode) Enabled() bool           { return n.enabled }
func (n *LocalNode) setEnabled(enabled bool) { n.enabled = enabled }
