package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	cacheSubdir = "pokemon-info"
	userAgent   = "pokemon-info"

	// Token bucket: the burst covers one full detail chain (3 requests)
	// so a single lookup pays no artificial delay; the refill rate is the
	// sustained request ceiling when rapidly browsing.
	rateBurst      = 3.0
	rateRefillPerS = 5.0
)

// cacheTTL is resolved at startup: env var > config file > default
// (see config.go). Reassigned live by the /config page on save.
var cacheTTL = resolveCacheTTL()

var (
	cacheDir string
)

func init() {
	dir, err := os.UserCacheDir()
	if err != nil {
		return
	}
	cacheDir = filepath.Join(dir, cacheSubdir)
}

func cacheKey(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Path == "" || u.Path == "/" {
		return hashURL(rawURL)
	}
	path := strings.Trim(u.Path, "/")
	path = strings.TrimPrefix(path, "api/v2/")
	path = strings.ReplaceAll(path, "/", "-")
	if path == "" {
		return hashURL(rawURL)
	}
	return path + ".json"
}

func hashURL(rawURL string) string {
	h := sha256.Sum256([]byte(rawURL))
	return hex.EncodeToString(h[:])[:16] + ".json"
}

func cachedGet(rawURL string) ([]byte, error) {
	if cacheDir != "" {
		p := filepath.Join(cacheDir, cacheKey(rawURL))
		if info, err := os.Stat(p); err == nil && time.Since(info.ModTime()) < cacheTTL {
			return os.ReadFile(p)
		}
	}

	body, err := rateLimitedGet(rawURL)
	if err != nil {
		return nil, err
	}

	if cacheDir != "" {
		p := filepath.Join(cacheDir, cacheKey(rawURL))
		if os.MkdirAll(cacheDir, 0755) == nil {
			_ = os.WriteFile(p, body, 0644)
		}
	}
	return body, nil
}

// clearDiskCache removes every cached response. Best-effort: errors are
// ignored because the cache is disposable.
func clearDiskCache() {
	if cacheDir == "" {
		return
	}
	_ = os.RemoveAll(cacheDir)
	_ = os.MkdirAll(cacheDir, 0755)
}

// tokenBucket bounds the request rate: burst requests pass immediately,
// then requests space out at the refill rate. Waiters sleep outside the
// lock; advancing `last` makes concurrent waiters queue one interval apart.
type tokenBucket struct {
	mu     sync.Mutex
	tokens float64
	last   time.Time
}

var apiBucket = &tokenBucket{tokens: rateBurst, last: time.Now()}

func (tb *tokenBucket) take() {
	tb.mu.Lock()
	now := time.Now()
	tb.tokens += now.Sub(tb.last).Seconds() * rateRefillPerS
	if tb.tokens > rateBurst {
		tb.tokens = rateBurst
	}
	tb.last = now
	if tb.tokens >= 1 {
		tb.tokens--
		tb.mu.Unlock()
		return
	}
	deficit := 1 - tb.tokens
	tb.tokens = 0
	wait := time.Duration(deficit / rateRefillPerS * float64(time.Second))
	tb.last = now.Add(wait)
	tb.mu.Unlock()
	time.Sleep(wait)
}

func rateLimitedGet(rawURL string) ([]byte, error) {
	apiBucket.take()

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request to %s failed with status %s: %s", rawURL, resp.Status, string(body))
	}

	return io.ReadAll(resp.Body)
}
