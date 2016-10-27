package models

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
	assert.NotNil(signed.Data)
	assert.True(len(signed.Data) > 4) // It must contain at least "John"
	assert.NotNil(signed.Signature)
	assert.NotNil(signed.Signer)

	// does this action validate?
	valid, err := ValidateAction(signed)
	require.Nil(err, "%+v", err)
	assert.Equal(signed.Signer, valid.Signer)
	assert.NotNil(valid.Action)
	ca, ok := valid.Action.(CreateAccountAction)
	if assert.True(ok) {
		assert.Equal(action.Name, ca.Name)
	} else {
		fmt.Println("Let's try a pointer")
		captr, ok := valid.Action.(*CreateAccountAction)
		if assert.True(ok) && assert.NotNil(captr) {
			assert.Equal(action.Name, captr.Name)
		}
	}
}
