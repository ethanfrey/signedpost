package tests

/**
This file is base HEAVILY on tendermint/tendermint/rpc/tests/helpers.go
However, I wanted to use public variables, so this could be a basis
of tests in various packages.
**/

import (
	"github.com/ethanfrey/tenderize/client"

	cfg "github.com/tendermint/go-config"
	p2p "github.com/tendermint/go-p2p"
	"github.com/tendermint/tendermint/config/tendermint_test"
	nm "github.com/tendermint/tendermint/node"
)

var (
	config cfg.Config
)

// GetConfig returns a config for the test cases as a singleton
func GetConfig() cfg.Config {
	if config == nil {
		config = tendermint_test.ResetConfig("rpc_test_client_test")
	}
	return config
}

// GetClient gets a rpc client pointing to the test node
func GetClient() *client.LightClient {
	rpcAddr := GetConfig().GetString("rpc_laddr")
	return client.New(rpcAddr, "/websocket")
}

// StartNode starts a test node in a go routine and returns when it is initialized
// TODO: can one pass an Application in????
func StartNode() {
	// start a node
	ready := make(chan struct{})
	go NewNode(ready)
	<-ready
}

// NewNode creates a new node and sleeps forever
func NewNode(ready chan struct{}) {
	// Create & start node
	node := nm.NewNodeDefault(GetConfig())
	protocol, address := nm.ProtocolAndAddress(config.GetString("node_laddr"))
	l := p2p.NewDefaultListener(protocol, address, true)
	node.AddListener(l)
	node.Start()

	// Run the RPC server.
	node.StartRPC()
	ready <- struct{}{}

	// Sleep forever
	ch := make(chan struct{})
	<-ch
}
