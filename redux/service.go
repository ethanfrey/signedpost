package redux

import (
	"github.com/ethanfrey/signedpost/txn"
	merkle "github.com/tendermint/go-merkle"
	tmsp "github.com/tendermint/tmsp/types"
)

// Service contains all static info to process transactions
type Service struct {
	// TODO: logger, block height
	store       *merkle.IAVLTree
	blockHeight int64
}

func (c *Service) GetDB() *merkle.IAVLTree {
	return c.store
}

func (c *Service) GetHeight() int64 {
	return c.blockHeight
}

func (c *Service) SetHeight(h int64) {
	c.blockHeight = h
}

// Apply will take any authentication action and apply it to the store
// TODO: change result type??
func (c *Service) Apply(tx txn.ValidatedAction) tmsp.Result {
	switch action := tx.GetAction().(type) {
	case txn.CreateAccountAction:
		return c.CreateAccount(action, tx.GetSigner())
	case txn.AddPostAction:
		return c.AppendPost(action, tx.GetSigner())
	}
	return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, "Unknown action")
}
