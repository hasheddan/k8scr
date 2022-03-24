# Set the shell to bash always
SHELL := /bin/bash

build: test
	@CGO_ENABLED=0 go build -o ./build/k8scr ./cmd/k8scr
	@CGO_ENABLED=0 go build -o ./build/distribution ./cmd/k8scr-distribution

image:
	@docker build . -f distribution.Dockerfile -t hasheddan/k8scr-distribution:latest --load

all: build image

tidy:
	@go mod tidy

test:
	@go test -v ./...

.PHONY: tidy clean build image all
