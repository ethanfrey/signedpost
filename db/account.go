package db

import (
	"encoding/json"
	"fmt"

	merkle "github.com/tendermint/go-merkle"
)

// Account is a named account that can publish blog entries
type Account struct {
	PK         []byte `json:"-"`           // this is the public key of the owner
	Name       string `json:"name"`        // this is a name to search for
	EntryCount int    `json:"num_entries"` // total number of entries (de-normalize for speed)
	Public     bool   `json:"public"`      // if set to false, only the owner can read blog
}

func (acct *Account) ID() []byte {
	return acct.PK
}

func (acct *Account) Serialize() []byte {
	data, err := json.Marshal(acct)
	if err != nil {
		panic(err)
	}
	return data
}

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
func FindAccountByPK(store *merkle.IAVLTree, pk []byte) *Account {
	_, data, exists := store.Get(pk)
	if !exists || data == nil {
		return nil
	}
	acct, err := LoadAccount(pk, data)
	if err == nil {
		return acct
	}
	fmt.Println("FindAccountByPK:", err.Error())
	return nil
}

// FindAccountByName does a table-scan for name match (later secondary index?)
func FindAccountByName(store *merkle.IAVLTree, name string) *Account {
	var match *Account
	store.Iterate(func(key []byte, value []byte) bool {
		acct, err := LoadAccount(key, value)
		if err != nil && acct.Name == name {
			match = acct
			return true
		}
		return false
	})
	return match
}

// Create or update account
func (acct *Account) Save(store *merkle.IAVLTree) bool {
	data := acct.Serialize()
	return store.Set(acct.ID(), data)
}
