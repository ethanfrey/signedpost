package store

// Entry represents one verifiably immutable blog entry (so no typos ;)
type Entry struct {
	AccountPK []byte `json:"account"`
	Number    int    `json:"number"`

	PublishedBlock int64  `json:"published"`
	Title          string `json:"title"`
	Content        string `json:"content"`
}
