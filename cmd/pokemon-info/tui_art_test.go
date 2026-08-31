package main

import (
	"strings"
	"testing"
)

func TestArtLookup(t *testing.T) {
	a, ok := findArt("bulbasaur", "small")
	if !ok {
		t.Fatal("no art found for bulbasaur")
	}
	if !strings.Contains(a.content, "▀") && !strings.Contains(a.content, "▄") {
		t.Fatalf("content has no halfblocks: %q", a.content[:100])
	}
	t.Logf("bulbasaur small: %d wide, %d lines", a.width, strings.Count(a.content, "\n")+1)

	if _, ok := findArt("bulbasaur", "large"); !ok {
		t.Error("no large art for bulbasaur")
	}

	for _, slug := range []string{
		"basculin-red-striped",
		"maushold-family-of-four",
		"raticate-totem-alola",
		"venusaur-mega",
		"deoxys-normal",
	} {
		if _, ok := findArt(slug, "small"); !ok {
			t.Errorf("fallback miss for %s", slug)
		}
	}

	if _, ok := findArt("does-not-exist", "small"); ok {
		t.Error("expected miss for unknown slug")
	}
}
