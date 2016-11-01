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
