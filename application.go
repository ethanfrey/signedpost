package signedpost

import (
	merkle "github.com/tendermint/go-merkle"
	tmsp "github.com/tendermint/tmsp/types"
)

// Application is the TMSP application for modifying state
type Application struct {
	tree *merkle.IAVLTree
}

// Info is a placeholder
func (app *Application) Info() string {
	return "Welcome to signed post"
}

// SetOption is ignored for now
func (app *Application) SetOption(key, value string) string {
	return "ignored"
}

// AppendTx actually does something
func (app *Application) AppendTx(tx []byte) tmsp.Result {
	// TODO
	return tmsp.Result{}
}

// CheckTx validates a tx for the mempool
func (app *Application) CheckTx(tx []byte) tmsp.Result {
	// TODO
	return tmsp.Result{}
}

// Query checks the current state
func (app *Application) Query(query []byte) tmsp.Result {
	// TODO
	return tmsp.Result{}
}

// Commit returns the application Merkle root hash
func (app *Application) Commit() tmsp.Result {
	// TODO
	return tmsp.Result{}
}
