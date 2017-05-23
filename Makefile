.PHONY: build clean docker vendor

GOPATH := ${PWD}/vendor:${GOPATH}
export GOPATH

default: build

build:
	GOOS=linux go build -o bin/cookiemonster-linux-amd64 -v ./src/*.go
	GOOS=darwin go build -o bin/cookiemonster-darwin-amd64 -v ./src/*.go

clean:
	rm -rf cookiemonster ./bin

docker:
	docker build --no-cache -t oreo01:5000/cookiemonster .
	docker push oreo01:5000/cookiemonster:latest

vendor:
	rm -rf ./vendor
	GOPATH=${PWD}/vendor go get github.com/gorilla/mux
