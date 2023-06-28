package redis

import (
	"context"
	"crypto/tls"
	"net"
	"time"
)

// redis client
func NewClient(opt *Options) *Client {
	opt.init()

	c := Client{
		baseClient: &baseClient{
			opt: opt,
		},
	}
	// c.init()
	// c.connPool = newConnPool(opt, c.dialHook)

	return &c
}

// sdk 提供的默认创建连接的方法
func NewDialer(opt *Options) func(context.Context, string, string) (net.Conn, error) {
	// 返回一个函数
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		// net.Dialer go 官方的方法
		netDialer := &net.Dialer{
			Timeout:   opt.DialTimeout,
			KeepAlive: 5 * time.Minute,
		}
		if opt.TLSConfig == nil {
			return netDialer.DialContext(ctx, network, addr)
		}
		return tls.DialWithDialer(netDialer, network, addr, opt.TLSConfig)
	}
}
