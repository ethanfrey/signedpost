package view

import (
	"encoding/hex"

	"github.com/ethanfrey/signedpost/store"
)

func RenderPost(post *store.PostField) *Post {
	acct, err := store.AccountKeyFromPost(post.Key)
	if err != nil {
		panic(err)
	}
	return &Post{
		ID:             hex.EncodeToString(post.Key),
		AccountID:      hex.EncodeToString(acct),
		Number:         post.Number,
		PublishedBlock: post.PublishedBlock,
		Title:          post.Title,
		Content:        post.Content,
	}
}

func RenderPostList(posts []*store.PostField) *PostList {
	res := PostList{
		Count: int64(len(posts)),
		Items: make([]*Post, len(posts)),
	}
	for i := range posts {
		res.Items[i] = RenderPost(posts[i])
	}
	return &res
}

func RenderAccount(acct *store.AccountField) *Account {
	return &Account{
		ID:        hex.EncodeToString(acct.Key),
		Name:      acct.Name,
		PostCount: acct.EntryCount,
	}
}

func RenderAccountList(accts []*store.AccountField) *AccountList {
	res := AccountList{
		Count: int64(len(accts)),
		Items: make([]*Account, len(accts)),
	}
	for i := range accts {
		res.Items[i] = RenderAccount(accts[i])
	}
	return &res
}
