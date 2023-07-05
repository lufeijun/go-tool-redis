package redis_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/lufeijun/go-tool-redis/redis"
)

type redisHook struct{}

var _ redis.Hook = redisHook{}

func (redisHook) DialHook(hook redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		fmt.Printf("dialing %s %s\n", network, addr)
		conn, err := hook(ctx, network, addr)
		fmt.Printf("finished dialing %s %s\n", network, addr)
		return conn, err
	}
}

func (redisHook) ProcessHook(hook redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		fmt.Printf("starting processing: <%s>\n", cmd)
		err := hook(ctx, cmd)
		fmt.Printf("finished processing: <%s>\n", cmd)
		return err
	}
}

func (redisHook) ProcessPipelineHook(hook redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		fmt.Printf("pipeline starting processing: %v\n", cmds)
		err := hook(ctx, cmds)
		fmt.Printf("pipeline finished processing: %v\n", cmds)
		return err
	}
}

func Test1(t *testing.T) {
	RedisClientPool := redis.NewClient(&redis.Options{
		Addr:            "192.168.0.87:6379",
		Password:        "123456",
		DB:              1,
		PoolSize:        1,               // 连接池大小
		MinIdleConns:    1,               // 最小空闲连接数
		MaxIdleConns:    1,               // 最大空闲连接数
		ConnMaxLifetime: 3 * time.Second, // 连接最大生存时间
	})

	RedisClientPool.AddHook(redisHook{})

	ctx := context.Background()

	val, err := RedisClientPool.Get(ctx, "hello").Result()
	fmt.Println(val)
	if err != nil {
		fmt.Println("报错了===", err)
	}

}
