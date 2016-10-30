
vendor:
	go get github.com/Masterminds/glide
	glide install

test:
	go test -p 1 `glide novendor`

build:
	go install `glide novendor`
