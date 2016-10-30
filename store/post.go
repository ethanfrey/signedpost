package store

import (
	"encoding/binary"

	"github.com/ethanfrey/signedpost/utils"
	"github.com/pkg/errors"

	merkle "github.com/tendermint/go-merkle"
)

// PostField is the post along with the key for lookup
type PostField struct {
	Key []byte
	Post
}

// Post represents one verifiably immutable blog entry (so no typos ;)
type Post struct {
	Number         int64
	PublishedBlock uint64
	Title          string
	Content        string
}

var postPrefix = []byte("p")
var endPostPrefix = []byte("q")

// AccountKeyFromPK creates the db key from a the account it belongs to and it's order
func PostKeyFromAccount(acct []byte, num int64) ([]byte, error) {
	if len(acct) < 16 {
		return nil, errors.New("Invalid account key")
	}
	output := append(postPrefix, acct...)
	suffix := make([]byte, 4)
	binary.BigEndian.PutUint32(suffix, uint32(num))
	output = append(output, suffix...)
	return output, nil
}

// Serialize turns the structure into bytes for storage and signing
func (p Post) Serialize() ([]byte, error) {
	return utils.ToBinary(p)
}

// Deserialize recovers the data bytes
func (p *Post) Deserialize(data []byte) error {
	return utils.FromBinary(data, p)
}

// Save stores they data at the given address
func (p Post) Save(store merkle.Tree, key []byte) (bool, error) {
	data, err := p.Serialize()
	if err != nil {
		return false, err
	}
	if len(key) < 16 {
		return false, errors.New("Key must be at least 16 bytes")
	}
	return store.Set(key, data), nil
}

// FindPostByAcctNum looks up by primary key (index scan)
// Error on storage error, if no match, returns nil
func FindPostByAcctNum(store merkle.Tree, acct []byte, num int64) (*PostField, error) {
	key, err := PostKeyFromAccount(acct, num)
	if err != nil {
		return nil, err
	}
	return FindPostByKey(store, key)
}

// FindPostByKey looks up the post by the db key
func FindPostByKey(store merkle.Tree, key []byte) (*PostField, error) {
	_, data, exists := store.Get(key)
	if !exists || data == nil {
		return nil, nil
	}
	p := Post{}
	err := p.Deserialize(data)
	return &PostField{Key: key, Post: p}, err
}

// FindPostsForAccount does a partial-index scan for all posts on a given account
func FindPostsForAccount(store merkle.Tree, acct []byte) ([]*PostField, error) {
	res := []*PostField{}
	start, _ := PostKeyFromAccount(acct, 0)
	end, err := PostKeyFromAccount(acct, 65000)
	if err != nil {
		return nil, err
	}
	store.IterateRange(start, end, true, func(key []byte, value []byte) bool {
		p := Post{}
		err = p.Deserialize(value)
		if err != nil {
			return true
		}
		res = append(res, &PostField{Key: key, Post: p})
		return false
	})
	return res, err
}
