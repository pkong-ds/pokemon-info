.PHONY: build test catalog demos clean

build:
	go build -o dist/pokemon-info ./cmd/pokemon-info

test:
	go vet ./... && go test ./...

catalog:
	go run ./cmd/prepare --output-format yaml --resource pokemon --output-file cmd/pokemon-info/pokemons.yaml
	go run ./cmd/prepare --output-format yaml --resource move --output-file cmd/pokemon-info/moves.yaml

demos:
	vhs doc/tapes/demo.tape
	vhs doc/tapes/moves.tape
	vhs doc/tapes/cli.tape

clean:
	rm -rf dist
