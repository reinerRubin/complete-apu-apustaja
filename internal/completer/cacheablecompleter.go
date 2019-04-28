package completer

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/reinerRubin/complete-apu-apustaja/internal/cache"
)

type CacheableCompleter struct {
	completer Completer
	cache     cache.Cache
	ttl       time.Duration
}

// TODO: move to config
func NewCacheableCompleter(
	completer Completer,
	cache cache.Cache,
	ttl time.Duration,
) *CacheableCompleter {
	return &CacheableCompleter{
		completer: completer,
		cache:     cache,
	}
}

func (c *CacheableCompleter) Complete(ctx context.Context, q *Query) (*Suggestions, error) {
	suggestions, err := c.get(ctx, q)
	if err != nil {
		// maybe it's better to ignore the err, but idk now
		return nil, err
	}
	if suggestions != nil {
		return suggestions, nil
	}

	suggestions, err = c.completer.Complete(ctx, q)
	if err != nil {
		return nil, err
	}

	err = c.set(ctx, q, suggestions)
	if err != nil {
		return nil, err
	}

	return suggestions, nil
}

func (c *CacheableCompleter) get(ctx context.Context, q *Query) (*Suggestions, error) {
	key := queryToCacheKey(q)
	body, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("cant get suggestions from cache: %s", err)
	}

	if len(body) == 0 {
		return nil, nil
	}

	suggestions := Suggestions{}
	if err := json.Unmarshal(body, &suggestions); err != nil {
		return nil, fmt.Errorf("cant unmarshal suggestions: %s", err)
	}

	return &suggestions, nil
}

func (c *CacheableCompleter) set(ctx context.Context, q *Query, suggestions *Suggestions) error {
	body, err := json.Marshal(suggestions)
	if err != nil {
		return fmt.Errorf("cant marshal suggestions: %s", err)
	}

	return c.cache.Set(ctx, queryToCacheKey(q), body, c.ttl)
}

func queryToCacheKey(q *Query) cache.Key {
	key := ""

	// idk if an order is matter for types, but I try to reduce a cache size
	normalizedTypes := make([]string, len(q.Types))
	copy(normalizedTypes, q.Types)
	sort.Sort(sort.StringSlice(normalizedTypes))

	key += strings.Join(normalizedTypes, "") // TODO sort them
	key += q.Term
	key += q.Locale

	return cache.Key(key)
}
