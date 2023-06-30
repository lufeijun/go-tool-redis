package redis_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lufeijun/go-tool-redis/redis"
)

func TestOne(t *testing.T) {

	RedisClientPool := redis.NewClient(&redis.Options{
		Addr:         "192.168.0.87:6379",
		Password:     "123456",
		DB:           1,
		PoolSize:     5, // 连接池大小
		MinIdleConns: 2, // 最小空闲连接数
		MaxIdleConns: 3, // 最大空闲连接数
	})

	ctx := context.Background()

	val, err := RedisClientPool.Get(ctx, "name").Result()

	fmt.Println(val)

	if err != nil {
		fmt.Println("报错了===", err)
	}

	time.Sleep(10 * time.Second)

}
