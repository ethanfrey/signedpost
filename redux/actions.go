package redux

import (
	"github.com/ethanfrey/signedpost/store"
	"github.com/ethanfrey/signedpost/txn"
	crypto "github.com/tendermint/go-crypto"
	tmsp "github.com/tendermint/tmsp/types"
)

// CreateAccount creates a new account based on the signing public key
func (ctx *Service) CreateAccount(tx txn.CreateAccountAction, signer crypto.PubKey) tmsp.Result {
	if signer == nil {
		return tmsp.NewError(tmsp.CodeType_Unauthorized, "Must sign transaction")
	}

	// make sure none with this name or pk already....
	exists, err := store.FindAccountByPK(ctx.GetDB(), signer)
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
	key, _ := store.AccountKeyFromPK(signer)
	account.Save(ctx.GetDB(), key)
	// return the new pk as response
	return tmsp.NewResultOK(key, "")
}
