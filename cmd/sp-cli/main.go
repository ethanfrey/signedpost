package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethanfrey/signedpost/txn"
	crypto "github.com/tendermint/go-crypto"
)

type txPost struct {
	TX string `json:"tx"`
}

func main() {
	keyFile := flag.String("key", "/tmp/demo-key", "File location for private key to sign with (generated if missing)")
	acctName := flag.String("name", "Demo", "Name with which to create account (tx to sign)")
	flag.Parse()

	// first, make sure we have a private key
	var key crypto.PrivKey
	var err error
	if file, err := os.Open(*keyFile); err == nil {
		defer file.Close()
		// TODO: load from file
		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Printf("File error: %+v\n", err)
			return
		}
		key, err = crypto.PrivKeyFromBytes(data)
	} else {
		key = crypto.GenPrivKeyEd25519()
		outf, err := os.Create(*keyFile)
		if err != nil {
			fmt.Printf("File error: %+v\n", err)
			return
		}
		_, err = outf.Write(key.Bytes())
		outf.Close()
	}
	if err != nil {
		fmt.Printf("Key error: %+v\n", err)
		return
	}

	// create and sign transaction
	tx := txn.CreateAccountAction{Name: *acctName}
	data, err := txn.Send(tx, key)
	if err != nil {
		fmt.Printf("TX error: %+v\n", err)
		return
	}

	// encode for output
	res := txPost{
		TX: hex.EncodeToString(data),
	}
	out, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("JSON error: %+v\n", err)
		return
	}
	fmt.Println("*****")
	fmt.Println(string(out))
	fmt.Println("*****")
}
