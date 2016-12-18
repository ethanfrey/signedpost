package signedpost

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/ethanfrey/signedpost/utils"
	"github.com/ethanfrey/tenderize/client"
	"github.com/ethanfrey/tenderize/sign"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const minTxLength = 20
const maxBlocks = 50

// Proxy validates queries and sends appropriate ones to the tendermint core
type Proxy struct {
	client client.Client
}

// NewProxy creates a Proxy pointing to the url of the rpc server on tendermint core
func NewProxy(baseURL string) Proxy {
	return Proxy{
		client: client.New(baseURL, "/websocket"),
	}
}

type txPost struct {
	TX string `json:"tx"`
}

// PostTransaction validates the posted transaction and submits it to the tendermint consensus engine
// Format: {"tx": "deadbeef"} - hex encoded transaction as tx key in a json blob
func (p Proxy) PostTransaction(rw http.ResponseWriter, r *http.Request) {
	var res *ctypes.ResultBroadcastTx
	var tx []byte
	post := txPost{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&post)
	if err == nil {
		tx, err = hex.DecodeString(post.TX)
		if len(tx) < minTxLength {
			err = errors.New("tx data truncated or missing")
		}
	}
	if err == nil {
		var val sign.ValidatedAction
		val, err = sign.Receive(tx)
		if err == nil && val.IsAnon() {
			err = errors.New("All transactions require a valid signature")
		}
	}
	// at this point, we have an error, or we known body is an acceptable transaction
	if err == nil {
		res, err = p.client.BroadcastTxSync(tx)
	}
	utils.RenderQuery(rw, res, err)
}

// GetStatus returns the status of the tendermint core
func (p Proxy) GetStatus(rw http.ResponseWriter, r *http.Request) {
	status, err := p.client.Status()
	utils.RenderQuery(rw, status, err)
}

// GetValidators returns the current validator set
func (p Proxy) GetValidators(rw http.ResponseWriter, r *http.Request) {
	vals, err := p.client.Validators()
	utils.RenderQuery(rw, vals, err)
}

// GetBlock validates the parameters and gets one block from tendermint core
func (p Proxy) GetBlock(rw http.ResponseWriter, r *http.Request) {
	var res *ctypes.ResultBlock
	h, err := strconv.Atoi(r.URL.Query().Get("height"))
	if err == nil {
		res, err = p.client.Block(h)
	}
	utils.RenderQuery(rw, res, err)
}

// GetChain returns a list of block headers between minHeight and maxHeight.
// Here, we require both and don't allow a query for more than 50 blocks at once
func (p Proxy) GetChain(rw http.ResponseWriter, r *http.Request) {
	var res *ctypes.ResultBlockchainInfo
	min, err := strconv.Atoi(r.URL.Query().Get("minHeight"))
	if err == nil {
		max, err := strconv.Atoi(r.URL.Query().Get("maxHeight"))
		if err == nil {
			if max-min < 0 {
				err = errors.New("maxHeight must be greater than minHeight")
			} else if max-min > maxBlocks {
				err = errors.Errorf("You cannot query more than %d blocks at once", maxBlocks)
			} else {
				res, err = p.client.BlockchainInfo(min, max)
			}
		}
	}
	utils.RenderQuery(rw, res, err)
}

// AddChainRoutes adds a routes for tendermint core interactions to the router
func (p Proxy) AddChainRoutes(r *mux.Router) {
	tndr := r.PathPrefix("/tndr").Subrouter()
	tndr.HandleFunc("/tx", p.PostTransaction).Methods("POST")
	tndr.HandleFunc("/status", p.GetStatus).Methods("GET")
	tndr.HandleFunc("/validators", p.GetValidators).Methods("GET")
	tndr.HandleFunc("/block", p.GetBlock).Methods("GET")
	tndr.HandleFunc("/blockchain", p.GetChain).Methods("GET")
}
