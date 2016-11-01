package view

// Account is the json object we return for one account
type Account struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PostCount int64  `json:"posts"`
}

// AccountList represent a list of accounts (from a search)
type AccountList struct {
	Items []*Account `json:"items"`
	Count int64      `json:"count"`
}

// Post is the json object we return for one post
type Post struct {
	ID             string `json:"id"`
	AccountID      string `json:"account"`
	Number         int64  `json:"number"`
	PublishedBlock uint64 `json:"published_block"`
	Title          string `json:"title"`
	Content        string `json:"content"`
}

// PostList represent a list of posts (for a user)
type PostList struct {
	Items []*Post `json:"items"`
	Count int64   `json:"count"`
}
