VERSION?=$(shell git describe --tags --always --dirty)

all: build

build: brightness
	go build -ldflags "-X main.version=$(VERSION)"

install: build
	go install -ldflags "-X main.version=$(VERSION)"

dist: dist/necd_darwin_amd64 dist/necd_linux_amd64

dist/necd_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o dist/necd_darwin_amd64 -ldflags "-X main.version=$(VERSION)"

dist/necd_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o dist/necd_linux_amd64 -ldflags "-X main.version=$(VERSION)"

rel: dist
	hub release create -a dist $(VERSION)


brightness:
	gcc -std=c99 -o brightness c/brightness.c -framework IOKit -framework ApplicationServices
