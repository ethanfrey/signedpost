package store

import (
	"bytes"
	"testing"

	"github.com/ethanfrey/tenderize/mom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crypto "github.com/tendermint/go-crypto"
	merkle "github.com/tendermint/go-merkle"
)

func TestAccount(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	priv := crypto.GenPrivKeyEd25519()
	pub := priv.PubKey()
	badpriv := crypto.GenPrivKeyEd25519()
	badpub := badpriv.PubKey()

	tree := merkle.NewIAVLTree(0, nil) // in-memory
	assert.Equal(0, tree.Size())

	// make sure empty searches work as expected
	match, err := FindAccount(tree, pub)
	assert.Nil(match)
	assert.Nil(err)
	matches, err := ListAccounts(tree, AccountMatchesName("Demo"))
	assert.Equal(0, len(matches))
	assert.Nil(err)

	// on set
	acct := NewAccount(pub, "Demo")
	updated, err := mom.Save(tree, acct)
	assert.False(updated)
	require.Nil(err)

	// on update
	acct.Name = "Demoed"
	updated, err = mom.Save(tree, acct)
	assert.True(updated)
	require.Nil(err)

	// TODO: add some more checks to tenderize for this case?
	// // cannot save to invalid address
	// _, err = mom.Save(tree, Account{})
	// assert.NotNil(err)

	// make sure it is stores under one key
	match, err = FindAccount(tree, pub)
	assert.Nil(err)
	assertAccount(t, acct, *match)
	match, err = FindAccount(tree, badpub)
	assert.Nil(err)
	assert.Nil(match)

	// raw mom query
	query := mom.Query{
		Key: AccountKey{},
	}
	models, err := mom.List(tree, query)
	assert.Nil(err)
	assert.Equal(1, len(models))

	// and try a few searches
	// all account...
	matches, err = ListAccounts(tree, nil)
	assert.Nil(err)
	if assert.Equal(1, len(matches)) {
		assertAccount(t, acct, matches[0])
	}

	// Exact match on name
	matches, err = ListAccounts(tree, AccountMatchesName("Demoed"))
	assert.Nil(err)
	if assert.Equal(1, len(matches)) {
		assertAccount(t, acct, matches[0])
	}

	// contains match on substring
	matches, err = ListAccounts(tree, AccountContainsName("deMO"))
	assert.Nil(err)
	if assert.Equal(1, len(matches)) {
		assertAccount(t, acct, matches[0])
	}

	// exact match on substring
	matches, err = ListAccounts(tree, AccountMatchesName("Demo"))
	assert.Nil(err)
	assert.Equal(0, len(matches))
}

func assertAccount(t *testing.T, acct Account, match Account) {
	assert := assert.New(t)
	assert.True(bytes.Equal(acct.ID, match.ID))
	assert.Equal(acct.Name, match.Name)
	assert.Equal(acct.EntryCount, match.EntryCount)
}

func TestMultipleAccounts(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tree := merkle.NewIAVLTree(0, nil) // in-memory
	assert.Equal(0, tree.Size())

	alice := makeAccount(t, tree, "Alice")
	bob := makeAccount(t, tree, "Bob")

	assert.Equal(2, tree.Size())

	accts, err := ListAccounts(tree, nil)
	require.Nil(err)
	require.Equal(2, len(accts))

	ai, bi := 0, 1
	if bytes.Compare(alice.ID, bob.ID) > 0 {
		ai, bi = 1, 0
	}
	assertAccount(t, alice, accts[ai])
	assertAccount(t, bob, accts[bi])

	// and one more makes three...
	makeAccount(t, tree, "Carl")
	accts, err = ListAccounts(tree, nil)
	require.Nil(err)
	require.Equal(3, len(accts))
}

func makeAccount(t *testing.T, tree merkle.Tree, name string) (acct Account) {
	assert := assert.New(t)
	pub := crypto.GenPrivKeyEd25519().PubKey()
	acct = NewAccount(pub, name)

	updated, err := mom.Save(tree, acct)
	assert.False(updated)
	assert.Nil(err)

	return acct
}
