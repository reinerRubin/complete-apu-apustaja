package cache

import (
	"context"
	"testing"
	"time"
)

// TODO: make this tests generic to cache interface
func TestInMemoryCacheBasic(t *testing.T) {
	cache := NewInMemoryCache(10 * time.Second)
	cache.Start()
	defer cache.Stop()

	key := Key("key1")
	body := []byte("suchBody")

	ctx := context.Background()
	err := cache.Set(ctx, key, body, 2*time.Second)
	if err != nil {
		t.Fatalf("cant set key to cache: %s", err)
	}

	probablyBody, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("cant get from cache: %s", err)
	}

	if string(probablyBody) != string(body) {
		t.Fatalf("body does not match. expected: %s, actual: %s",
			string(body), string(probablyBody))
	}
}

func TestInMemoryExpirationInClearPeriod(t *testing.T) {
	cache := NewInMemoryCache(200 * time.Millisecond)
	cache.Start()
	defer cache.Stop()

	key := Key("key1")
	body := []byte("suchBody")

	ctx := context.Background()
	err := cache.Set(ctx, key, body, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("cant set key to cache: %s", err)
	}

	<-time.After(201 * time.Millisecond)

	probablyBody, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("cant get from cache: %s", err)
	}

	if probablyBody != nil {
		t.Fatalf("body must be nil but %s", probablyBody)
	}
}

func TestInMemoryExpirationWithoutClearPeriod(t *testing.T) {
	cache := NewInMemoryCache(10 * time.Second)
	cache.Start()
	defer cache.Stop()

	key := Key("key1")
	body := []byte("suchBody")

	ctx := context.Background()
	err := cache.Set(ctx, key, body, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("cant set key to cache: %s", err)
	}

	<-time.After(11 * time.Millisecond)

	probablyBody, err := cache.Get(ctx, key)
	if err != nil {
		t.Fatalf("cant get from cache: %s", err)
	}

	if probablyBody != nil {
		t.Fatalf("body must be nil but %s", probablyBody)
	}
}
