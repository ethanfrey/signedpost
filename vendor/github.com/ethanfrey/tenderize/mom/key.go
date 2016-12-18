package mom

import (
	wutil "github.com/ethanfrey/tenderize/wire"
	"github.com/tendermint/go-merkle"
)

// Key is designed to let you easily build compound indexes using the one database key
type Key interface {
	// Range assumes the key has one or more elements at the zero (nil) value
	// It should return two keys, the first is with those zero elements set the the minimum possible value
	// The second return value should be with the zero elements set to their maximum possible value
	Range() (min Key, max Key)
}

// mkey is designed to bridge Keys for go-wire
type mkey struct {
	Key
}

// Query can define a detailed range query based on a key
type Query struct {
	Key     Key
	Reverse bool
	Filter  func(Model) bool
}

// KeyToBytes returns the Key as a byte array for use in go-merkle
// Returns an error if some required fields are not set
func KeyToBytes(key Key) ([]byte, error) {
	return wutil.ToBinary(mkey{key})
}

// Load attempts to find the data matching the given key
// If the key or store data cannot be parsed, returns error
// If there is no data, Model is nil
func Load(store merkle.Tree, key Key) (Model, error) {
	k, err := KeyToBytes(key)
	if err != nil {
		return nil, err
	}

	_, data, exists := store.Get(k)
	if !exists || data == nil {
		return nil, nil
	}

	return ModelFromBytes(data)
}

// ByteRange attempts to take the range and serialize it, returns error if either fails
func ByteRange(key Key) (start []byte, end []byte, err error) {
	s, e := key.Range()
	start, err = KeyToBytes(s)
	if err != nil {
		return
	}
	end, err = KeyToBytes(e)
	return
}

// List returns all items that match this query
func List(store merkle.Tree, q Query) ([]Model, error) {
	return filter(store, q.Key, !q.Reverse, q.Filter)
}

// filter returns a list of all items that fit in the range of the key, optionally applying a filter
// if filter is nil, then return all items in that range
func filter(store merkle.Tree, key Key, ascending bool, filter func(Model) bool) ([]Model, error) {
	res := []Model{}
	start, end, err := ByteRange(key)
	if err != nil {
		return nil, err
	}

	store.IterateRange(start, end, ascending, func(k []byte, v []byte) bool {
		item, err := ModelFromBytes(v)
		if err == nil && (filter == nil || filter(item)) {
			res = append(res, item)
		}
		return false
	})
	return res, nil
}
