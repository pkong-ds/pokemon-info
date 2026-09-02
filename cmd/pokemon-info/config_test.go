package main

import (
	"testing"
	"time"
)

func TestParseCacheTTLDays(t *testing.T) {
	d, err := parseCacheTTL("7d")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 7*24*time.Hour {
		t.Fatalf("expected 168h, got %v", d)
	}
}

func TestParseCacheTTLGoDuration(t *testing.T) {
	d, err := parseCacheTTL("168h")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != 168*time.Hour {
		t.Fatalf("expected 168h, got %v", d)
	}
}

func TestParseCacheTTLCaseInsensitive(t *testing.T) {
	if _, err := parseCacheTTL(" 30D "); err != nil {
		t.Fatalf("expected case-insensitive parse, got %v", err)
	}
}

func TestParseCacheTTLInvalid(t *testing.T) {
	for _, s := range []string{"", "abc", "0d", "-7d", "0h", "7x"} {
		if _, err := parseCacheTTL(s); err == nil {
			t.Fatalf("expected error for %q", s)
		}
	}
}

func TestTTLStringRoundTrips(t *testing.T) {
	for _, s := range []string{"30d", "7d", "1d"} {
		d, err := parseCacheTTL(s)
		if err != nil {
			t.Fatalf("parse %q: %v", s, err)
		}
		if got := ttlString(d); got != s {
			t.Fatalf("round trip %q -> %q", s, got)
		}
	}
	if got := ttlString(90 * time.Minute); got != "1h30m0s" {
		t.Fatalf("expected Go duration string, got %q", got)
	}
}

func TestConfigFileSaveLoad(t *testing.T) {
	if _, ok := configPath(); !ok {
		t.Skip("no user config dir")
	}
	orig := loadConfig()
	t.Cleanup(func() { _ = saveConfig(orig) })

	if err := saveConfig(userConfig{CacheTTL: "14d"}); err != nil {
		t.Fatalf("save: %v", err)
	}
	cfg := loadConfig()
	if cfg.CacheTTL != "14d" {
		t.Fatalf("expected 14d, got %q", cfg.CacheTTL)
	}
}

func TestResolveCacheTTLPrecedence(t *testing.T) {
	if _, ok := configPath(); !ok {
		t.Skip("no user config dir")
	}
	orig := loadConfig()
	t.Cleanup(func() { _ = saveConfig(orig) })

	// file wins over default
	_ = saveConfig(userConfig{CacheTTL: "14d"})
	t.Setenv("POKEMON_INFO_CACHE_TTL", "")
	if got := resolveCacheTTL(); got != 14*24*time.Hour {
		t.Fatalf("expected file value 14d, got %v", got)
	}
	if src := cacheTTLSource(); src != ttlSourceFile {
		t.Fatalf("expected file source, got %q", src)
	}

	// env wins over file
	t.Setenv("POKEMON_INFO_CACHE_TTL", "48h")
	if got := resolveCacheTTL(); got != 48*time.Hour {
		t.Fatalf("expected env value 48h, got %v", got)
	}
	if src := cacheTTLSource(); src != ttlSourceEnv {
		t.Fatalf("expected env source, got %q", src)
	}
}
