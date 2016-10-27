package db

import (
	"encoding/json"

	"github.com/pkg/errors"

	merkle "github.com/tendermint/go-merkle"
)

var accountPrefix = []byte("u")

// Account is a named account that can publish blog entries
type Account struct {
	PK         []byte `json:"-"`           // this is the public key of the owner
	Name       string `json:"name"`        // this is a name to search for
	EntryCount int    `json:"num_entries"` // total number of entries (de-normalize for speed)
}

func accountPKToID(pk []byte) ([]byte, error) {
	if pk == nil || len(pk) < 1 {
		return nil, errors.New("Empty private key")
	}
	return append(accountPrefix, pk...), nil
}

func accountIDtoPK(id []byte) ([]byte, error) {
	if id == nil || len(id) < 2 || id[0] != accountPrefix[0] {
		return nil, errors.New("Invalid account ID")
	}
	return id[1 : len(id)-1], nil
}

// ID gives you the db id of the account
func (acct *Account) ID() []byte {
	id, _ := accountPKToID(acct.PK)
	return id
}

// Serialize turns the structure into bytes for storage and signing
func (acct *Account) Serialize() []byte {
	data, err := json.Marshal(acct)
	if err != nil {
		panic(err)
	}
	return data
}

// LoadAccount takes db data and puts it into the structure
func LoadAccount(key, value []byte) (*Account, error) {
	acct := new(Account)
	err := json.Unmarshal(value, acct)
	if err != nil {
		return nil, err
	}
	acct.PK = key
	return acct, nil
}

// FindAccountByPK looks up by primary key (index scan)
func FindAccountByPK(store *merkle.IAVLTree, pk []byte) (*Account, error) {
	id, err := accountPKToID(pk)
	if err != nil {
		return nil, err
	}
	_, data, exists := store.Get(id)
	if !exists || data == nil {
		return nil, errors.New("No account for this pk")
	}
	return LoadAccount(pk, data)
}

// FindAccountByName does a table-scan for name match (later secondary index?)
func FindAccountByName(store *merkle.IAVLTree, name string) (*Account, error) {
	var match *Account
	store.Iterate(func(key []byte, value []byte) bool {
		acct, err := LoadAccount(key, value)
		if err != nil && acct.Name == name {
			match = acct
			return true
		}
		return false
	})
	return match, nil
}

// Save will create or update account
func (acct *Account) Save(store *merkle.IAVLTree) bool {
	data := acct.Serialize()
	return store.Set(acct.ID(), data)
}
