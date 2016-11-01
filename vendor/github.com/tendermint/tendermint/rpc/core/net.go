package core

import (
	"fmt"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

//-----------------------------------------------------------------------------

func NetInfo() (*ctypes.ResultNetInfo, error) {
	listening := p2pSwitch.IsListening()
	listeners := []string{}
	for _, listener := range p2pSwitch.Listeners() {
		listeners = append(listeners, listener.String())
	}
	peers := []ctypes.Peer{}
	for _, peer := range p2pSwitch.Peers().List() {
		peers = append(peers, ctypes.Peer{
			NodeInfo:         *peer.NodeInfo,
			IsOutbound:       peer.IsOutbound(),
			ConnectionStatus: peer.Connection().Status(),
		})
	}
	return &ctypes.ResultNetInfo{
		Listening: listening,
		Listeners: listeners,
		Peers:     peers,
	}, nil
}

//-----------------------------------------------------------------------------

// Dial given list of seeds if we have no outbound peers
func DialSeeds(seeds []string) (*ctypes.ResultDialSeeds, error) {
	outbound, _, _ := p2pSwitch.NumPeers()
	if outbound != 0 {
		return nil, fmt.Errorf("Already have some outbound peers")
	}
	// starts go routines to dial each seed after random delays
	p2pSwitch.DialSeeds(seeds)
	return &ctypes.ResultDialSeeds{}, nil
}

//-----------------------------------------------------------------------------

func Genesis() (*ctypes.ResultGenesis, error) {
	return &ctypes.ResultGenesis{genDoc}, nil
}
