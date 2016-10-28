package txs

import (
	"github.com/tendermint/go-crypto"
	merkle "github.com/tendermint/go-merkle"
)

// Context contains info on the current request
type Context struct {
	// Signer is the public key that signed the transaction (may be nil)
	signer crypto.PubKey
	store  *merkle.IAVLTree
}

// IsAnon is set if there is no signature
func (c *Context) IsAnon() bool {
	return c == nil || c.signer == nil
}

// Signer gets the signer's public key for authentication
func (c *Context) Signer() crypto.PubKey {
	return c.signer
}

func (c *Context) GetDB() *merkle.IAVLTree {
	return c.store
}
