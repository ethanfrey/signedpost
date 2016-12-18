package redux

import (
	"fmt"

	"github.com/ethanfrey/signedpost/txn"
	"github.com/ethanfrey/tenderize/sign"
	merkle "github.com/tendermint/go-merkle"
	tmsp "github.com/tendermint/tmsp/types"
)

// Service contains all static info to process transactions
type Service struct {
	// TODO: logger, block height
	store       merkle.Tree
	blockHeight uint64
}

func New(tree merkle.Tree, height uint64) *Service {
	return &Service{
		store:       tree,
		blockHeight: height,
	}
}

func (s *Service) GetDB() merkle.Tree {
	return s.store
}

func (s *Service) GetHeight() uint64 {
	return s.blockHeight
}

func (s *Service) SetHeight(h uint64) {
	s.blockHeight = h
}

func (s *Service) Info() string {
	return fmt.Sprintf("size:%v", s.store.Size())
}

func (s *Service) Hash() []byte {
	if s.store.Size() == 0 {
		return nil
	}
	return s.store.Hash()
}

func (s *Service) Copy() *Service {
	return &Service{
		store:       s.store.Copy(),
		blockHeight: s.blockHeight,
	}
}

// Apply will take any authentication action and apply it to the store
// TODO: change result type??
func (s *Service) Apply(tx sign.ValidatedAction) tmsp.Result {
	switch action := tx.GetAction().(type) {
	case txn.CreateAccountAction:
		return s.CreateAccount(action, tx.GetSigner())
	case txn.AddPostAction:
		return s.AppendPost(action, tx.GetSigner())
	}
	return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, "Unknown action")
}
