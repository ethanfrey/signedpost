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
	wire.RegisterInterface(
		struct{ SignedAction }{},
		wire.ConcreteType{FFignedAction{}, 0x01},
	)
	wire.RegisterInterface(
		struct{ ValidatedAction }{},
		wire.ConcreteType{validatedAction{}, 0x01},
	)
}

type SignedAction interface {
	GetSigner() crypto.PubKey
	IsAnon() bool
	GetActionData() []byte
	Validate() (ValidatedAction, error)
	Serialize() ([]byte, error)
}

type ValidatedAction interface {
	SignedAction
	GetAction() Action
}

// validatedAction is returned after properly parsing a SignedAction and validating the signature
type validatedAction struct {
	action Action
	valid  bool
	SignedAction
}

func (v validatedAction) GetAction() Action {
	return v.action
}

func (v validatedAction) GetSigner() crypto.PubKey {
	if !v.valid {
		return nil
	}
	return v.SignedAction.GetSigner()
}

func (v validatedAction) IsAnon() bool {
	return v.GetSigner() == nil
}

// FFignedAction contains a serialized action along with a signature for authorization
type FFignedAction struct {
	actionData []byte
	signature  crypto.Signature
	signer     crypto.PubKey
}

func (tx FFignedAction) GetSigner() crypto.PubKey {
	return tx.signer
}

func (tx FFignedAction) IsAnon() bool {
	return tx.signer == nil
}

func (tx FFignedAction) GetActionData() []byte {
	return tx.actionData
}

// Serialize gives a wire version of this action, reversed by Deserialize
func (tx FFignedAction) Serialize() ([]byte, error) {
	return utils.ToBinary(tx)
}

// Deserialize will set the content of this FFignedAction to the bytes on the wire
func (tx *FFignedAction) Deserialize(data []byte) error {
	return utils.FromBinary(data, tx)
}

// Validate will deserialize the contained action, and validate the signature or return an error
func (tx FFignedAction) Validate() (ValidatedAction, error) {
	action := validatedAction{
		SignedAction: tx,
	}
	valid := tx.signer.VerifyBytes(tx.actionData, tx.signature)
	if !valid {
		return action, errors.New("Invalid signature")
	}

	var n int
	var err error
	buf := bytes.NewBuffer(tx.actionData)
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
	res := FFignedAction{}
	buf := new(bytes.Buffer)

	wire.WriteBinary(actionWrapper{action}, buf, &n, &err)
	if err != nil {
		return res, errors.Wrap(err, "Sign Action")
	}

	res.actionData = buf.Bytes()
	res.signature = privKey.Sign(res.actionData)
	res.signer = privKey.PubKey()
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
	tx := FFignedAction{}
	err := tx.Deserialize(data)
	if err != nil {
		return nil, err
	}
	return tx.Validate()
}
