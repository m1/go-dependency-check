package cache

import (
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	*redis.Client
}

func NewRedis(host string) (*Redis, error) {
	r := &Redis{
		Client: redis.NewClient(&redis.Options{
			Addr: host,
		}),
	}

	_, err := r.Ping().Result()

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r Redis) Get(key string) (*string, error) {
	cache, err := r.Client.Get(key).Result()
	if err == redis.Nil {
		return nil, ErrCacheKeyNotFound
	} else if err != nil {
		return nil, err
	}

	return &cache, nil
}

func (r Redis) Set(key string, data string, duration time.Duration) error {
	status := r.Client.Set(key, data, duration)
	return status.Err()
}
