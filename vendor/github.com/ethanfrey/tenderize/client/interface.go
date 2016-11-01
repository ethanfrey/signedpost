package client

import (
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

// Client is a generic interface used to access an tendermint node
// you can use this interface to allow easy mocking for tests without a full node
type Client interface {
	Status() (*ctypes.ResultStatus, error)
	NetInfo() (*ctypes.ResultNetInfo, error)
	DialSeeds(seeds []string) (*ctypes.ResultDialSeeds, error)
	BlockchainInfo(minHeight, maxHeight int) (*ctypes.ResultBlockchainInfo, error)
	Genesis() (*ctypes.ResultGenesis, error)
	Block(height int) (*ctypes.ResultBlock, error)
	Validators() (*ctypes.ResultValidators, error)

	BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTx, error)
	BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTx, error)
	BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error)

	TMSPQuery(query []byte) (*ctypes.ResultTMSPQuery, error)
	TMSPInfo() (*ctypes.ResultTMSPInfo, error)

	// subscribe to events (how to read them depends on implementation...)
	Subscribe(event string) error
	Unsubscribe(event string) error

	// are these needed by clients??
	// DumpConsensusState() (*ctypes.ResultDumpConsensusState, error)
	// UnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error)
	// NumUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error)

	// Also, skip "unsafe" methods
}
