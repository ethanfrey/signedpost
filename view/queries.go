package view

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/ethanfrey/signedpost/store"
	merkle "github.com/tendermint/go-merkle"
)

// AllAccounts returns what you expect
func AllAccounts(tree merkle.Tree) (*AccountList, error) {
	accts, err := store.AllAccounts(tree)
	if err != nil {
		return nil, err
	}
	return RenderAccountList(accts), nil
}

// AccountByKey returns an exact match
func AccountByKey(tree merkle.Tree, key []byte) (*Account, error) {
	acct, err := store.FindAccountByKey(tree, key)
	if err != nil {
		return nil, err
	}
	if acct == nil {
		return nil, errors.New("Not Found")
	}
	return RenderAccount(acct), nil
}

// AccountByName searches for similar names
func AccountByName(tree merkle.Tree, name string) (*AccountList, error) {
	accts, err := store.SearchAccountByName(tree, name)
	if err != nil {
		return nil, err
	}
	return RenderAccountList(accts), nil
}

// PostsForAccount returns all posts that belong to this account
func PostsForAccount(tree merkle.Tree, acct []byte) (*PostList, error) {
	posts, err := store.FindPostsForAccount(tree, acct)
	if err != nil {
		return nil, err
	}
	return RenderPostList(posts), nil
}

// PostByKey returns an exact match
func PostByKey(tree merkle.Tree, key []byte) (*Post, error) {
	post, err := store.FindPostByKey(tree, key)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("Not Found")
	}
	return RenderPost(post), nil
}

// RenderQuery writes the output of any of the above issues to the client
func RenderQuery(rw http.ResponseWriter, res interface{}, err error) {
	var data []byte
	if err == nil {
		data, err = json.Marshal(res)
	}
	if err != nil {
		rw.WriteHeader(400)
		rw.Write([]byte(fmt.Sprintf("%+v", err)))
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(data)
}
