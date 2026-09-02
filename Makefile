.PHONY: build test catalog demos clean

VERSION ?= $(shell git describe --tags --always --dirty)

build:
	go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o dist/pokemon-info ./cmd/pokemon-info

test:
	go vet ./... && go test ./...

catalog:
	go run ./cmd/prepare --output-format yaml --resource pokemon --output-file cmd/pokemon-info/pokemons.yaml
	go run ./cmd/prepare --output-format yaml --resource move --output-file cmd/pokemon-info/moves.yaml

demos:
	vhs docs/tapes/demo.tape
	vhs docs/tapes/moves.tape

clean:
	rm -rf dist
