package client

import (
	"encoding/json"

	"github.com/tendermint/go-rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	// "github.com/tendermint/tendermint/rpc/test/client_test.go"
)

type LightClient struct {
	remote   string
	endpoint string
	rpc      *rpcclient.ClientJSONRPC
	ws       *rpcclient.WSClient
}

func New(remote, wsEndpoint string) *LightClient {
	return &LightClient{
		rpc:      rpcclient.NewClientJSONRPC(remote),
		remote:   remote,
		endpoint: wsEndpoint,
	}
}

func (c *LightClient) Status() (*ctypes.ResultStatus, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("status", []interface{}{}, tmResult)
	if err != nil {
		return nil, err
	}
	// note: panics if rpc doesn't match.  okay???
	return (*tmResult).(*ctypes.ResultStatus), nil
}

func (c *LightClient) TMSPInfo() (*ctypes.ResultTMSPInfo, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("tmsp_info", []interface{}{}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultTMSPInfo), nil
}

func (c *LightClient) TMSPQuery(query []byte) (*ctypes.ResultTMSPQuery, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("tmsp_query", []interface{}{query}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultTMSPQuery), nil
}

func (c *LightClient) BroadcastTxCommit(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return c.broadcastTX("broadcast_tx_commit", tx)
}

func (c *LightClient) BroadcastTxAsync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return c.broadcastTX("broadcast_tx_async", tx)
}

func (c *LightClient) BroadcastTxSync(tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	return c.broadcastTX("broadcast_tx_sync", tx)
}

func (c *LightClient) broadcastTX(route string, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call(route, []interface{}{tx}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultBroadcastTx), nil
}

func (c *LightClient) NetInfo() (*ctypes.ResultNetInfo, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("net_info", nil, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultNetInfo), nil
}

func (c *LightClient) DialSeeds(seeds []string) (*ctypes.ResultDialSeeds, error) {
	tmResult := new(ctypes.TMResult)
	// TODO: is this the correct way to marshall seeds?
	_, err := c.rpc.Call("dial_seeds", []interface{}{seeds}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultDialSeeds), nil
}

func (c *LightClient) BlockchainInfo(minHeight, maxHeight int) (*ctypes.ResultBlockchainInfo, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("blockchain", []interface{}{minHeight, maxHeight}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultBlockchainInfo), nil
}

func (c *LightClient) Genesis() (*ctypes.ResultGenesis, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("genesis", nil, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultGenesis), nil
}

func (c *LightClient) Block(height int) (*ctypes.ResultBlock, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("block", []interface{}{height}, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultBlock), nil
}

func (c *LightClient) Validators() (*ctypes.ResultValidators, error) {
	tmResult := new(ctypes.TMResult)
	_, err := c.rpc.Call("validators", nil, tmResult)
	if err != nil {
		return nil, err
	}
	return (*tmResult).(*ctypes.ResultValidators), nil
}

/** websocket event stuff here... **/

// StartWebsocket starts up a websocket and a listener goroutine
// if already started, do nothing
func (c *LightClient) StartWebsocket() error {
	var err error
	if c.ws == nil {
		ws := rpcclient.NewWSClient(c.remote, c.endpoint)
		_, err = ws.Start()
		if err == nil {
			c.ws = ws
		}
	}
	return err
}

// StopWebsocket stops the websocket connection
func (c *LightClient) StopWebsocket() {
	if c.ws != nil {
		c.ws.Stop()
		c.ws = nil
	}
}

// GetEventChannels returns the results and error channel from the websocket
func (c *LightClient) GetEventChannels() (chan json.RawMessage, chan error) {
	if c.ws == nil {
		return nil, nil
	}
	return c.ws.ResultsCh, c.ws.ErrorsCh
}

func (c *LightClient) Subscribe(event string) error {
	return c.ws.Subscribe(event)
}

func (c *LightClient) Unsubscribe(event string) error {
	return c.ws.Unsubscribe(event)
}
