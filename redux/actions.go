package redux

import (
	"github.com/ethanfrey/signedpost/store"
	"github.com/ethanfrey/signedpost/txn"
	tmsp "github.com/tendermint/tmsp/types"
)

// CreateAccount creates a new account based on the signing public key
func CreateAccount(ctx *Service, tx txn.CreateAccountAction) tmsp.Result {
	if ctx.IsAnon() {
		return tmsp.NewError(tmsp.CodeType_Unauthorized, "Must sign transaction")
	}

	// make sure none with this name or pk already....
	exists, err := store.FindAccountByPK(ctx.GetDB(), ctx.Signer())
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}
	if exists != nil {
		return tmsp.NewError(tmsp.CodeType_BaseDuplicateAddress,
			"Account exists for this public key")
	}

	exists, err = store.FindAccountByName(ctx.GetDB(), tx.Name)
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}
	if exists != nil {
		return tmsp.NewError(tmsp.CodeType_BaseDuplicateAddress,
			"Account name already taken")
	}

	// all safe, go save it
	account := store.Account{
		Name:       tx.Name,
		EntryCount: 0,
	}
	key, _ := store.AccountKeyFromPK(ctx.Signer())
	account.Save(ctx.GetDB(), key)
	// return the new pk as response
	return tmsp.NewResultOK(key, "")
}
