package redux

import (
	"testing"

	"github.com/ethanfrey/signedpost/store"
	"github.com/ethanfrey/signedpost/txn"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"
	merkle "github.com/tendermint/go-merkle"
)

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)
	alice := crypto.GenPrivKeyEd25519()
	bob := crypto.GenPrivKeyEd25519()
	tree := merkle.NewIAVLTree(0, nil) // in-memory
	assert.Equal(0, tree.Size())

	srv := &Service{
		store: tree,
	}
	tx := txn.CreateAccountAction{Name: "Alice"}

	// anon is prevented
	r := srv.CreateAccount(tx, nil)
	assert.True(r.IsErr(), "%+v", r.Code)
	assert.Equal(0, tree.Size())

	// success for self-creation
	r = srv.CreateAccount(tx, alice.PubKey())
	assert.False(r.IsErr(), r.Error())
	assert.Equal(1, tree.Size())

	// let's check this account by key
	data, err := store.FindAccountByPK(tree, alice.PubKey())
	assert.Nil(err)
	if assert.NotNil(data) {
		assert.Equal(data.Name, "Alice")
	}

	// let's check this account by name
	data, err = store.FindAccountByName(tree, "Alice")
	assert.Nil(err)
	if assert.NotNil(data) {
		assert.Equal(data.Name, "Alice")
	}

	// error by second name
	tx2 := txn.CreateAccountAction{Name: "Bob"}
	r = srv.CreateAccount(tx, alice.PubKey())
	assert.True(r.IsErr(), "%+v", r.Code)

	// cannot claim the same name (taken)
	r = srv.CreateAccount(tx, bob.PubKey())
	assert.True(r.IsErr(), "%+v", r.Code)
	// but he can claim his own name
	r = srv.CreateAccount(tx2, bob.PubKey())
	assert.False(r.IsErr(), r.Error())

	// TODO: add queries
}
