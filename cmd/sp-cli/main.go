package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/ethanfrey/signedpost/txn"
	"github.com/ethanfrey/tenderize/sign"
)

type txPost struct {
	TX string `json:"tx"`
}

var (
	app     = kingpin.New("sp-cli", "A simple command line client for the signed post tendermint app")
	server  = app.Flag("server", "URL of signed post server").Default("http://localhost:54321").String()
	keyFile = app.Flag("key", "File location for private key to sign with (generated if missing)").Required().String()

	user = app.Command("account", "Create an account")
	name = user.Arg("name", "The username for the account").Required().String()

	post    = app.Command("post", "Add a new post")
	title   = post.Arg("title", "The title of the post").Required().String()
	content = post.Arg("content", "The post content").Required().String()
)

// ParseKey reads a keyfile or creates one if needed
func ParseKey(keyfile string) (crypto.PrivKey, error) {
	// first, make sure we have a private key
	var key crypto.PrivKey
	var err error
	if file, err := os.Open(*keyFile); err == nil {
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, errors.Wrap(err, "Loading key")
		}
		key, err = crypto.PrivKeyFromBytes(data)
		if err != nil {
			return nil, errors.Wrap(err, "Parsing key file")
		}
	} else {
		key = crypto.GenPrivKeyEd25519()
		outf, err := os.Create(*keyFile)
		if err != nil {
			return nil, errors.Wrap(err, "Creating key file")
		}
		_, err = outf.Write(key.Bytes())
		if err != nil {
			return nil, errors.Wrap(err, "Writing key")
		}
		outf.Close()
	}
	return key, err
}

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	// make sure we have a key to sign
	key, err := ParseKey(*keyFile)
	if err != nil {
		kingpin.Fatalf("Key error: %+v\n", err)
	}

	// generate and sign transaction based on input
	var data []byte
	switch cmd {
	case user.FullCommand():
		tx := txn.CreateAccountAction{Name: *name}
		data, err = sign.Send(tx, key)
	case post.FullCommand():
		tx := txn.AddPostAction{Title: *title, Content: *content}
		data, err = sign.Send(tx, key)
	}
	if err != nil {
		kingpin.Fatalf("Creating transaction: %+v\n", err)
	}

	// encode tx as json for output
	res := txPost{
		TX: hex.EncodeToString(data),
	}
	out, err := json.Marshal(res)
	if err != nil {
		kingpin.Fatalf("JSON error: %v\n", err)
	}

	// actually post to the server
	endpoint := *server + "/tndr/tx"
	client := http.Client{}
	resp, err := client.Post(endpoint, "application/json", bytes.NewBuffer(out))
	if err != nil {
		kingpin.Fatalf("HTTP Error: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		kingpin.Fatalf("Bad HTTP response: %d\n", resp.StatusCode)
	}

	// now we can print what happened!
	var msg ctypes.ResultBroadcastTx
	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		kingpin.Fatalf("Parse error tx response: %v\n", err)
	}
	fmt.Printf("Log: %s\n", msg.Log)
	fmt.Printf("ID: %s\n", hex.EncodeToString(msg.Data))
}
