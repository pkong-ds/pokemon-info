.PHONY: build
build:
	go build -o ../dist/pokemon_info

.PHONY: clean
clean:
	rm -rf ../dist


.PHONY: autocomplete-zsh
autocomplete-zsh:
	../dist/pokemon_info completion zsh > ~/.zsh/completion/_pokemon_info
# have below in .zshrc

# fpath=(~/.zsh/completion $fpath)
# autoload -U compinit && compinit


