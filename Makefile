.PHONY: build clean docker vendor

default: build-local

all: vendor build docker

build: build-darwin build-linux

build-local:
	gofmt -w ./src/*.go
	gofmt -w ./src/worker/*.go
	go build -o bin/cookiemonster -v ./src/*.go
	go build -o bin/cookiemonsterworker -v ./src/worker/*.go

build-darwin:
	GOOS=darwin CGO_ENABLED=0 go build -o bin/cookiemonster-darwin-amd64 -v ./src/*.go
	GOOS=darwin CGO_ENABLED=0 go build -o bin/cookiemonsterworker-darwin-amd64 -v ./src/worker/*.go

build-linux:
	GOOS=linux CGO_ENABLED=0 go build -o bin/cookiemonster-linux-amd64 -v ./src/*.go
	GOOS=linux CGO_ENABLED=0 go build -o bin/cookiemonsterworker-linux-amd64 -v ./src/worker/*.go

clean:
	rm -rf ./bin ./vendor

docker:
	docker build --no-cache -t oreo01:5000/cookiemonster -f Dockerfile.cookiemonster .
	docker build --no-cache -t oreo01:5000/cookiemonsterworker -f Dockerfile.cookiemonsterworker .
	docker push oreo01:5000/cookiemonster:latest
	docker push oreo01:5000/cookiemonsterworker:latest

vendor:
	glide install
