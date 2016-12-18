package store

import (
	"bytes"

	"github.com/ethanfrey/tenderize/mom"
)

var (
	accountIDLength = 20
	minAccountID    = bytes.Repeat([]byte{0}, accountIDLength)
	maxAccountID    = bytes.Repeat([]byte{255}, accountIDLength)
)

func init() {
	// IMPORTANT: you must call this in the init, so all serialization works
	mom.RegisterModels(Account{}, Post{})
}
