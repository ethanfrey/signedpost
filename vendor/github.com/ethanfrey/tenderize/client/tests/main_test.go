package tests

import (
	"os"
	"testing"

	"github.com/tendermint/log15"
)

func TestMain(m *testing.M) {
	// turn down the logging, so we can see something else on failed tests
	h := log15.LvlFilterHandler(log15.LvlWarn, log15.StdoutHandler)
	log15.Root().SetHandler(h)
	// start a tendermint node (and dummy app) in the background to test against
	StartNode()
	os.Exit(m.Run())
}
