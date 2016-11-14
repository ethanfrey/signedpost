package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/tendermint/go-merkle"
	"github.com/tendermint/tmsp/server"

	"github.com/ethanfrey/signedpost"
)

// MakeServer creates an http server
func MakeServer(listen string, app *signedpost.Application, proxy signedpost.Proxy) *http.Server {
	r := mux.NewRouter()
	app.AddQueryRoutes(r)
	proxy.AddChainRoutes(r)
	wrap := cors.Default().Handler(r)

	s := &http.Server{
		Addr:    listen,
		Handler: wrap,
	}
	return s

}

func main() {
	tmspPtr := flag.String("addr", "tcp://0.0.0.0:46658", "Address for tmsp server to listen")
	protoPtr := flag.String("tmsp", "socket", "socket | grpc")
	rpcPtr := flag.String("rpc", "localhost:46657", "Address of tendermint core rpc server")
	servePtr := flag.String("http", ":54321", "Port to serve the custom http application")
	flag.Parse()

	// these should come from command-line
	tree := merkle.NewIAVLTree(0, nil)

	app := signedpost.NewApp(tree)
	proxy := signedpost.NewProxy(*rpcPtr)

	// start tmsp server
	_, err := server.NewServer(*tmspPtr, *protoPtr, app)
	if err != nil {
		fmt.Printf("TMSP server failed: %+v\n", err)
		return
	}

	// start http server
	srv := MakeServer(*servePtr, app, proxy)
	fmt.Println("Starting http server on port", *servePtr)
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Printf("HTTP server failed: %+v\n", err)
		return
	}
	fmt.Println("Finished")
}
