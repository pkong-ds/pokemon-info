package main

import (
	"fmt"
	"os"
	"strings"
)

const maxSuggestions = 10

type catalogEntry interface {
	entryName() string
	entrySlug() string
}

func (e PokemonEntry) entryName() string { return e.Name }
func (e PokemonEntry) entrySlug() string { return e.Slug }

func (e MoveEntry) entryName() string { return e.Name }
func (e MoveEntry) entrySlug() string { return e.Slug }

func resolveEntry[T catalogEntry](entries []T, input string) *T {
	for i := range entries {
		if strings.EqualFold(entries[i].entryName(), input) || strings.EqualFold(entries[i].entrySlug(), input) {
			return &entries[i]
		}
	}
	return nil
}

func suggestNames[T catalogEntry](entries []T, input string) []string {
	var names []string
	for _, e := range entries {
		if strings.HasPrefix(strings.ToLower(e.entryName()), strings.ToLower(input)) {
			names = append(names, e.entryName())
		}
	}
	return names
}

func resolveWithSuggestions[T catalogEntry](entries []T, input, kind string) *T {
	if e := resolveEntry(entries, input); e != nil {
		return e
	}
	suggestions := suggestNames(entries, input)
	if len(suggestions) == 0 {
		fmt.Fprintf(os.Stderr, "Error: %s '%s' not found in our list.\n", kind, input)
		return nil
	}
	if len(suggestions) == 1 {
		return resolveEntry(entries, suggestions[0])
	}
	fmt.Fprint(os.Stderr, suggestionMessage(kind, input, suggestions))
	return nil
}

func suggestionMessage(kind, input string, suggestions []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Error: %s '%s' is ambiguous. Did you mean:\n", kind, input)
	if len(suggestions) > maxSuggestions {
		for _, s := range suggestions[:maxSuggestions] {
			fmt.Fprintf(&b, "  %s\n", s)
		}
		fmt.Fprintf(&b, "  ...and %d more (press TAB to see all)\n", len(suggestions)-maxSuggestions)
		return b.String()
	}
	for _, s := range suggestions {
		fmt.Fprintf(&b, "  %s\n", s)
	}
	return b.String()
}
