package ioc

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var c Config
	err := viper.UnmarshalKey("redis", &c)
	if err != nil {
		panic(fmt.Errorf("init redis failed %v", err))
	}

	return redis.NewClient(&redis.Options{
		Addr: c.Addr,
	})
}
