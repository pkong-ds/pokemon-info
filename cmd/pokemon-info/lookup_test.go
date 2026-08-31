package main

import (
	"strings"
	"testing"
)

func TestResolveEntry(t *testing.T) {
	entries := []PokemonEntry{
		{ID: 25, Slug: "pikachu", Name: "Pikachu"},
		{ID: 122, Slug: "mr-mime", Name: "Mr. Mime"},
	}
	tests := []struct {
		input string
		want  string
	}{
		{"Pikachu", "Pikachu"},
		{"PIKACHU", "Pikachu"},
		{"pikachu", "Pikachu"},
		{"mr-mime", "Mr. Mime"},
		{"MR-MIME", "Mr. Mime"},
		{"Mr. Mime", "Mr. Mime"},
		{"pika", ""},
		{"mrmime", ""},
		{"raichu", ""},
	}
	for _, tt := range tests {
		got := resolveEntry(entries, tt.input)
		if tt.want == "" {
			if got != nil {
				t.Errorf("resolveEntry(%q) = %v, want nil", tt.input, got.Name)
			}
			continue
		}
		if got == nil || got.Name != tt.want {
			t.Errorf("resolveEntry(%q) = %v, want %s", tt.input, got, tt.want)
		}
	}
}

func TestResolveEntryMoves(t *testing.T) {
	entries := []MoveEntry{
		{ID: 85, Slug: "thunderbolt", Name: "Thunderbolt"},
	}
	if got := resolveEntry(entries, "Thunderbolt"); got == nil {
		t.Error("resolveEntry(moves, \"Thunderbolt\") = nil, want entry")
	}
	if got := resolveEntry(entries, "surf"); got != nil {
		t.Errorf("resolveEntry(moves, \"surf\") = %v, want nil", got.Name)
	}
}

func TestSuggestNames(t *testing.T) {
	entries := []PokemonEntry{
		{Slug: "pikachu", Name: "Pikachu"},
		{Slug: "pidgey", Name: "Pidgey"},
		{Slug: "pidgeotto", Name: "Pidgeotto"},
		{Slug: "mr-mime", Name: "Mr. Mime"},
	}
	tests := []struct {
		input string
		want  []string
	}{
		{"pi", []string{"Pikachu", "Pidgey", "Pidgeotto"}},
		{"PI", []string{"Pikachu", "Pidgey", "Pidgeotto"}},
		{"pid", []string{"Pidgey", "Pidgeotto"}},
		{"mr", []string{"Mr. Mime"}},
		{"xy", nil},
		{"", []string{"Pikachu", "Pidgey", "Pidgeotto", "Mr. Mime"}},
	}
	for _, tt := range tests {
		got := suggestNames(entries, tt.input)
		if strings.Join(got, ",") != strings.Join(tt.want, ",") {
			t.Errorf("suggestNames(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestResolveWithSuggestions(t *testing.T) {
	entries := []PokemonEntry{
		{Slug: "pikachu", Name: "Pikachu"},
		{Slug: "pidgey", Name: "Pidgey"},
		{Slug: "pidgeotto", Name: "Pidgeotto"},
	}
	if got := resolveWithSuggestions(entries, "pika", "Pokémon"); got == nil || got.Name != "Pikachu" {
		t.Errorf("resolveWithSuggestions(%q) = %v, want Pikachu (unique prefix resolves)", "pika", got)
	}
	if got := resolveWithSuggestions(entries, "pid", "Pokémon"); got != nil {
		t.Error("resolveWithSuggestions(ambiguous) = non-nil, want nil")
	}
	if got := resolveWithSuggestions(entries, "zzz", "Pokémon"); got != nil {
		t.Error("resolveWithSuggestions(no match) = non-nil, want nil")
	}
}

func TestSuggestionMessage(t *testing.T) {
	suggestions := make([]string, 15)
	for i := range suggestions {
		suggestions[i] = strings.Repeat("x", i+1)
	}
	msg := suggestionMessage("Pokémon", "char", suggestions)
	if !strings.Contains(msg, "is ambiguous") {
		t.Errorf("message missing ambiguity header: %q", msg)
	}
	if !strings.Contains(msg, suggestions[maxSuggestions-1]) {
		t.Error("message missing last shown suggestion")
	}
	if strings.Contains(msg, suggestions[maxSuggestions]) {
		t.Error("message shows suggestion beyond cap")
	}
	if !strings.Contains(msg, "and 5 more") {
		t.Errorf("message missing overflow count: %q", msg)
	}
	short := suggestionMessage("Move", "x", []string{"A", "B"})
	if !strings.Contains(short, "A") || !strings.Contains(short, "B") {
		t.Errorf("short message missing candidates: %q", short)
	}
	if strings.Contains(short, "more") {
		t.Errorf("short message has overflow line: %q", short)
	}
}

func TestEmbeddedCatalog(t *testing.T) {
	if len(allPokemonEntries) < 1000 {
		t.Errorf("embedded Pokémon catalog too small: %d entries", len(allPokemonEntries))
	}
	if resolveEntry(allPokemonEntries, "Pikachu") == nil {
		t.Error("Pikachu missing from embedded catalog")
	}
	if resolveEntry(allPokemonEntries, "mr-mime") == nil {
		t.Error("slug mr-mime missing from embedded catalog")
	}
	if len(allMoveEntries) < 500 {
		t.Errorf("embedded move catalog too small: %d entries", len(allMoveEntries))
	}
	if resolveEntry(allMoveEntries, "Thunderbolt") == nil {
		t.Error("Thunderbolt missing from embedded catalog")
	}
}
