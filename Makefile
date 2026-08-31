.PHONY: build test catalog clean

build:
	go build -o dist/pokemon-info ./cmd/pokemon-info

test:
	go vet ./... && go test ./...

catalog:
	go run ./cmd/prepare --output-format yaml --resource pokemon --output-file cmd/pokemon-info/pokemons.yaml
	go run ./cmd/prepare --output-format yaml --resource move --output-file cmd/pokemon-info/moves.yaml

clean:
	rm -rf dist
