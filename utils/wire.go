package utils

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/tendermint/go-wire"
)

// ToBinary serialize the object to a byte slice using go-wire
// obj: the object to serialize
func ToBinary(obj interface{}) ([]byte, error) {
	// return wire.BinaryBytes(obj), nil
	var err error
	w, n := new(bytes.Buffer), new(int)
	wire.WriteBinary(obj, w, n, &err)
	if err != nil {
		return nil, errors.Wrap(err, "To Binary")
	}
	return w.Bytes(), nil
}

// FromBinary deserialize the object from a byte slice using go-wire
// ptr: a pointer to the object to be filled (filled with data)
func FromBinary(data []byte, ptr interface{}) error {
	return errors.Wrap(wire.ReadBinaryBytes(data, ptr), "From Binary")
}
