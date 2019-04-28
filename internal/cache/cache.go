package cache

import (
	"context"
	"time"
)

type (
	Cache interface {
		Get(context.Context, Key) ([]byte, error)
		Set(context.Context, Key, []byte, time.Duration) error
	}

	Key string
)
