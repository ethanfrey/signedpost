# Signed Post

Demo app using the tendermint blockchain.

This allows people to sign posts using their public key.  In the most banal sense, this is an
immutible blogging platform on the blockchain.  And as such, serve as a demo for building an
interactive web-app based on tendermint technology.

However, it is also designed to provide legally valid proof of the existence of any document
at a given time, backed by the validation of a shared blockchain. When combined with proof of
real-world identity tied to a public-key, this could actually be a useful tool.

An old-fashioned version was to send yourself a sealed letter and use the stamp from the post office
on the sealed envelope as proof as to the date of the document.  This is a high-tech version for
signing your post.

## Quick Start

First, make sure to compile the apps

```
go install github.com/tendermint/tendermint/cmd/tendermint
go install github.com/ethanfrey/signedpost/cmd/sp-server
go install github.com/ethanfrey/signedpost/cmd/sp-cli
```

In one shell run:
```
sp-server
```

In another shell run:
```
# this is to keep us clean frm any other demos...
export TMROOT=`pwd`/tmdata
rm -rf $TMROOT

# and now start a fresh chain for this server
tendermint init
tendermint node
```

And in a final (client) shell, run:
```
# 1. make sure the server and client work...

curl -XGET localhost:54321/tndr/status
curl -XGET localhost:54321/tndr/block?height=22
sp-cli --help

# 2. create some accounts
sp-cli --key sample.key account Fred
# This gives Fred's ID
sp-cli --key sample.key account John
# This returns an error
sp-cli --key alice.key account Alice
# a new id for alice -> ALICE_ID

# 3. query these accounts
curl -XGET localhost:54321/accounts | jq
curl -XGET localhost:54321/accounts?username=Fre |jq
# use the id you get above
curl -XGET localhost:54321/accounts/$ALICE_ID

# 4. add some posts
sp-cli --key alice.key post "Hello world" "Life is good!"
sp-cli --key alice.key post "One more time" "For good luck"
# -> store post id as POST_ID

# 5. check the update
curl -XGET localhost:54321/posts/$POST_ID | jq
curl -XGET localhost:54321/accounts?username=Al | jq
curl -XGET localhost:54321/accounts/$ALICE_ID/posts | jq
```

Okay, now this worked.  But json is kinda boring...  Well, leave your tendermint app running, and open up yet another shell.

Go to github to find my example [react frontend viewer](https://github.com/ethanfrey/signedpost-react). You need npm locally, the rest of the instructions are in that repo.
(Note that only the cli can actually sign transactions, as I think giving your private key to the browser defeats the security of the blockchain)

## Data Storage

There are two types in the database.

*Account* which connects a human readable username with a public key in a first-come, first-serve basis.
This name cannot be changed.  One can only add posts to an account.

*Post* is tied to an account and leave an "immutable" (very difficult to fake) record of a document.
Any account can contain an arbitrary number of `Posts`. Each post also contains the blockheight it was
added, which can then be used to verify and timestamp it as needed.

## REST API

For querying and easy UI construction, we expose a very standard REST API.

* `GET /accounts/` returns a list of all accounts
* `GET /accounts/?username=XYZ` returns a list of all accounts containing the string `XYZ` in the username
* `GET /accounts/{id}` returns details for account with the given id
* `GET /accounts/{id}/posts/` returns a list of all posts for the given account
* `GET /accounts/{id}/posts/{pid}` returns the full details for the named post.

Notes:

* All lists may return summary information (not the full details of the structure)
* Pagination should be added to all lists by v0.2
* The objects returned are in json format and without proofs, the full-crypto version has a more complex API

Wishes:

* Some way to do pub-sub on eg. a list of posts, to be notified if an account makes a new post (tweet)

## Crypto API

However, to get the power of the blockchain, we need to expose some other features.  These require
a crypto library to sign/verify, as well as a deterministic serialization format for proper signing,
so we cannot rely on JSON.  This should be used to get true proofs. This could be implemented only in
a native app version, while the web version would allow easy browsing.

Custom data-aware endpoints:

* `GET /crypt/accounts/{id}` gets a merkle-proof of the account details (including number of posts)
* `GET /crypt/posts/{pid}` gets a merkle-proof of the post details (including block height it was added)

Proxies to tendermint core for validation:

* `POST /tndr/tx` allows one to post a new transaction to the engine (proxy to `broadcast_tx_sync`).  Post must look like `{"tx": "0123beef"}` hex-encoded form of the transaction
* `GET /tndr/block?height={h}` gets the given block
* `GET /tndr/blockchain?minHeight={min}&maxHeight={max}` gets a list of blocks
* `GET /tndr/status` gets the current blockchain status
* `GET /tndr/validators` gets the validator set

## Clients

The first client will be a web-based SPA to display the data through the REST API,
along with a golang cli tool to create accounts and add posts.

The second client will either be a more complex app (web or mobile) that uses the
crypto API to securely post data and validate and timestamp any claims.

# Dev Zone

This info is just for people wising to understand and modify this code.  Not for running the app.

## Package Layout

The top level contains the server objects, which are constructed from the pieces below.  The packages have the following responsibilities:

* txn - (de)serialization, signing, and validating signatures for all transactions that can modify state, used by client to construct, and server to verify
* store - the actual binary structures we store internally, along with functions for querying and modifying the data store
* redux - this holds the reducer, which applies `txn` Actions to the Data `store`.  The main class here is `Service`, which wraps a `go-merkle` tree
* view - these are query functions and http helpers for reading the state of the app.
* utils - common utilities (may disappear later if not really needed)
* cmd - all commands (main packages)

Top level package:

* `application.go` - Implementation of a TMSP application
* `rest.go` - Implementation of a JSON REST API to view the data
* `chain.go` - Implementation of a tendermint proxy, allowing writing transactions to the blockchain, and querying the blockchain state.

## Roadmap

### v0.1.0 (in progress)

* Implementation of users and posts in tmsp app
* Implementation of json api to view data
* Go cli to generate, sign, and submit valid transactions
* Javascript client (react) to display state via json api

### v0.2.0 (planned)

* Android native app capable of submitting transactions as well as viewing state
* Websockets to "watch" a query for live updates (auto-update on post by a given user)
* Integrate live update in both clients

### Future ideas

* Protocol buffer instead of go-wire and portable crypto libraries
  * Easier for signing/verifying with no-golang clients
* Validating "light-client" features
  * Get and verify proof for a given state
  * Verify block headers
  * Display and validate block where post was submitted (audit trail)
* iOS app?
* Desktop app?
* Multiple content-types of post
  * Encoded text, timestamp validated, key may be released much later to publish
  * Special data fields with app-specific meaning?

## Licensing

Please note that all code here is currently under the GPLv3.

It is intended as an example app, and a free place to work on best practices without commercial interest. All contributions (and public forks) are welcome in order to evolve best practices in tendermint apps.  However, if you wish to integrate any of this code in a commercial application, please check with the author first.
