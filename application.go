package signedpost

import (
	"github.com/ethanfrey/signedpost/redux"
	"github.com/ethanfrey/signedpost/txn"
	merkle "github.com/tendermint/go-merkle"
	tmsp "github.com/tendermint/tmsp/types"
)

// Application is the TMSP application for modifying state
type Application struct {
	commited *redux.Service
	check    *redux.Service
}

// NewApp creates a new tmsp application
func NewApp(tree merkle.Tree) *Application {
	a := Application{
		commited: redux.New(tree, 0),
	}
	a.check = a.commited.Copy()
	return &a
}

// Info is a placeholder
func (app *Application) Info() string {
	return app.commited.Info()
}

// SetOption is ignored for now
func (app *Application) SetOption(key, value string) string {
	return "ignored"
}

// AppendTx actually does something
func (app *Application) AppendTx(tx []byte) tmsp.Result {
	action, err := txn.Receive(tx)
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}
	return app.commited.Apply(action)
}

// CheckTx validates a tx for the mempool
func (app *Application) CheckTx(tx []byte) tmsp.Result {
	action, err := txn.Receive(tx)
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}
	return app.check.Apply(action)
}

// Query returns contents behind given key
func (app *Application) Query(query []byte) tmsp.Result {
	_, val, exists := app.commited.GetDB().Get(query)
	if !exists {
		return tmsp.NewError(tmsp.CodeType_BaseUnknownAddress, "")
	}
	return tmsp.NewResultOK(val, "")
}

// Commit returns the application Merkle root hash
func (app *Application) Commit() tmsp.Result {
	app.check = app.commited.Copy()
	hash := app.commited.Hash()
	return tmsp.NewResultOK(hash, "")
}

// InitChain make it blockchain aware
// (but we ignore all but BeginBlock)
func (app *Application) InitChain(validators []*tmsp.Validator) {}

// BeginBlock signals the beginning of a block, update service so we tag posts properly
func (app *Application) BeginBlock(height uint64) {
	// TODO: this is never called in the current code, so we make do with EndBlock, implying a begin block
}

// EndBlock signals the end of a block, ignored now
// diffs: changed validators from app to TendermintCore
func (app *Application) EndBlock(height uint64) (diffs []*tmsp.Validator) {
	app.commited.SetHeight(height + 1)
	return nil
}
