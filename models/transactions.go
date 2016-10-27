package models

import "time"

// CreateAccountAction is used once to claim a username for a given public key
type CreateAccountAction struct {
	Name string `json:"name"` // this is a name to search for
}

// AddEntryAction is used for an existing account to append an entry to its list
type AddEntryAction struct {
	Published time.Time `json:"published"` // this must be close to the time it went to the blockchain
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}
