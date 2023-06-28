package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/lufeijun/go-tool-redis/redis/pool"
)

type baseClient struct {
	opt      *Options
	connPool pool.Pooler

	onClose func() error // hook called when client is closed
}

// 克隆
func (c *baseClient) clone() *baseClient {
	clone := *c
	return &clone
}

func (c *baseClient) withTimeout(timeout time.Duration) *baseClient {
	opt := c.opt.clone()
	opt.ReadTimeout = timeout
	opt.WriteTimeout = timeout

	clone := c.clone()
	clone.opt = opt

	return clone
}

func (c *baseClient) String() string {
	return fmt.Sprintf("Redis<%s db:%d>", c.getAddr(), c.opt.DB)
}

func (c *baseClient) getAddr() string {
	return c.opt.Addr
}

// 新连接
func (c *baseClient) newConn(ctx context.Context) (*pool.Conn, error) {
	// 从连接池中获取一个连接
	cn, err := c.connPool.NewConn(ctx)
	if err != nil {
		return nil, err
	}

	// 初始化连接
	err = c.initConn(ctx, cn)
	if err != nil {
		_ = c.connPool.CloseConn(cn)
		return nil, err
	}

	return cn, nil
}

// 初始化连接
func (c *baseClient) initConn(ctx context.Context, cn *pool.Conn) error {
	// 如果已经初始化过了，直接返回
	if cn.Inited {
		return nil
	}
	cn.Inited = true

	username, password := c.opt.Username, c.opt.Password
	if c.opt.CredentialsProvider != nil {
		username, password = c.opt.CredentialsProvider()
	}

	// 单链接池
	connPool := pool.NewSingleConnPool(c.connPool, cn)
	conn := newConn(c.opt, connPool)

	var auth bool
	protocol := c.opt.Protocol
	// By default, use RESP3 in current version.
	if protocol < 2 {
		protocol = 3
	}

	// for redis-server versions that do not support the HELLO command,
	// RESP2 will continue to be used.
	if err := conn.Hello(ctx, protocol, username, password, "").Err(); err == nil {
		auth = true
	} else if !isRedisError(err) {
		// When the server responds with the RESP protocol and the result is not a normal
		// execution result of the HELLO command, we consider it to be an indication that
		// the server does not support the HELLO command.
		// The server may be a redis-server that does not support the HELLO command,
		// or it could be DragonflyDB or a third-party redis-proxy. They all respond
		// with different error string results for unsupported commands, making it
		// difficult to rely on error strings to determine all results.
		return err
	}

	_, err := conn.Pipelined(ctx, func(pipe Pipeliner) error {
		if !auth && password != "" {
			if username != "" {
				pipe.AuthACL(ctx, username, password)
			} else {
				pipe.Auth(ctx, password)
			}
		}

		if c.opt.DB > 0 {
			pipe.Select(ctx, c.opt.DB)
		}

		if c.opt.readOnly {
			pipe.ReadOnly(ctx)
		}

		if c.opt.ClientName != "" {
			pipe.ClientSetName(ctx, c.opt.ClientName)
		}

		return nil
	})
	if err != nil {
		return err
	}

	if c.opt.OnConnect != nil {
		return c.opt.OnConnect(ctx, conn)
	}
	return nil
}
