package txn

import "github.com/tendermint/go-wire"

func init() {
	wire.RegisterInterface(
		actionWrapper{},
		wire.ConcreteType{O: CreateAccountAction{}, Byte: 0x01},
		wire.ConcreteType{O: &CreateAccountAction{}, Byte: 0x02},
		wire.ConcreteType{O: AddEntryAction{}, Byte: 0x03},
		wire.ConcreteType{O: &AddEntryAction{}, Byte: 0x04},
	)
}

// actionWrapper is needed by go-wire to handle the interface
type actionWrapper struct {
	Action
}

// Action tries to limit the types we support to desired ones
type Action interface {
	IsAction() error
}

// CreateAccountAction is used once to claim a username for a given public key
type CreateAccountAction struct {
	Name string // this is a name to search for
}

// IsAction fulfills interface for go-wire
func (c CreateAccountAction) IsAction() error {
	return nil
}

// AddEntryAction is used for an existing account to append an entry to its list
type AddEntryAction struct {
	Title   string
	Content string
}

// IsAction fulfills interface for go-wire
func (c AddEntryAction) IsAction() error {
	return nil
}
