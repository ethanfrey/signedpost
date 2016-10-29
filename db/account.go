package db

import (
	"github.com/ethanfrey/signedpost/utils"
	"github.com/pkg/errors"

	crypto "github.com/tendermint/go-crypto"
	merkle "github.com/tendermint/go-merkle"
)

var accountPrefix = []byte("u")
var endAccountPrefix = []byte("v")

// Field represents the Key, Value pair for a leaf node
type Field struct {
	Key   []byte
	Value []byte
}

// Account is a named account that can publish blog entries
// This can be serialized with go-wire
type Account struct {
	Name       string // this is a name to search for
	EntryCount int    // total number of entries (de-normalize for speed)
}

// AccountKeyFromPK creates the db key from a public key
func AccountKeyFromPK(pk crypto.PubKey) ([]byte, error) {
	if pk == nil {
		return nil, errors.New("Empty private key")
	}
	addr := pk.Address()
	return append(accountPrefix, addr...), nil
}

// Serialize turns the structure into bytes for storage and signing
func (acct Account) Serialize() ([]byte, error) {
	return utils.ToBinary(acct)
}

// Deserialize recovers the data bytes
func (acct *Account) Deserialize(data []byte) error {
	return utils.FromBinary(data, acct)
}

// Save stores they data at the given address
func (acct Account) Save(store *merkle.IAVLTree, key []byte) (bool, error) {
	data, err := acct.Serialize()
	if err != nil {
		return false, err
	}
	if len(key) < 16 {
		return false, errors.New("Key must be at least 16 bytes")
	}
	return store.Set(key, data), nil
}

// FindAccountByPK looks up by primary key (index scan)
// Error on storage error, if no match, returns nil
func FindAccountByPK(store *merkle.IAVLTree, pk crypto.PubKey) (*Account, error) {
	key, err := AccountKeyFromPK(pk)
	if err != nil {
		return nil, err
	}
	return FindAccountByKey(store, key)
}

// FindAccountByKey looks up the account by the db key
func FindAccountByKey(store *merkle.IAVLTree, key []byte) (*Account, error) {
	_, data, exists := store.Get(key)
	if !exists || data == nil {
		return nil, nil
	}
	acct := Account{}
	err := acct.Deserialize(data)
	return &acct, err
}

// FindAccountByName does a table-scan over accounts for name match (later secondary index?)
func FindAccountByName(store *merkle.IAVLTree, name string) (*Account, error) {
	var match *Account
	acct := new(Account)
	store.IterateRange(accountPrefix, endAccountPrefix, true, func(key []byte, value []byte) bool {
		err := acct.Deserialize(value)
		if err == nil && acct.Name == name {
			match = acct
			return true
		}
		return false
	})
	return match, nil
}
