NAME=vsv-decoder
VERSION=$(shell git tag -l --points-at HEAD)
VERSION:=$(if $(VERSION),$(VERSION),"DEV")
BUILD=$(shell git rev-parse --short HEAD)
DATE=$(shell date -u +%Y.%m.%d_%H:%M:%S)
LD_FLAGS="-w -X main.execName=$(NAME) -X main.version=$(VERSION) -X main.build=$(BUILD) -X main.date=$(DATE)"

.PHONY: all build clean release run

all:
	GOOS=darwin  GOARCH=amd64 go build -tags release -ldflags $(LD_FLAGS) -o build/${NAME}-$(VERSION)-darwin-amd64      ./cmd/encoder
	GOOS=darwin  GOARCH=arm64 go build -tags release -ldflags $(LD_FLAGS) -o build/${NAME}-$(VERSION)-darwin-arm64      ./cmd/encoder
	GOOS=linux   GOARCH=amd64 go build -tags release -ldflags $(LD_FLAGS) -o build/${NAME}-$(VERSION)-linux-amd64       ./cmd/encoder
	GOOS=linux   GOARCH=arm64 go build -tags release -ldflags $(LD_FLAGS) -o build/${NAME}-$(VERSION)-linux-arm64       ./cmd/encoder
	GOOS=windows GOARCH=amd64 go build -tags release -ldflags $(LD_FLAGS) -o build/${NAME}-$(VERSION)-windows-amd64.exe ./cmd/encoder
	GOOS=windows GOARCH=arm64 go build -tags release -ldflags $(LD_FLAGS) -o build/${NAME}-$(VERSION)-windows-arm64.exe ./cmd/encoder
	cd build; sha256sum * > sha256sums.txt

build:
	go build -tags release -ldflags $(LD_FLAGS) -o build/${NAME} ./cmd/encoder

clean:
	rm -rf build/*

release:
	cd build; gh release create $(VERSION) -d -t v$(VERSION) *

run:
	go run -tags release -ldflags $(LD_FLAGS) ./cmd/encoder $(ARGS)
