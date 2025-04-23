package server

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var rdb redis.UniversalClient

type Config struct {
	Redis struct {
		Cluster  bool   `json:"cluster"`
		Addr     string `json:"address"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
	Proto struct {
		Dir string `json:"dir"`
	}
}

func InitRedis(config *Config) {
	if config.Redis.Cluster {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{config.Redis.Addr},
			Password: config.Redis.Password,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     config.Redis.Addr,
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
		})
	}
}
