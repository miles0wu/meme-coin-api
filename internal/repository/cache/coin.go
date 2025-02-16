package cache

import (
	"context"
	"fmt"
	"github.com/miles0wu/meme-coin-api/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

//go:generate mockgen -source=./coin.go -package=cachemocks -destination=./mocks/coin.mock.go CoinCache
type CoinCache interface {
	Set(ctx context.Context, c domain.Coin) error
	Get(ctx context.Context, id int64) (domain.Coin, error)
	IncrPopularityScoreIfPresent(ctx context.Context, id int64) error
}

type RedisCoinCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewRedisCoinCache(client redis.Cmdable) CoinCache {
	return &RedisCoinCache{client: client}
}

func (c *RedisCoinCache) key(id int64) string {
	return fmt.Sprintf("coin:detail:%d", id)
}

func (c *RedisCoinCache) scoreKey(id int64) string {
	return fmt.Sprintf("coin:popularity_score:%d", id)
}

func (c *RedisCoinCache) Set(ctx context.Context, coin domain.Coin) error {
	//TODO implement me
	panic("implement me")
}

func (c *RedisCoinCache) Get(ctx context.Context, id int64) (domain.Coin, error) {
	//TODO implement me
	panic("implement me")
}

func (c *RedisCoinCache) IncrPopularityScoreIfPresent(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}
