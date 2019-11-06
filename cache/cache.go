package cache

import (
	"errors"
	"time"
)

var (
	ErrCacheKeyNotFound = errors.New("cache key not found")
)

type Cache interface {
	Get(key string) (*string, error)
	Set(key, data string, duration time.Duration) error
}
