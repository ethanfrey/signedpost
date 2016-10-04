package db

// Account is a named account that can publish blog entries
type Account struct {
	PK         []byte `json:"pk"`          // this is the public key of the owner
	Name       string `json:"name"`        // this is a name to search for
	EntryCount int    `json:"num_entries"` // total number of entries (de-normalize for speed)
	Public     bool   `json:"public"`      // if set to false, only the owner can read blog
}
