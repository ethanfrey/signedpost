package store

import (
	"strings"

	"github.com/ethanfrey/tenderize/mom"
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-merkle"
)

// Account is a named account that can publish blog entries
// This can be serialized with go-wire
type Account struct {
	ID         []byte
	Name       string // this is a name to search for
	EntryCount int64  // total number of entries (de-normalize for speed)
}

// AccountKey wraps the immutible ID
type AccountKey struct {
	ID []byte
}

// Key returns just the address of the account
func (a Account) Key() mom.Key {
	return AccountKey{ID: a.ID}
}

// Range returns all account keys iff the is not set, just the specified account otherwise
func (k AccountKey) Range() (min, max mom.Key) {
	if len(k.ID) == accountIDLength {
		return k, k
	}
	// TODO: if len > 0 but < accountIDLength, then use the prefix and fill the rest with 0 or 255 for min, max
	return AccountKey{ID: minAccountID}, AccountKey{ID: maxAccountID}
}

// NewAccount generates the account id from a public key
func NewAccount(pk crypto.PubKey, name string) Account {
	addr := pk.Address()
	return Account{ID: addr, Name: name}
}

// AccountMatchesName returns a mom.Query filter for exact matches of account name
func AccountMatchesName(name string) func(mom.Model) bool {
	return func(m mom.Model) bool {
		acct, ok := m.(Account)
		return ok && acct.Name == name
	}
}

// AccountContainsName returns a mom.Query filter for accounts that include this name substring
func AccountContainsName(name string) func(mom.Model) bool {
	lname := strings.ToLower(name)
	return func(m mom.Model) bool {
		acct, ok := m.(Account)
		return ok && strings.Contains(strings.ToLower(acct.Name), lname)
	}
}

// FindAccount looks up by primary key (index scan)
// Error on storage error, if no match, returns nil
func FindAccount(store merkle.Tree, pk crypto.PubKey) (*Account, error) {
	key := NewAccount(pk, "").Key()
	model, err := mom.Load(store, key)
	if err != nil || model == nil {
		return nil, err
	}
	res := model.(Account)
	return &res, nil
}

// ListAccounts makes a search over all accounts, and casts them to the proper type
// note an empty response returns no error
func ListAccounts(store merkle.Tree, filter func(mom.Model) bool) ([]Account, error) {
	query := mom.Query{
		Key:    AccountKey{},
		Filter: filter,
	}
	models, err := mom.List(store, query)
	if err != nil {
		return nil, err
	}
	res := make([]Account, len(models))
	for i := range models {
		res[i] = models[i].(Account)
	}
	return res, nil
}
