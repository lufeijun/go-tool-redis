package redis

import "github.com/lufeijun/go-tool-redis/redis/pool"

func newConn(opt *Options, connPool pool.Pooler) *Conn {
	c := Conn{
		baseClient: baseClient{
			opt:      opt,
			connPool: connPool,
		},
	}

	c.cmdable = c.Process
	c.statefulCmdable = c.Process
	c.initHooks(hooks{
		dial:       c.baseClient.dial,
		process:    c.baseClient.process,
		pipeline:   c.baseClient.processPipeline,
		txPipeline: c.baseClient.processTxPipeline,
	})

	return &c
}
