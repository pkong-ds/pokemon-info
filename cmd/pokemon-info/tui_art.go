package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
	"regexp"
	"strings"
	"sync"
)

// Embedded Pokémon colorscripts (ANSI truecolor half-block art), generated
// from pokemon-colorscripts (phoneybadber, MIT) plus poketex's gen 9 art
// (Caruban, see poketex PR #58). small/ and large/ sets, regular colors only.
//
//go:embed colorscripts.tar.gz
var colorscriptsArchive []byte

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*m`)

type artSet struct {
	content string
	width   int
}

var (
	artOnce  sync.Once
	artFiles map[string]artSet
)

func loadArt() {
	artOnce.Do(func() {
		artFiles = map[string]artSet{}
		gz, err := gzip.NewReader(bytes.NewReader(colorscriptsArchive))
		if err != nil {
			return
		}
		defer gz.Close()

		tr := tar.NewReader(gz)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}
			data, err := io.ReadAll(tr)
			if err != nil {
				return
			}
			artFiles[hdr.Name] = newArt(string(data))
		}
	})
}

func newArt(content string) artSet {
	lines := strings.Split(content, "\n")
	for len(lines) > 0 && strings.TrimSpace(ansiRe.ReplaceAllString(lines[len(lines)-1], "")) == "" {
		lines = lines[:len(lines)-1]
	}
	width := 0
	for _, l := range lines {
		if n := len(ansiRe.ReplaceAllString(l, "")); n > width {
			width = n
		}
	}
	return artSet{content: strings.Join(lines, "\n") + "\x1b[0m", width: width}
}

// artCandidates returns the colorscript names to try for a slug, from most
// to least specific: exact, totem forms collapsed, then progressively fewer
// trailing segments ("basculin-red-striped" → "basculin").
func artCandidates(slug string) []string {
	out := []string{slug}
	if i := strings.Index(slug, "-totem-"); i >= 0 {
		out = append(out, slug[:i]+slug[i+len("-totem-"):])
	}
	for {
		i := strings.LastIndex(slug, "-")
		if i <= 0 {
			break
		}
		slug = slug[:i]
		out = append(out, slug)
	}
	return out
}

func findArt(slug, set string) (artSet, bool) {
	loadArt()
	for _, cand := range artCandidates(slug) {
		if a, ok := artFiles[set+"/"+cand]; ok {
			return a, true
		}
	}
	return artSet{}, false
}
