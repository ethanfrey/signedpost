package controllers

import (
	"github.com/ethanfrey/signedpost/db"
	"github.com/ethanfrey/signedpost/models"
)

// FindAccount returns the account that matches, or nil on nothing
func FindAccount(ctx *Context, query *models.AccountQuery) *models.AccountResponse {
	var match *db.Account
	if query.PK != nil && len(query.PK) > 0 {
		match, _ = db.FindAccountByPK(ctx.GetDB(), query.PK)
	} else if query.Name != "" {
		match, _ = db.FindAccountByName(ctx.GetDB(), query.Name)
	}

	if match == nil {
		return nil
	}
	// TODO: add proof
	return &models.AccountResponse{match}
}
