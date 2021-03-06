package signedpost

import (
	"testing"

	"github.com/ethanfrey/signedpost/store"
	"github.com/ethanfrey/signedpost/txn"
	"github.com/ethanfrey/tenderize/mom"
	"github.com/ethanfrey/tenderize/sign"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crypto "github.com/tendermint/go-crypto"
	merkle "github.com/tendermint/go-merkle"
)

func TestApplication(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	earl := crypto.GenPrivKeyEd25519()
	tree := merkle.NewIAVLTree(0, nil) // in-memory

	app := NewApp(tree)
	// make sure initial hash is nil
	assert.Nil(app.Commit().Data)

	// TODO: use BeginBlock(2), when that works...
	app.EndBlock(1)
	utx := txn.CreateAccountAction{Name: "Grey"}
	data, err := sign.Send(utx, earl)
	require.Nil(err, "%+v", err)
	require.NotNil(data)
	ures := app.AppendTx(data)
	assert.False(ures.IsErr(), ures.Error())
	ukey := ures.Data

	// make sure commit hash is updated
	hash := app.Commit().Data
	assert.NotEqual(nil, hash)
	// make sure we can query
	qres := app.Query(ukey)
	assert.False(qres.IsErr(), qres.Error())
	model, err := mom.ModelFromBytes(qres.Data)
	if assert.Nil(err, "%+v", err) {
		acct, ok := model.(store.Account)
		if assert.True(ok) {
			assert.Equal(acct.Name, "Grey")
		}
	}

	// now add the post
	ptx := txn.AddPostAction{
		Title:   "Good post",
		Content: "Some imporant info",
	}
	pdata, err := sign.Send(ptx, earl)
	require.Nil(err, "%+v", err)
	require.NotNil(pdata)
	// make sure check works, but doesn't update data
	pres := app.CheckTx(pdata)
	assert.False(pres.IsErr(), pres.Error())
	assert.Equal(hash, app.Commit().Data)

	// now, really append
	pres = app.AppendTx(pdata)
	assert.False(pres.IsErr(), pres.Error())
	assert.NotEqual(hash, app.Commit().Data)
}
