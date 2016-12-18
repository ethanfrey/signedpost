package mom

import (
	wutil "github.com/ethanfrey/tenderize/wire"
	"github.com/tendermint/go-merkle"
	"github.com/tendermint/go-wire"
)

// Model is an abstraction over an object to be stored in MerkleDB
type Model interface {
	// Key returns the db key for this model.
	// This key may have zero values (and designed for range queries), depending on the state of the model
	// The Key should not change with any allowed transformation of the model
	Key() Key
}

// mwire is designed to bridge Models for go-wire
type mwire struct {
	Model
}

// Save attempts to save the given model in the given store
// updated is false on insert, otherwise true
// error is non-nil if key or value cannot be serialized
func Save(store merkle.Tree, model Model) (updated bool, err error) {
	key, err := KeyToBytes(model.Key())
	if err != nil {
		return false, err
	}

	data, err := ModelToBytes(model)
	if err != nil {
		return false, err
	}

	return store.Set(key, data), nil
}

// ModelToBytes converts the model into bytes to store in the db
// If there are invalid values in the model you can return an error
func ModelToBytes(model Model) ([]byte, error) {
	return wutil.ToBinary(mwire{model})
}

// ModelFromBytes sets the model contents to the passed in data
// Returns error if the data doesn't match this model
func ModelFromBytes(data []byte) (Model, error) {
	// Is there an easier way
	holder := mwire{}
	err := wutil.FromBinary(data, &holder)
	return holder.Model, err
}

// RegisterModels takes a list of Models and registers them with go-wire for Serialization
// The control byte is based on the order, so if you want to maintain compatibility with an
// existing data store, do not change the position of any items.  You can use nil as a
// placeholder to not use that byte anymore.
//
// This also registers the serailzers for the model keys with the same prefix
//
// eg. RegisterModels(Account{}, nil, Status{}) gives account Byte 1, Status Byte 3
func RegisterModels(models ...Model) {
	// prepare for max size, but might be shorter with nil
	regMods := make([]wire.ConcreteType, 0, len(models))
	regKeys := make([]wire.ConcreteType, 0, len(models))

	for i, model := range models {
		// add the model and key for all non-nil values
		if model != nil {
			b := byte(i + 1)
			regMods = append(regMods, wire.ConcreteType{O: model, Byte: b})
			regKeys = append(regKeys, wire.ConcreteType{O: model.Key(), Byte: b})
		}
	}

	// now register this with go-wire
	wire.RegisterInterface(mwire{}, regMods...)
	wire.RegisterInterface(mkey{}, regKeys...)
}
