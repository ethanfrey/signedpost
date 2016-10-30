package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crypto "github.com/tendermint/go-crypto"
	merkle "github.com/tendermint/go-merkle"
)

func TestPost(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	priv := crypto.GenPrivKeyEd25519()
	pub := priv.PubKey()
	acct, err := AccountKeyFromPK(pub)
	require.Nil(err)

	tree := merkle.NewIAVLTree(0, nil) // in-memory
	assert.Equal(0, tree.Size())

	// create an account
	a := Account{Name: "Demo"}
	_, err = a.Serialize()
	require.Nil(err)
	updated, err := a.Save(tree, acct)
	assert.False(updated)
	require.Nil(err)

	// make sure we don't find any posts, or a named post
	posts, err := FindPostsForAccount(tree, acct)
	require.Nil(err)
	assert.Equal(0, len(posts))
	first, err := FindPostByAcctNum(tree, acct, 1)
	require.Nil(err)
	assert.Nil(first)

	// let's add one
	p := Post{
		Number:  1,
		Title:   "First Post",
		Content: "Some verified text, please",
	}
	key, err := PostKeyFromAccount(acct, 1)
	require.Nil(err, "%+v", err)
	updated, err = p.Save(tree, key)
	require.Nil(err, "%+v", err)
	assert.False(updated)

	// make sure we find it, or a named post
	posts, err = FindPostsForAccount(tree, acct)
	require.Nil(err)
	assert.Equal(1, len(posts))
	assertPost(t, &p, posts[0])
	first, err = FindPostByAcctNum(tree, acct, 1)
	require.Nil(err)
	assertPost(t, &p, first)
}

func assertPost(t *testing.T, post, match *Post) {
	assert := assert.New(t)
	if assert.NotNil(post) && assert.NotNil(match) {
		assert.Equal(post.Title, match.Title)
		assert.Equal(post.Number, match.Number)
		assert.Equal(post.Content, match.Content)
	}
}
