package sign

import (
	wutil "github.com/ethanfrey/tenderize/wire"
	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
)

func init() {
	registerGoCryptoGoWire()
}

// we must register these types here, to make sure they parse (maybe go-wire issue??)
// TODO: fix go-wire, remove this code
func registerGoCryptoGoWire() {
	wire.RegisterInterface(
		struct{ crypto.PubKey }{},
		wire.ConcreteType{O: crypto.PubKeyEd25519{}, Byte: crypto.PubKeyTypeEd25519},
		wire.ConcreteType{O: crypto.PubKeySecp256k1{}, Byte: crypto.PubKeyTypeSecp256k1},
	)
	wire.RegisterInterface(
		struct{ crypto.Signature }{},
		wire.ConcreteType{O: crypto.SignatureEd25519{}, Byte: crypto.SignatureTypeEd25519},
		wire.ConcreteType{O: crypto.SignatureSecp256k1{}, Byte: crypto.SignatureTypeSecp256k1},
	)
}

// Action tries to limit the types we support to desired ones
type Action interface {
	IsAction() error
}

// actionWrapper is needed by go-wire to handle the interface
type actionWrapper struct {
	Action
}

// ActionToBytes converts the action into bytes to store in the db
// If there are invalid values in the action you can return an error
func ActionToBytes(action Action) ([]byte, error) {
	return wutil.ToBinary(actionWrapper{action})
}

// ActionFromBytes sets the action contents to the passed in data
// Returns error if the data doesn't match this action
func ActionFromBytes(data []byte) (Action, error) {
	holder := actionWrapper{}
	err := wutil.FromBinary(data, &holder)
	return holder.Action, err
}

// RegisterActions takes a list of all Actions we support and registers them with go-wire for Serialization
// The control byte is based on the order, so if you want to maintain compatibility with an
// existing data store, do not change the position of any items.  You can use nil as a
// placeholder to not use that byte anymore.
func RegisterActions(actions ...Action) {
	// prepare for max size, but might be shorter with nil
	regActs := make([]wire.ConcreteType, 0, len(actions))

	for i, action := range actions {
		// add the model and key for all non-nil values
		if action != nil {
			regActs = append(regActs, wire.ConcreteType{O: action, Byte: byte(i + 1)})
		}
	}
	// now register this with go-wire
	wire.RegisterInterface(actionWrapper{}, regActs...)
}

// ValidatedAction is returned after properly parsing a SignedAction and validating the signature
type ValidatedAction struct {
	action Action
	valid  bool
	SignedAction
}

// GetAction returns the action which was validated
func (v ValidatedAction) GetAction() Action {
	return v.action
}

// GetSigner returns the public key that signed the action, or nil if unvalidated
func (v ValidatedAction) GetSigner() crypto.PubKey {
	if !v.valid {
		return nil
	}
	return v.SignedAction.GetSigner()
}

// IsAnon returns false iff it was properly validated
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
	return wutil.ToBinary(tx)
}

// Deserialize will set the content of this SignedAction to the bytes on the wire
func (tx *SignedAction) Deserialize(data []byte) error {
	return wutil.FromBinary(data, tx)
}

// Validate will deserialize the contained action, and validate the signature or return an error
func (tx SignedAction) Validate() (ValidatedAction, error) {
	res := ValidatedAction{
		SignedAction: tx,
	}
	valid := tx.Signer.VerifyBytes(tx.ActionData, tx.Signature)
	if !valid {
		return res, errors.New("Invalid signature")
	}

	var err error
	res.action, err = ActionFromBytes(tx.ActionData)
	if err == nil {
		res.valid = true
	}
	return res, err
}

// SignAction will serialize the action and sign it with your key
func SignAction(action Action, privKey crypto.PrivKey) (res SignedAction, err error) {
	res.ActionData, err = ActionToBytes(action)
	if err != nil {
		return res, err
	}
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
