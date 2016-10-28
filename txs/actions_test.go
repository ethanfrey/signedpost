package txs

import (
	"testing"

	"github.com/ethanfrey/signedpost/db"
	"github.com/ethanfrey/signedpost/models"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"
	merkle "github.com/tendermint/go-merkle"
)

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)
	alice := crypto.GenPrivKeyEd25519()
	bob := crypto.GenPrivKeyEd25519()
	tree := merkle.NewIAVLTree(0, nil) // in-memory

	anon := &Context{
		store: tree,
	}
	tx := models.CreateAccountAction{Name: "Alice"}

	// anon is prevented
	r := CreateAccount(anon, tx)
	assert.True(r.IsErr(), "%+v", r.Code)

	// success for self-creation
	ctx := &Context{
		store:  tree,
		signer: alice.PubKey(),
	}
	r = CreateAccount(ctx, tx)
	assert.False(r.IsErr(), r.Error())

	// let's check this account by key
	data, err := db.FindAccountByPK(tree, ctx.Signer().Address())
	assert.Nil(err)
	if assert.NotNil(data) {
		assert.Equal(data.Name, "Alice")
	}

	// let's check this account by name
	data, err = db.FindAccountByName(tree, "Alice")
	assert.Nil(err)
	if assert.NotNil(data) {
		assert.Equal(data.Name, "Alice")
	}

	// error by second name
	tx2 := models.CreateAccountAction{Name: "Bob"}
	r = CreateAccount(ctx, tx)
	assert.True(r.IsErr(), "%+v", r.Code)

	// but bob can make a new account
	ctx2 := &Context{
		store:  tree,
		signer: bob.PubKey(),
	}
	// cannot claim the same name (taken)
	r = CreateAccount(ctx2, tx)
	assert.True(r.IsErr(), "%+v", r.Code)
	// but he can claim his own name
	r = CreateAccount(ctx2, tx2)
	assert.False(r.IsErr(), r.Error())

	// TODO: add queries
}
