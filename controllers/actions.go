package controllers

import (
	"github.com/ethanfrey/bloggermint/db"
	"github.com/ethanfrey/bloggermint/models"
	tmsp "github.com/tendermint/tmsp/types"
)

func CreateAccount(ctx *Context, tx models.CreateAccountAction) tmsp.Result {
	if ctx.IsAnon() {
		return tmsp.NewError(tmsp.CodeType_Unauthorized, "Must sign transaction")
	}

	// TODO: make sure none with this name or pk already....

	account := db.Account{
		PK:         ctx.Signer(),
		Name:       tx.Name,
		Public:     tx.Public,
		EntryCount: 0,
	}

	// TODO: save account

	// return the new pk as response
	return tmsp.NewResultOK(account.PK, "")
}
