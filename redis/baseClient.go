package redis

import "github.com/lufeijun/go-tool-redis/redis/pool"

type baseClient struct {
	opt      *Options
	connPool pool.Pooler

	onClose func() error // hook called when client is closed
}
