package models

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
)

func init() {
	// we must register these types here, to make sure they parse (maybe go-wire issue??)
	wire.RegisterInterface(
		struct{ crypto.PubKey }{},
		wire.ConcreteType{crypto.PubKeyEd25519{}, crypto.PubKeyTypeEd25519},
		wire.ConcreteType{crypto.PubKeySecp256k1{}, crypto.PubKeyTypeSecp256k1},
	)
	wire.RegisterInterface(
		struct{ crypto.Signature }{},
		wire.ConcreteType{crypto.SignatureEd25519{}, crypto.SignatureTypeEd25519},
		wire.ConcreteType{crypto.SignatureSecp256k1{}, crypto.SignatureTypeSecp256k1},
	)

}

// ValidatedAction is returned after properly parsing a SignedAction and validating the signature
type ValidatedAction struct {
	Action Action
	Signer crypto.PubKey
}

// SignedAction contains a serialized action along with a signature for authorization
type SignedAction struct {
	Data      []byte
	Signature crypto.Signature
	Signer    struct{ crypto.PubKey }
}

// Serialize gives a wire version of this action, reversed by Deserialize
func (tx SignedAction) Serialize() ([]byte, error) {
	return ToBinary(tx)
}

// Deserialize will set the content of this SignedAction to the bytes on the wire
func (tx *SignedAction) Deserialize(data []byte) error {
	return FromBinary(data, tx)
}

// Validate will deserialize the contained action, and validate the signature or return an error
func (tx SignedAction) Validate() (ValidatedAction, error) {
	action := ValidatedAction{}
	valid := tx.Signer.VerifyBytes(tx.Data, tx.Signature)
	if !valid {
		return action, errors.New("Invalid signature")
	}

	var n int
	var err error
	buf := bytes.NewBuffer(tx.Data)
	res := wire.ReadBinary(actionWrapper{}, buf, 0, &n, &err)
	if err != nil {
		return action, errors.Wrap(err, "Parsing valid action")
	}

	action.Signer = tx.Signer
	wrap, ok := res.(actionWrapper)
	if !ok {
		return action, errors.New("Data is not action type")
	}
	action.Action = wrap.Action
	return action, nil
}

// SignAction will serialize the action and sign it with your key
func SignAction(action Action, privKey crypto.PrivKey) (SignedAction, error) {
	var n int
	var err error
	res := SignedAction{}
	buf := new(bytes.Buffer)

	wire.WriteBinary(actionWrapper{action}, buf, &n, &err)
	if err != nil {
		return res, errors.Wrap(err, "Sign Action")
	}

	res.Data = buf.Bytes()
	res.Signature = privKey.Sign(res.Data)
	res.Signer = struct{ crypto.PubKey }{privKey.PubKey()}
	return res, nil
}
