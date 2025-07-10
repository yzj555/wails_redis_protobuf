package server

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
)

var ctx = context.Background()

var rdb redis.UniversalClient

type Redis struct {
	Name     string `json:"name"`
	Cluster  bool   `json:"cluster"`
	Addr     string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

var Config *AppConfig

type AppConfig struct {
	CurrentRedis string   `json:"currentRedis"`
	RedisServer  []*Redis `json:"redisServer"`
	Proto        struct {
		Dir string `json:"dir"`
	}
}

func InitRedis() {
	if len(Config.RedisServer) == 0 {
		panic("redis server config is empty")
	}
	defaultRedis := Config.RedisServer[0]
	Config.CurrentRedis = defaultRedis.Name
	rdb = GetRedisClient(defaultRedis)
}

func GetRedisClient(Redis *Redis) redis.UniversalClient {
	var temp redis.UniversalClient
	if Redis.Cluster {
		temp = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{Redis.Addr},
			Password: Redis.Password,
		})
	} else {
		temp = redis.NewClient(&redis.Options{
			Addr:     Redis.Addr,
			Password: Redis.Password,
			DB:       Redis.DB,
		})
	}
	return temp
}

func ChangeRedisServer(name string) {
	for _, server := range Config.RedisServer {
		if server.Name == name {
			Config.CurrentRedis = server.Name
			rdb = GetRedisClient(server)
			return
		}
	}
	panic("redis server not found")
}

// 扫描所有键（集群模式遍历所有主节点，单机模式直接扫描）
func scanAllKeys(ctx context.Context, pattern string, count int64) ([]string, error) {
	allKeys := make(map[string]struct{})
	var mu sync.Mutex
	var wg sync.WaitGroup
	var scanErr error

	var temp []string

	// 集群模式处理
	if clusterClient, ok := rdb.(*redis.ClusterClient); ok {
		err := clusterClient.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
			wg.Add(1)
			go func(c *redis.Client) {
				defer wg.Done()
				keys, err := scanSingleNode(ctx, c, pattern, count)
				if err != nil {
					scanErr = err
					return
				}
				mu.Lock()
				temp = append(temp, keys...)
				mu.Unlock()
			}(client)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("cluster nodes traversal failed: %v", err)
		}
		wg.Wait()
		if scanErr != nil {
			return nil, scanErr
		}
	} else {
		// 单机模式处理
		keys, err := scanSingleNode(ctx, rdb.(*redis.Client), pattern, count)
		if err != nil {
			return nil, err
		}
		temp = append(temp, keys...)
	}
	result := make([]string, 0)
	// 结果去重
	for _, key := range temp {
		if _, ok := allKeys[key]; ok {
			continue
		}
		allKeys[key] = struct{}{}
		result = append(result, key)
	}
	return result, nil
}

// 扫描单个Redis节点
func scanSingleNode(ctx context.Context, client *redis.Client, pattern string, count int64) ([]string, error) {
	var keys []string
	var cursor uint64

	for {
		var err error
		var batch []string

		batch, cursor, err = client.Scan(ctx, cursor, pattern, count).Result()
		if err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}

		keys = append(keys, batch...)
		if cursor == 0 {
			break
		}
	}

	return keys, nil
}
