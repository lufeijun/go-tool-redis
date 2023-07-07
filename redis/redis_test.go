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

	val, err := RedisClientPool.Get(ctx, "redis").Result()

	fmt.Println(val)

	if err != nil {
		fmt.Println("报错了===", err)
	}

	// time.Sleep(10 * time.Second)

}

// 看一下是否会主动断开连接
func TestOne1(t *testing.T) {

	RedisClientPool := redis.NewClient(&redis.Options{
		Addr:            "192.168.0.87:6379",
		Password:        "123456",
		DB:              1,
		PoolSize:        1,               // 连接池大小
		MinIdleConns:    1,               // 最小空闲连接数
		MaxIdleConns:    1,               // 最大空闲连接数
		ConnMaxLifetime: 3 * time.Second, // 连接最大生存时间
	})

	ctx := context.Background()

	val, err := RedisClientPool.Get(ctx, "name").Result()
	fmt.Println(val)
	if err != nil {
		fmt.Println("报错了===", err)
	}
	time.Sleep(5 * time.Second)

	val, err = RedisClientPool.Get(ctx, "name").Result()
	fmt.Println(val)
	if err != nil {
		fmt.Println("报错了===", err)
	}

	time.Sleep(5 * time.Second)

}

func TestOne2(t *testing.T) {

	RedisClientPool := redis.NewClient(&redis.Options{
		Addr:            "192.168.0.87:6379",
		Password:        "123456",
		DB:              1,
		PoolSize:        1,               // 连接池大小
		MinIdleConns:    1,               // 最小空闲连接数
		MaxIdleConns:    1,               // 最大空闲连接数
		ConnMaxLifetime: 5 * time.Second, // 连接最大生存时间
	})

	time.Sleep(3 * time.Second)

	ctx := context.Background()

	for i := 0; i < 10; i++ {
		val, err := RedisClientPool.Set(ctx, "redis", "世界和平", 0).Result()

		fmt.Println(val)

		if err != nil {
			fmt.Println("报错了===", err)
		}

		time.Sleep(1 * time.Second)
	}

	// time.Sleep(10 * time.Second)

}

func TestOne3(t *testing.T) {

	RedisClientPool := redis.NewClient(&redis.Options{
		Addr:            "192.168.0.87:6379",
		Password:        "123456",
		DB:              1,
		PoolSize:        1,               // 连接池大小
		MinIdleConns:    1,               // 最小空闲连接数
		MaxIdleConns:    1,               // 最大空闲连接数
		ConnMaxLifetime: 5 * time.Second, // 连接最大生存时间
	})

	ctx := context.Background()

	val, err := RedisClientPool.HGet(ctx, "student1", "name").Result()
	if err != nil {
		fmt.Println("报错了===", err)
	}
	fmt.Println("name===", val)

	val, err = RedisClientPool.HGet(ctx, "student1", "age").Result()
	if err != nil {
		fmt.Println("报错了===", err)
	}
	fmt.Println("age===", val)

}

func TestOne4(t *testing.T) {

	RedisClientPool := redis.NewClient(&redis.Options{
		Addr:         "192.168.0.87:6379",
		Password:     "123456",
		DB:           1,
		PoolSize:     1, // 连接池大小
		MinIdleConns: 1, // 最小空闲连接数
		MaxIdleConns: 1, // 最大空闲连接数
		// ConnMaxLifetime: 5 * time.Second, // 连接最大生存时间
	})

	// fmt.Println("======start=====", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("======start=====")

	ctx := context.Background()

	val, err := RedisClientPool.HSet(ctx, "student2", "name", "张三", "age", 22, "city", "北京").Result()
	if err != nil {
		fmt.Println("报错了===", err)
	}

	time.Sleep(15 * time.Second)

	name, err := RedisClientPool.HGet(ctx, "student2", "name").Result()
	if err != nil {
		fmt.Println("报错了===", err)
	}

	fmt.Println("name===", name)

	fmt.Println("成功设置字段===", val)

}

// 获取超过 reader buffer 的数据
func TestOne6(t *testing.T) {
	ctx := context.Background()

	RedisClientPool := redis.NewClient(&redis.Options{
		Addr:         "192.168.0.87:6379",
		Password:     "123456",
		DB:           1,
		PoolSize:     1, // 连接池大小
		MinIdleConns: 1, // 最小空闲连接数
		MaxIdleConns: 1, // 最大空闲连接数
	})

	val, err := RedisClientPool.Get(ctx, "hello").Result()
	fmt.Println(len(val))
	// fmt.Println(val)
	if err != nil {
		fmt.Println("报错了===", err)
	}

}
