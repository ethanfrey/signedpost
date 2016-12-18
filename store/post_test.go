package store

import (
	"testing"

	"github.com/ethanfrey/tenderize/mom"
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

	tree := merkle.NewIAVLTree(0, nil) // in-memory
	assert.Equal(0, tree.Size())

	// create an account
	acct := NewAccount(pub, "Fred")
	updated, err := mom.Save(tree, acct)
	require.Nil(err)

	// make sure we don't find any posts, or a named post
	myPosts := PostsForAccount(acct, 0)
	firstPost := PostsForAccount(acct, 1)

	posts, err := ListPosts(tree, myPosts, nil)
	require.Nil(err)
	assert.Equal(0, len(posts))

	first, err := ListPosts(tree, firstPost, nil)
	require.Nil(err)
	assert.Equal(0, len(first))

	// let's add two
	p := Post{
		Account: acct.Key(),
		Number:  1,
		Title:   "First Post",
		Content: "Some verified text, please",
	}
	updated, err = mom.Save(tree, p)
	require.Nil(err, "%+v", err)
	assert.False(updated)

	p2 := Post{
		Account: acct.Key(),
		Number:  2,
		Title:   "Something else",
		Content: "Now we have some data!",
	}
	updated, err = mom.Save(tree, p2)
	require.Nil(err, "%+v", err)
	assert.False(updated)

	// make sure we find it, or a named post
	posts, err = ListPosts(tree, myPosts, nil)
	require.Nil(err)
	if assert.Equal(2, len(posts)) {
		assertPost(t, p, posts[0])
		assertPost(t, p2, posts[1])
	}

	first, err = ListPosts(tree, firstPost, nil)
	require.Nil(err)
	if assert.Equal(1, len(first)) {
		assertPost(t, p, first[0])
	}
}

func assertPost(t *testing.T, post Post, match Post) {
	assert := assert.New(t)
	assert.EqualValues(post.Account, match.Account)
	assert.Equal(post.Title, match.Title)
	assert.Equal(post.Number, match.Number)
	assert.Equal(post.Content, match.Content)
}
