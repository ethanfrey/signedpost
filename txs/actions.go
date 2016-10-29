package txs

import (
	"github.com/ethanfrey/signedpost/db"
	"github.com/ethanfrey/signedpost/models"
	tmsp "github.com/tendermint/tmsp/types"
)

// CreateAccount creates a new account based on the signing public key
func CreateAccount(ctx *Context, tx models.CreateAccountAction) tmsp.Result {
	if ctx.IsAnon() {
		return tmsp.NewError(tmsp.CodeType_Unauthorized, "Must sign transaction")
	}

	// make sure none with this name or pk already....
	exists, err := db.FindAccountByPK(ctx.GetDB(), ctx.Signer())
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}
	if exists != nil {
		return tmsp.NewError(tmsp.CodeType_BaseDuplicateAddress,
			"Account exists for this public key")
	}

	exists, err = db.FindAccountByName(ctx.GetDB(), tx.Name)
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}
	if exists != nil {
		return tmsp.NewError(tmsp.CodeType_BaseDuplicateAddress,
			"Account name already taken")
	}

	// all safe, go save it
	account := db.Account{
		Name:       tx.Name,
		EntryCount: 0,
	}
	key, _ := db.AccountKeyFromPK(ctx.Signer())
	account.Save(ctx.GetDB(), key)
	// return the new pk as response
	return tmsp.NewResultOK(key, "")
}
