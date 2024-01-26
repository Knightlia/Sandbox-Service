COMMIT   := $(shell git rev-list --tags --max-count 1)
VERSION  := $(shell git describe --tags $(COMMIT))
LD_FLAGS := -X main.version=$(VERSION)

build:
	go build -ldflags="$(LD_FLAGS) -s -w"

run: build
	./sandbox-service
