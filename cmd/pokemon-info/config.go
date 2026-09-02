package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const defaultCacheTTL = 30 * 24 * time.Hour

// userConfig is persisted at os.UserConfigDir()/pokemon-info/config.json.
// Values are duration strings; "Nd" day shorthand is accepted.
type userConfig struct {
	CacheTTL string `json:"cache_ttl,omitempty"`
}

func configPath() (string, bool) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", false
	}
	return filepath.Join(dir, "pokemon-info", "config.json"), true
}

// loadConfig is best-effort: unreadable or malformed files yield the
// zero config, never an error.
func loadConfig() userConfig {
	var cfg userConfig
	p, ok := configPath()
	if !ok {
		return cfg
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, &cfg)
	return cfg
}

func saveConfig(cfg userConfig) error {
	p, ok := configPath()
	if !ok {
		return fmt.Errorf("no user config directory")
	}
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}

// cacheTTL sources, most to least specific.
const (
	ttlSourceEnv     = "env POKEMON_INFO_CACHE_TTL"
	ttlSourceFile    = "config file"
	ttlSourceDefault = "default"
)

// resolveCacheTTL applies precedence: env var > config file > default.
func resolveCacheTTL() time.Duration {
	if v := os.Getenv("POKEMON_INFO_CACHE_TTL"); v != "" {
		if d, err := parseCacheTTL(v); err == nil {
			return d
		}
	}
	if d, err := parseCacheTTL(loadConfig().CacheTTL); err == nil {
		return d
	}
	return defaultCacheTTL
}

// cacheTTLSource reports where the effective TTL comes from, for the
// /config page.
func cacheTTLSource() string {
	if v := os.Getenv("POKEMON_INFO_CACHE_TTL"); v != "" {
		if _, err := parseCacheTTL(v); err == nil {
			return ttlSourceEnv
		}
	}
	if _, err := parseCacheTTL(loadConfig().CacheTTL); err == nil {
		return ttlSourceFile
	}
	return ttlSourceDefault
}

// parseCacheTTL accepts Go durations ("168h") and day shorthand ("7d").
func parseCacheTTL(s string) (time.Duration, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0, fmt.Errorf("empty value")
	}
	if strings.HasSuffix(s, "d") {
		n, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		if err != nil || n <= 0 {
			return 0, fmt.Errorf("invalid day count")
		}
		return time.Duration(n) * 24 * time.Hour, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("invalid duration")
	}
	return d, nil
}

// ttlString renders a duration back in the friendliest form: whole days
// as "Nd", anything else as a Go duration string.
func ttlString(d time.Duration) string {
	day := 24 * time.Hour
	if d%day == 0 {
		return fmt.Sprintf("%dd", int(d/day))
	}
	return d.String()
}
