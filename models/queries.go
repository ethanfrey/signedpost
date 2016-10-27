package models

import "github.com/ethanfrey/signedpost/db"

// AccountQuery finds an account by pk or name
type AccountQuery struct {
	PK   []byte `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// AccountResponse is the return value from an AccountQuery
type AccountResponse struct {
	*db.Account // TODO: modify later... eg. proof
}

// EntryQuery finds an entry by account pk and (optional) number
type EntryQuery struct {
	PK     []byte `json:"account"`
	Number int    `json:"number,omitempty"` // if not present, then the latest one
}

// EntryResponse is the return value from an EntryQuery
type EntryResponse struct {
	db.Entry // TODO: modify later... eg. proof
}

type NotFoundResponse struct {
}
