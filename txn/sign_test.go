package txn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-crypto"
)

func TestSignVerify(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	privKey := crypto.GenPrivKeyEd25519()

	// let's sign the action and make sure that works
	action := CreateAccountAction{Name: "John"}
	signed, err := SignAction(action, privKey)
	require.Nil(err)
	assert.NotNil(signed.GetActionData())
	assert.True(len(signed.GetActionData()) > 4) // It must contain at least "John"
	assert.NotNil(signed.GetSigner())
	assert.False(signed.IsAnon())

	// does this action validate?
	valid, err := signed.Validate()
	require.Nil(err, "%+v", err)
	assert.Equal(signed.GetSigner(), valid.GetSigner())
	act := valid.GetAction()
	assert.NotNil(act)
	ca, ok := act.(CreateAccountAction)
	if assert.True(ok) {
		assert.Equal(action.Name, ca.Name)
	} else {
		fmt.Println("Let's try a pointer")
		captr, ok := act.(*CreateAccountAction)
		if assert.True(ok) && assert.NotNil(captr) {
			assert.Equal(action.Name, captr.Name)
		}
	}
}

func TestSignSerialization(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	privKey := crypto.GenPrivKeyEd25519()

	// let's sign the action and make sure that works
	action := AddEntryAction{Title: "First Post", Content: "Some text here"}
	signed, err := SignAction(action, privKey)
	require.Nil(err)

	wire, err := signed.Serialize()
	require.Nil(err, "%+v", err)
	require.Equal(129, len(wire))

	// make sure the data is there
	parsed, err := Receive(wire)
	require.Nil(err, "%+v", err)
	assert.Equal(signed.GetActionData(), parsed.GetActionData())
	assert.Equal(signed.GetSigner(), parsed.GetSigner())

	// serialize a second object and make sure the same wire
	a2 := AddEntryAction{Title: "First Post"}
	wire2, err := Send(a2, privKey)
	require.Nil(err, "%+v", err)
	assert.NotEqual(wire, wire2)
	// 14 chars less, means shorter data (why 15?)
	assert.Equal(len(wire)-15, len(wire2))

	a2.Content = "Some text here"
	wire3, err := Send(a2, privKey)
	require.Nil(err, "%+v", err)
	assert.Equal(wire, wire3)
}
