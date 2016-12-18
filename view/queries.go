package view

import (
	"github.com/pkg/errors"

	"github.com/ethanfrey/signedpost/store"
	"github.com/ethanfrey/tenderize/mom"
	merkle "github.com/tendermint/go-merkle"
)

// AllAccounts returns what you expect
func AllAccounts(tree merkle.Tree) (*AccountList, error) {
	accts, err := store.ListAccounts(tree, nil)
	if err != nil {
		return nil, err
	}
	return RenderAccountList(accts), nil
}

// AccountByKey returns an exact match
func AccountByKey(tree merkle.Tree, key []byte) (*Account, error) {
	model, err := mom.Load(tree, store.AccountKey{ID: key})
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, errors.New("Not Found")
	}
	return RenderAccount(model.(store.Account)), nil
}

// AccountByName searches for similar names
func AccountByName(tree merkle.Tree, name string) (*AccountList, error) {
	accts, err := store.ListAccounts(tree, store.AccountContainsName(name))
	if err != nil {
		return nil, err
	}
	return RenderAccountList(accts), nil
}

// PostsForAccount returns all posts that belong to this account
func PostsForAccount(tree merkle.Tree, acct []byte) (*PostList, error) {
	key := store.PostKey{Account: store.AccountKey{ID: acct}}
	posts, err := store.ListPosts(tree, key, nil)
	if err != nil {
		return nil, err
	}
	return RenderPostList(posts), nil
}

// PostByKey returns an exact match
func PostByKey(tree merkle.Tree, key []byte) (*Post, error) {
	postKey, err := mom.KeyFromBytes(key)
	if err != nil {
		return nil, err
	}
	posts, err := store.ListPosts(tree, postKey, nil)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, errors.New("Not Found")
	}
	return RenderPost(posts[0]), nil
}
