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

	// make sure none with this name or pk already....
	exists := db.FindAccountByPK(ctx.GetDB(), ctx.Signer())
	if exists != nil {
		return tmsp.NewError(tmsp.CodeType_BaseDuplicateAddress,
			"Account exists for this public key")
	}
	exists = db.FindAccountByName(ctx.GetDB(), tx.Name)
	if exists != nil {
		return tmsp.NewError(tmsp.CodeType_BaseDuplicateAddress,
			"Account name already taken")
	}

	// all safe, go save it
	account := db.Account{
		PK:         ctx.Signer(),
		Name:       tx.Name,
		Public:     tx.Public,
		EntryCount: 0,
	}
	account.Save(ctx.GetDB())
	// return the new pk as response
	return tmsp.NewResultOK(account.PK, "")
}

// TogglePublic changes the oublic flag on an existing account
func TogglePublic(ctx *Context, tx models.TogglePublicAction) tmsp.Result {
	if ctx.IsAnon() {
		return tmsp.NewError(tmsp.CodeType_Unauthorized, "Must sign transaction")
	}

	// find the account for this user....
	acct := db.FindAccountByPK(ctx.GetDB(), ctx.Signer())
	if acct == nil {
		return tmsp.NewError(tmsp.CodeType_BaseUnknownAddress,
			"No account exists for this public key")
	}

	// update the value and save it if needed
	if acct.Public != tx.Public {
		acct.Public = tx.Public
		acct.Save(ctx.GetDB())
	}

	return tmsp.NewResultOK(nil, "")
}
