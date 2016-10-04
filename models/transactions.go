package models

import "time"

type CreateAccountAction struct {
	Name   string `json:"name"`   // this is a name to search for
	Public bool   `json:"public"` // if set to false, only the owner can read blog
}

type TogglePublicAction struct {
	Public bool `json:"public"` // new state of public flag
}

type AddEntryAction struct {
	Published time.Time `json:"published"` // this must be close to the time it went to the blockchain
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}
