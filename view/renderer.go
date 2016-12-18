package view

import (
	"encoding/hex"

	"github.com/ethanfrey/signedpost/store"
	"github.com/ethanfrey/tenderize/mom"
)

func RenderPost(post store.Post) *Post {
	// acct, err := store.AccountKeyFromPost(post.ID)
	// if err != nil {
	// 	panic(err)
	// }
	pKey, err := mom.KeyToBytes(post.Key())
	if err != nil {
		panic(err)
	}

	aKey, err := mom.KeyToBytes(post.Account)
	if err != nil {
		panic(err)
	}

	return &Post{
		ID:             hex.EncodeToString(pKey),
		AccountID:      hex.EncodeToString(aKey),
		Number:         post.Number,
		PublishedBlock: post.PublishedBlock,
		Title:          post.Title,
		Content:        post.Content,
	}
}

func RenderPostList(posts []store.Post) *PostList {
	res := PostList{
		Count: int64(len(posts)),
		Items: make([]*Post, len(posts)),
	}
	for i := range posts {
		res.Items[i] = RenderPost(posts[i])
	}
	return &res
}

func RenderAccount(acct store.Account) *Account {
	aKey, err := mom.KeyToBytes(acct.Key())
	if err != nil {
		panic(err)
	}

	return &Account{
		ID:        hex.EncodeToString(aKey),
		Name:      acct.Name,
		PostCount: acct.EntryCount,
	}
}

func RenderAccountList(accts []store.Account) *AccountList {
	res := AccountList{
		Count: int64(len(accts)),
		Items: make([]*Account, len(accts)),
	}
	for i := range accts {
		res.Items[i] = RenderAccount(accts[i])
	}
	return &res
}
