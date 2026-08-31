.PHONY: build
build:
	go build -o ../dist/prepare

.PHONY: clean
clean:
	rm -rf ../dist

.PHONY: run
run:
	../dist/prepare --output-format yaml -o ../cmd/pokemons.yaml
