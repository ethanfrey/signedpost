package signedpost

import (
	"encoding/hex"
	"net/http"

	"github.com/ethanfrey/signedpost/utils"
	"github.com/ethanfrey/signedpost/view"
	"github.com/gorilla/mux"
)

/*
This file contains all REST API calls exposed by the app (simple json queries)
*/

func (app *Application) SearchAccounts(rw http.ResponseWriter, r *http.Request) {
	var accts *view.AccountList
	var err error
	name := r.URL.Query().Get("username")
	if name == "" {
		accts, err = view.AllAccounts(app.commited.GetDB())
	} else {
		accts, err = view.AccountByName(app.commited.GetDB(), name)
	}
	utils.RenderQuery(rw, accts, err)
}

func (app *Application) AccountByKey(rw http.ResponseWriter, r *http.Request) {
	var acct *view.Account
	q := mux.Vars(r)["acct"]
	key, err := hex.DecodeString(q)
	if err == nil {
		acct, err = view.AccountByKey(app.commited.GetDB(), key)
	}
	utils.RenderQuery(rw, acct, err)
}

func (app *Application) PostByKey(rw http.ResponseWriter, r *http.Request) {
	var post *view.Post
	q := mux.Vars(r)["post"]
	key, err := hex.DecodeString(q)
	if err == nil {
		post, err = view.PostByKey(app.commited.GetDB(), key)
	}
	utils.RenderQuery(rw, post, err)
}

func (app *Application) PostsForAccount(rw http.ResponseWriter, r *http.Request) {
	var posts *view.PostList
	q := mux.Vars(r)["acct"]
	key, err := hex.DecodeString(q)
	if err == nil {
		posts, err = view.PostsForAccount(app.commited.GetDB(), key)
	}
	utils.RenderQuery(rw, posts, err)
}

// AddQueryRoutes add all routes for reading the app state (unsigned)
func (app *Application) AddQueryRoutes(r *mux.Router) {
	r.HandleFunc("/accounts", app.SearchAccounts).Methods("GET")
	r.HandleFunc("/accounts/{acct}", app.AccountByKey).Methods("GET")
	r.HandleFunc("/accounts/{acct}/posts", app.PostsForAccount).Methods("GET")
	r.HandleFunc("/posts/{post}", app.PostByKey).Methods("GET")
}
