package txn

import "github.com/ethanfrey/tenderize/sign"

func init() {
	sign.RegisterActions(CreateAccountAction{}, AddPostAction{})
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
