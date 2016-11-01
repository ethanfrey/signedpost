package rpcclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-rpc/types"
	"github.com/tendermint/go-wire"
)

// TODO: Deprecate support for IP:PORT or /path/to/socket
func makeHTTPDialer(remoteAddr string) (string, func(string, string) (net.Conn, error)) {

	parts := strings.SplitN(remoteAddr, "://", 2)
	var protocol, address string
	if len(parts) != 2 {
		log.Warn("WARNING (go-rpc): Please use fully formed listening addresses, including the tcp:// or unix:// prefix")
		protocol = rpctypes.SocketType(remoteAddr)
		address = remoteAddr
	} else {
		protocol, address = parts[0], parts[1]
	}

	trimmedAddress := strings.Replace(address, "/", ".", -1) // replace / with . for http requests (dummy domain)
	return trimmedAddress, func(proto, addr string) (net.Conn, error) {
		return net.Dial(protocol, address)
	}
}

// We overwrite the http.Client.Dial so we can do http over tcp or unix.
// remoteAddr should be fully featured (eg. with tcp:// or unix://)
func makeHTTPClient(remoteAddr string) (string, *http.Client) {
	address, dialer := makeHTTPDialer(remoteAddr)
	return "http://" + address, &http.Client{
		Transport: &http.Transport{
			Dial: dialer,
		},
	}
}

//------------------------------------------------------------------------------------

type Client interface {
}

//------------------------------------------------------------------------------------

// JSON rpc takes params as a slice
type ClientJSONRPC struct {
	address string
	client  *http.Client
}

func NewClientJSONRPC(remote string) *ClientJSONRPC {
	address, client := makeHTTPClient(remote)
	return &ClientJSONRPC{
		address: address,
		client:  client,
	}
}

func (c *ClientJSONRPC) Call(method string, params []interface{}, result interface{}) (interface{}, error) {
	return c.call(method, params, result)
}

func (c *ClientJSONRPC) call(method string, params []interface{}, result interface{}) (interface{}, error) {
	// Make request and get responseBytes
	request := rpctypes.RPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      "",
	}
	requestBytes := wire.JSONBytes(request)
	requestBuf := bytes.NewBuffer(requestBytes)
	// log.Info(Fmt("RPC request to %v (%v): %v", c.remote, method, string(requestBytes)))
	httpResponse, err := c.client.Post(c.address, "text/json", requestBuf)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	responseBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}
	// 	log.Info(Fmt("RPC response: %v", string(responseBytes)))
	return unmarshalResponseBytes(responseBytes, result)
}

//-------------------------------------------------------------

// URI takes params as a map
type ClientURI struct {
	address string
	client  *http.Client
}

func NewClientURI(remote string) *ClientURI {
	address, client := makeHTTPClient(remote)
	return &ClientURI{
		address: address,
		client:  client,
	}
}

func (c *ClientURI) Call(method string, params map[string]interface{}, result interface{}) (interface{}, error) {
	return c.call(method, params, result)
}

func (c *ClientURI) call(method string, params map[string]interface{}, result interface{}) (interface{}, error) {
	values, err := argsToURLValues(params)
	if err != nil {
		return nil, err
	}
	//log.Info(Fmt("URI request to %v (%v): %v", c.address, method, values))
	resp, err := c.client.PostForm(c.address+"/"+method, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return unmarshalResponseBytes(responseBytes, result)
}

//------------------------------------------------

func unmarshalResponseBytes(responseBytes []byte, result interface{}) (interface{}, error) {
	// read response
	// if rpc/core/types is imported, the result will unmarshal
	// into the correct type
	// log.Notice("response", "response", string(responseBytes))
	var err error
	response := &rpctypes.RPCResponse{}
	err = json.Unmarshal(responseBytes, response)
	if err != nil {
		return nil, errors.New(Fmt("Error unmarshalling rpc response: %v", err))
	}
	errorStr := response.Error
	if errorStr != "" {
		return nil, errors.New(Fmt("Response error: %v", errorStr))
	}
	// unmarshal the RawMessage into the result
	result = wire.ReadJSONPtr(result, *response.Result, &err)
	if err != nil {
		return nil, errors.New(Fmt("Error unmarshalling rpc response result: %v", err))
	}
	return result, nil
}

func argsToURLValues(args map[string]interface{}) (url.Values, error) {
	values := make(url.Values)
	if len(args) == 0 {
		return values, nil
	}
	err := argsToJson(args)
	if err != nil {
		return nil, err
	}
	for key, val := range args {
		values.Set(key, val.(string))
	}
	return values, nil
}

func argsToJson(args map[string]interface{}) error {
	var n int
	var err error
	for k, v := range args {
		buf := new(bytes.Buffer)
		wire.WriteJSON(v, buf, &n, &err)
		if err != nil {
			return err
		}
		args[k] = buf.String()
	}
	return nil
}
