package redux

import (
	"github.com/tendermint/go-crypto"
	merkle "github.com/tendermint/go-merkle"
)

// Service contains info on the current request
type Service struct {
	// Signer is the public key that signed the transaction (may be nil)
	signer crypto.PubKey
	store  *merkle.IAVLTree
}

// IsAnon is set if there is no signature
func (c *Service) IsAnon() bool {
	return c == nil || c.signer == nil
}

// Signer gets the signer's public key for authentication
func (c *Service) Signer() crypto.PubKey {
	return c.signer
}

func (c *Service) GetDB() *merkle.IAVLTree {
	return c.store
}
