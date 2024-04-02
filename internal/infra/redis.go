package infra

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sapphire-server/internal/conf"
)

var (
	Redis *redis.Client
	Ctx   context.Context
)

func InitRedis() error {
	opt, err := redis.ParseURL(conf.GetRedisConfig())
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)
	Redis = rdb

	Ctx = context.Background()

	return nil
}
