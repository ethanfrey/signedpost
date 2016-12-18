package txn

import "github.com/ethanfrey/tenderize/sign"

func init() {
	sign.RegisterActions(CreateAccountAction{}, AddPostAction{})
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

// AddPostAction is used for an existing account to append an entry to its list
type AddPostAction struct {
	Title   string
	Content string
}

// IsAction fulfills interface for go-wire
func (c AddPostAction) IsAction() error {
	return nil
}
