package controllers

import (
	"github.com/ethanfrey/bloggermint/db"
	"github.com/ethanfrey/bloggermint/models"
)

// FindAccount returns the account that matches, or nil on nothing
func FindAccount(ctx *Context, query *models.AccountQuery) *models.AccountResponse {
	var match *db.Account
	if query.PK != nil && len(query.PK) > 0 {
		match = db.FindAccountByPK(ctx.GetDB(), query.PK)
	} else if query.Name != "" {
		match = db.FindAccountByName(ctx.GetDB(), query.Name)
	}

	if match == nil {
		return nil
	}
	// TODO: add proof
	return &models.AccountResponse{match}
}
