package models

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
)

// TODO: easier was to serialize with go-wire....

// SignedAction contains a serialized action along with a signature for authorization
type SignedAction struct {
	Data      []byte
	Signature crypto.Signature
	Signer    crypto.PubKey
}

// ValidatedAction is returned after properly parsing a SignedAction and validating the signature
type ValidatedAction struct {
	Action Action
	Signer crypto.PubKey
}

// SignAction will serialize the action and sign it with your key
func SignAction(action Action, privKey crypto.PrivKey) (SignedAction, error) {
	var n int
	var err error
	res := SignedAction{}
	buf := new(bytes.Buffer)

	wire.WriteBinary(ActionWrap{action}, buf, &n, &err)
	if err != nil {
		return res, errors.Wrap(err, "Sign Action")
	}

	res.Data = buf.Bytes()
	res.Signature = privKey.Sign(res.Data)
	res.Signer = privKey.PubKey()
	return res, nil
}

// ValidateAction will deserialize the action, and validate the signature or return an error
func ValidateAction(tx SignedAction) (ValidatedAction, error) {
	action := ValidatedAction{}
	valid := tx.Signer.VerifyBytes(tx.Data, tx.Signature)
	if !valid {
		return action, errors.New("Invalid signature")
	}

	var n int
	var err error
	buf := bytes.NewBuffer(tx.Data)
	res := wire.ReadBinary(ActionWrap{}, buf, 0, &n, &err)
	if err != nil {
		return action, errors.Wrap(err, "Parsing valid action")
	}

	action.Signer = tx.Signer
	wrap, ok := res.(ActionWrap)
	if !ok {
		return action, errors.New("Data is not action type")
	}
	action.Action = wrap.Action
	return action, nil
}
