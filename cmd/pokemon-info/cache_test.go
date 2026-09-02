package main

import (
	"testing"
	"time"
)

func TestLRUCacheEvictsPastCap(t *testing.T) {
	c := newLRUCache(3)
	c.set("a", 1)
	c.set("b", 2)
	c.set("c", 3)
	c.set("d", 4) // evicts "a"

	if c.has("a") {
		t.Fatal("expected oldest key to be evicted")
	}
	if v, ok := c.get("d"); !ok || v.(int) != 4 {
		t.Fatalf("expected d=4, got %v %v", v, ok)
	}
}

func TestLRUCacheGetPromotesRecency(t *testing.T) {
	c := newLRUCache(3)
	c.set("a", 1)
	c.set("b", 2)
	c.set("c", 3)

	if _, ok := c.get("a"); !ok { // "a" is now most recent
		t.Fatal("expected a to be present")
	}
	c.set("d", 4) // evicts "b" (least recent), not "a"

	if !c.has("a") {
		t.Fatal("expected promoted key to survive eviction")
	}
	if c.has("b") {
		t.Fatal("expected least-recently-used key to be evicted")
	}
}

func TestLRUCacheOverwriteKeepsSingleEntry(t *testing.T) {
	c := newLRUCache(3)
	c.set("a", 1)
	c.set("a", 2)
	if c.ll.Len() != 1 {
		t.Fatalf("expected 1 entry after overwrite, got %d", c.ll.Len())
	}
	if v, _ := c.get("a"); v.(int) != 2 {
		t.Fatalf("expected overwritten value 2, got %v", v)
	}
}

func TestLRUCacheClear(t *testing.T) {
	c := newLRUCache(3)
	c.set("a", 1)
	c.clear()
	if c.has("a") || c.ll.Len() != 0 {
		t.Fatal("expected empty cache after clear")
	}
}

func TestTokenBucketBurstAllowsChain(t *testing.T) {
	tb := &tokenBucket{tokens: rateBurst, last: time.Now()}
	start := time.Now()
	for i := 0; i < 3; i++ {
		tb.take()
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Fatalf("burst of %d should pass without delay, took %v", int(rateBurst), elapsed)
	}
}

func TestTokenBucketThrottlesPastBurst(t *testing.T) {
	tb := &tokenBucket{tokens: rateBurst, last: time.Now()}
	for i := 0; i < 3; i++ {
		tb.take()
	}
	start := time.Now()
	tb.take() // 4th request: bucket empty
	minWait := time.Duration((1 / rateRefillPerS) * float64(time.Second) * 0.8)
	if elapsed := time.Since(start); elapsed < minWait {
		t.Fatalf("expected ~%v wait past burst, took %v", minWait, elapsed)
	}
}
