package db

import (
	"time"
)

// Entry represents one verifiably immutable blog entry (so no typos ;)
type Entry struct {
	AccountPK []byte `json:"account"`
	Number    int    `json:"number"`

	Published time.Time `json:"published"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}
