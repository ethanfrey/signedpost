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

// CreateAccount creates a new account based on the signing public key
func (ctx *Service) AppendPost(tx txn.AddPostAction, signer crypto.PubKey) tmsp.Result {
	if signer == nil {
		return tmsp.NewError(tmsp.CodeType_Unauthorized, "Must sign transaction")
	}

	// make sure we can find account for this user
	acct, err := store.FindAccountByPK(ctx.GetDB(), signer)
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}
	if acct == nil {
		return tmsp.NewError(tmsp.CodeType_BaseUnknownAddress,
			"No account exists for this public key")
	}
	acctKey, _ := store.AccountKeyFromPK(signer)

	// fill out other info...
	num := acct.EntryCount + 1
	post := store.Post{
		Title:          tx.Title,
		Content:        tx.Content,
		Number:         num,
		PublishedBlock: ctx.GetHeight(),
	}
	key, _ := store.PostKeyFromAccount(acctKey, num)
	_, err = post.Save(ctx.GetDB(), key)
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}

	// if saved, we must update account
	acct.EntryCount = num
	acct.Save(ctx.GetDB(), acctKey)

	// return the post key as response
	return tmsp.NewResultOK(key, "")
}
