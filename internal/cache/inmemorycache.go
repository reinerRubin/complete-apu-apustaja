package cache

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// I would not do it in a production code. Redis with TTL would be a good solution.
// But it was kind of fun to implement
type (
	InMemoryCache struct {
		cache sync.Map

		stopOnce       sync.Once
		stopChannel    chan struct{}
		stoppedChannel chan struct{}

		clearCyclePeriod time.Duration
	}

	imMemoryCacheRecord struct {
		body      []byte
		expiredAt time.Time
	}
)

func NewInMemoryCache(clearCyclePeriod time.Duration) *InMemoryCache {
	return &InMemoryCache{
		stopChannel:    make(chan struct{}),
		stoppedChannel: make(chan struct{}),

		clearCyclePeriod: clearCyclePeriod,
	}
}

func (c *InMemoryCache) Get(ctx context.Context, key Key) ([]byte, error) {
	body, found := c.cache.Load(key)
	if !found {
		return nil, nil
	}

	cache, ok := body.(*imMemoryCacheRecord)
	if !ok {
		return nil, fmt.Errorf("cant extract data")
	}

	if cache.isExpired() {
		return nil, nil
	}

	return cache.body, nil
}

func (c *InMemoryCache) Set(ctx context.Context, key Key, body []byte, ttl time.Duration) error {
	c.cache.Store(key, &imMemoryCacheRecord{
		body:      body,
		expiredAt: time.Now().Add(ttl),
	})

	return nil
}

func (c *InMemoryCache) Start() {
	go c.run()
}

// see InMemoryCache description
func (c *InMemoryCache) run() {
	ticker := time.NewTicker(c.clearCyclePeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cache.Range(func(key, value interface{}) bool {
				cache, ok := value.(*imMemoryCacheRecord)
				if !ok {
					// idk
					log.Printf("cache contains some not cacheRecord records")
					c.cache.Delete(key)
				}

				if cache.isExpired() {
					c.cache.Delete(key)
				}

				return true
			})
		case <-c.stopChannel:
			close(c.stoppedChannel)
			return
		}
	}
}

func (c *InMemoryCache) Stop() error {
	c.stopOnce.Do(func() {
		close(c.stopChannel)
	})
	<-c.stoppedChannel

	return nil
}

func (cr *imMemoryCacheRecord) isExpired() bool {
	return cr.expiredAt.Before(time.Now())
}
