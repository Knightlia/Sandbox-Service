COMMIT := $(shell git rev-parse HEAD)
VERSION := $(shell git describe --tags $(COMMIT) 2> /dev/null || echo $(COMMIT))
COMMIT := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date +%FT%T%z)
LD_FLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)

build:
	go build -ldflags="$(LD_FLAGS)"

run:
	go run -ldflags="$(LD_FLAGS)" main.go --debug

test:
	go test ./...

coverage:
	go test -cover -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LD_FLAGS)"
