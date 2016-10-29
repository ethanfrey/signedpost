package db

import (
	"testing"

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
	key, err := AccountKeyFromPK(pub)
	require.Nil(err)

	tree := merkle.NewIAVLTree(0, nil) // in-memory
	assert.Equal(0, tree.Size())

	// make sure empty searches work as expected
	match, err := FindAccountByPK(tree, pub)
	assert.Nil(match)
	assert.Nil(err)
	match, err = FindAccountByKey(tree, key)
	assert.Nil(match)
	assert.Nil(err)
	match, err = FindAccountByKey(tree, []byte("foobar"))
	assert.Nil(match)
	// assert.NotNil(err)
	match, err = FindAccountByName(tree, "Demo")
	assert.Nil(match)
	assert.Nil(err)

	acct := Account{Name: "Demo"}
	_, err = acct.Serialize()
	require.Nil(err)

	// on set
	updated, err := acct.Save(tree, key)
	assert.False(updated)
	assert.Nil(err)

	// update proper
	acct.EntryCount = 2
	updated, err = acct.Save(tree, key)
	assert.True(updated)
	assert.Nil(err)

	// cannot save to invalid address
	_, err = acct.Save(tree, nil)
	assert.NotNil(err)

	// Now we search....
	match, err = FindAccountByPK(tree, pub)
	assert.Nil(err)
	assertAccount(t, &acct, match)
	match, err = FindAccountByKey(tree, key)
	assert.Nil(err)
	assertAccount(t, &acct, match)
	match, err = FindAccountByName(tree, "Demo")
	assert.Nil(err)
	assertAccount(t, &acct, match)
}

func assertAccount(t *testing.T, acct, match *Account) {
	assert := assert.New(t)
	if assert.NotNil(acct) && assert.NotNil(match) {
		assert.Equal(acct.Name, match.Name)
		assert.Equal(acct.EntryCount, match.EntryCount)
	}
}
