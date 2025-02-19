package cache

import (
	"context"
	"encoding/json"
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
	Del(ctx context.Context, id int64) error
}

type RedisCoinCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewRedisCoinCache(client redis.Cmdable) CoinCache {
	return &RedisCoinCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (c *RedisCoinCache) key(id int64) string {
	return fmt.Sprintf("coin:%d", id)
}

func (c *RedisCoinCache) Set(ctx context.Context, coin domain.Coin) error {
	bs, err := json.Marshal(coin)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, c.key(coin.Id), bs, c.expiration).Err()
}

func (c *RedisCoinCache) Get(ctx context.Context, id int64) (domain.Coin, error) {
	val, err := c.client.Get(ctx, c.key(id)).Bytes()
	if err != nil {
		return domain.Coin{}, err
	}
	var coin domain.Coin
	err = json.Unmarshal(val, &coin)
	if err != nil {
		return domain.Coin{}, err
	}
	return coin, nil
}

func (c *RedisCoinCache) Del(ctx context.Context, id int64) error {
	return c.client.Del(ctx, c.key(id)).Err()
}
