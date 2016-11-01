#! /bin/bash

# set glide.lock path
if [[ "$GLIDE" == "" ]]; then
	GLIDE=$GOPATH/src/github.com/tendermint/tendermint/glide.lock
fi

# get vendored commit for given lib

####################
# libs we depend on
####################

LIBS_GO_TEST=(go-clist go-common go-config go-crypto go-db go-events go-merkle go-p2p)
LIBS_MAKE_TEST=(go-rpc go-wire tmsp)

for lib in "${LIBS_GO_TEST[@]}"; do

	# checkout vendored version of lib
	bash scripts/glide/checkout.sh $GLIDE $lib

	echo "Testing $lib ..."
	go test --race github.com/tendermint/$lib/...
	if [[ "$?" != 0 ]]; then
		echo "FAIL"
		exit 1
	fi
done

for lib in "${LIBS_MAKE_TEST[@]}"; do
	getDep $lib

	echo "Testing $lib ..."
	cd $GOPATH/src/github.com/tendermint/$lib
	make test
	if [[ "$?" != 0 ]]; then
		echo "FAIL"
		exit 1
	fi
done

echo ""
echo "PASS"
