package store

import (
	"math"

	"github.com/ethanfrey/tenderize/mom"
	"github.com/tendermint/go-merkle"
)

// Post represents one verifiably immutable blog entry (so no typos ;)
type Post struct {
	Account        mom.Key
	Number         int64
	PublishedBlock uint64
	Title          string
	Content        string
}

// PostKey is the index of this Post structure
type PostKey struct {
	Account mom.Key
	Number  int64
}

// Key returns the index of the Post (account, number)
func (p Post) Key() mom.Key {
	return PostKey{
		Account: p.Account,
		Number:  p.Number,
	}
}

// Range contains all posts if nothing set in PostKey, all posts for one account if Account set, but not number
func (p PostKey) Range() (mom.Key, mom.Key) {
	// TODO: make this a bit cleaner?
	min, max := p, p
	min.Account, max.Account = p.Account.Range()

	if p.Number == 0 {
		min.Number = 1
		max.Number = math.MaxInt32
	}
	return min, max
}

// PostsForAccount returns a range key for this account. if number is not 0, only return that post
func PostsForAccount(acct Account, number int64) mom.Key {
	return Post{Account: acct.Key(), Number: number}.Key()
}

// ListPosts makes a search over all accounts, and casts them to the proper type
// note an empty response returns no error
func ListPosts(store merkle.Tree, key mom.Key, filter func(mom.Model) bool) ([]Post, error) {
	query := mom.Query{
		Key:    key,
		Filter: filter,
	}
	models, err := mom.List(store, query)
	if err != nil {
		return nil, err
	}
	res := make([]Post, len(models))
	for i := range models {
		res[i] = models[i].(Post)
	}
	return res, nil
}
