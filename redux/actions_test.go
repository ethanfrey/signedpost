package redux

import (
	"testing"

	"github.com/ethanfrey/signedpost/store"
	"github.com/ethanfrey/signedpost/txn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestAppendPost(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	alice := crypto.GenPrivKeyEd25519()
	pub := alice.PubKey()
	tree := merkle.NewIAVLTree(0, nil) // in-memory
	assert.Equal(0, tree.Size())

	srv := &Service{
		store:       tree,
		blockHeight: 2,
	}

	tx := txn.AddPostAction{
		Title:   "My First Post",
		Content: "some data",
	}

	// anon is prevented
	r := srv.AppendPost(tx, nil)
	assert.True(r.IsErr(), "%+v", r.Code)
	assert.Equal(0, tree.Size())

	// append with un-registered account is prevented
	r = srv.AppendPost(tx, pub)
	assert.True(r.IsErr(), "%+v", r.Code)
	assert.Equal(0, tree.Size())

	// success for self-creation
	utx := txn.CreateAccountAction{Name: "Alice"}
	r = srv.CreateAccount(utx, pub)
	assert.False(r.IsErr(), r.Error())
	assert.Equal(1, tree.Size())
	acctKey := r.Data

	// now, let's add a post...
	r = srv.AppendPost(tx, pub)
	assert.False(r.IsErr(), "%+v", r.Error())
	assert.Equal(2, tree.Size())
	postKey := r.Data

	// let's check the post
	pp, err := store.FindPostByKey(tree, postKey)
	require.Nil(err, "%+v", err)
	require.NotNil(pp)
	assert.Equal(tx.Title, pp.Title)
	assert.Equal(srv.GetHeight(), pp.PublishedBlock)
	assert.Equal(1, pp.Number)

	// get the account and check it was updated
	aa, err := store.FindAccountByKey(tree, acctKey)
	assert.Nil(err)
	if assert.NotNil(aa) {
		assert.Equal("Alice", aa.Name)
		assert.Equal(1, aa.EntryCount)
	}

	// add a second post and make sure we query both
	tx2 := txn.AddPostAction{
		Title:   "Quick Update",
		Content: "We can add multiple posts",
	}
	r = srv.AppendPost(tx2, pub)
	assert.False(r.IsErr(), "%+v", r.Error())
	assert.Equal(3, tree.Size())

	// get the account and check it was updated
	aa, err = store.FindAccountByKey(tree, acctKey)
	assert.Nil(err)
	if assert.NotNil(aa) {
		assert.Equal("Alice", aa.Name)
		assert.Equal(2, aa.EntryCount)
	}

	// let's check the post
	posts, err := store.FindPostsForAccount(tree, acctKey)
	require.Nil(err, "%+v", err)
	require.Equal(2, len(posts))
	assert.Equal(tx.Title, posts[0].Title)
	assert.Equal(1, posts[0].Number)
	assert.Equal(tx2.Title, posts[1].Title)
	assert.Equal(2, posts[1].Number)
}
