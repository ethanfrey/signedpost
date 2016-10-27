package models

import "time"

type CreateAccountAction struct {
	Name string `json:"name"` // this is a name to search for
}

type AddEntryAction struct {
	Published time.Time `json:"published"` // this must be close to the time it went to the blockchain
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}
