package txn

import (
	"bytes"

	"github.com/ethanfrey/signedpost/utils"
	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
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

type Signed interface {
	GetSigner() crypto.PubKey
	IsAnon() bool
	GetActionData() []byte
	Validate() (ValidatedAction, error)
	Serialize() ([]byte, error)
}

type Validated interface {
	Signed
	GetAction() Action
}

// ValidatedAction is returned after properly parsing a SignedAction and validating the signature
type ValidatedAction struct {
	action Action
	valid  bool
	SignedAction
}

func (v ValidatedAction) GetAction() Action {
	return v.action
}

func (v ValidatedAction) GetSigner() crypto.PubKey {
	if !v.valid {
		return nil
	}
	return v.SignedAction.GetSigner()
}

func (v ValidatedAction) IsAnon() bool {
	return v.GetSigner() == nil
}

// SignedAction contains a serialized action along with a signature for authorization
type SignedAction struct {
	ActionData []byte
	Signature  crypto.Signature
	Signer     crypto.PubKey
}

func (tx SignedAction) GetSigner() crypto.PubKey {
	return tx.Signer
}

func (tx SignedAction) IsAnon() bool {
	return tx.Signer == nil
}

func (tx SignedAction) GetActionData() []byte {
	return tx.ActionData
}

// Serialize gives a wire version of this action, reversed by Deserialize
func (tx SignedAction) Serialize() ([]byte, error) {
	return utils.ToBinary(tx)
}

// Deserialize will set the content of this SignedAction to the bytes on the wire
func (tx *SignedAction) Deserialize(data []byte) error {
	return utils.FromBinary(data, tx)
}

// Validate will deserialize the contained action, and validate the signature or return an error
func (tx SignedAction) Validate() (ValidatedAction, error) {
	action := ValidatedAction{
		SignedAction: tx,
	}
	valid := tx.Signer.VerifyBytes(tx.ActionData, tx.Signature)
	if !valid {
		return action, errors.New("Invalid signature")
	}

	var n int
	var err error
	buf := bytes.NewBuffer(tx.ActionData)
	res := wire.ReadBinary(actionWrapper{}, buf, 0, &n, &err)
	if err != nil {
		return action, errors.Wrap(err, "Parsing valid action")
	}

	wrap, ok := res.(actionWrapper)
	if !ok {
		return action, errors.New("Data is not action type")
	}
	action.action = wrap.Action
	action.valid = true
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

	res.ActionData = buf.Bytes()
	res.Signature = privKey.Sign(res.ActionData)
	res.Signer = privKey.PubKey()
	return res, nil
}

// Send will take the action and owner and prepare bytes
func Send(action Action, privKey crypto.PrivKey) ([]byte, error) {
	tx, err := SignAction(action, privKey)
	if err != nil {
		return nil, err
	}
	return tx.Serialize()
}

// Receive will take some bytes, parse them, and validate the signature
func Receive(data []byte) (ValidatedAction, error) {
	tx := SignedAction{}
	err := tx.Deserialize(data)
	if err != nil {
		return ValidatedAction{}, err
	}
	return tx.Validate()
}
