package lufeijun

import (
	"github.com/lufeijun/go-tool-redis/illuminate/redis"
)

// 获取一个客户端
func NewClient(opt *redis.Options) *redis.Client {
	opt.Init()

	c := redis.Client{
		// baseClient: &baseClient{
		// 	opt: opt,
		// },
	}
	// c.init()
	// c.connPool = newConnPool(opt, c.dialHook)

	return &c
}

// RedisClientPool = redis.NewClient(&redis.Options{
// 	Addr:         "192.168.0.87:6379",
// 	Password:     "123456",
// 	DB:           1,
// 	PoolSize:     5, // 连接池大小
// 	MinIdleConns: 2, // 最小空闲连接数
// 	MaxIdleConns: 3, // 最大空闲连接数
// })
