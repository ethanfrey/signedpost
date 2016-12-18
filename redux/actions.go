package redux

import (
	"github.com/ethanfrey/signedpost/store"
	"github.com/ethanfrey/signedpost/txn"
	"github.com/ethanfrey/tenderize/mom"
	crypto "github.com/tendermint/go-crypto"
	tmsp "github.com/tendermint/tmsp/types"
)

// CreateAccount creates a new account based on the signing public key
func (ctx *Service) CreateAccount(tx txn.CreateAccountAction, signer crypto.PubKey) tmsp.Result {
	if signer == nil {
		return tmsp.NewError(tmsp.CodeType_Unauthorized, "Must sign transaction")
	}

	// make sure none with this name or pk already....
	exists, err := store.FindAccount(ctx.GetDB(), signer)
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}
	if exists != nil {
		return tmsp.NewError(tmsp.CodeType_BaseDuplicateAddress,
			"Account exists for this public key")
	}

	matches, err := store.ListAccounts(ctx.GetDB(), store.AccountMatchesName(tx.Name))
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}
	if len(matches) > 0 {
		return tmsp.NewError(tmsp.CodeType_BaseDuplicateAddress,
			"Account name already taken")
	}

	// all safe, go save it
	account := store.NewAccount(signer, tx.Name)
	mom.Save(ctx.GetDB(), account)
	// return the new pk as response
	key, _ := mom.KeyToBytes(account.Key())
	return tmsp.NewResultOK(key, "")
}

// AppendPost adds a post to an existing account
func (ctx *Service) AppendPost(tx txn.AddPostAction, signer crypto.PubKey) tmsp.Result {
	if signer == nil {
		return tmsp.NewError(tmsp.CodeType_Unauthorized, "Must sign transaction")
	}

	// make sure we can find account for this user
	acct, err := store.FindAccount(ctx.GetDB(), signer)
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}
	if acct == nil {
		return tmsp.NewError(tmsp.CodeType_BaseUnknownAddress,
			"No account exists for this public key")
	}

	// fill out other info...
	num := acct.EntryCount + 1
	post := store.Post{
		Account:        acct.Key(),
		Title:          tx.Title,
		Content:        tx.Content,
		Number:         num,
		PublishedBlock: ctx.GetHeight(),
	}
	_, err = mom.Save(ctx.GetDB(), post)
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}

	// if saved, we must update account
	acct.EntryCount = num
	_, err = mom.Save(ctx.GetDB(), *acct)
	if err != nil {
		return tmsp.NewError(tmsp.CodeType_BaseInvalidInput, err.Error())
	}

	// return the post key as response
	key, _ := mom.KeyToBytes(post.Key())
	return tmsp.NewResultOK(key, "")
}
