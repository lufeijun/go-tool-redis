package redis

import (
	"context"
	"crypto/tls"
	"net"
	"runtime"
	"strings"
	"time"
)

type Options struct {

	// ------------基础配置部分--------------
	// 网络类型，tcp 、unix
	Network string
	// 路由地址。待端口号的
	Addr string

	// 设置客户端连接的名称
	ClientName string

	Username string
	Password string
	DB       int

	// ------------连接池部分--------------

	// 创建连接的函数
	Dialer func(ctx context.Context, network, addr string) (net.Conn, error)

	// 创建连接时的钩子函数
	OnConnect func(ctx context.Context, cn *Conn) error

	// 超时时间
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	PoolTimeout     time.Duration // 客户端等待超时时间
	ConnMaxIdleTime time.Duration // 最大空闲时间，超过这个时间，连接会进行关闭
	ConnMaxLifetime time.Duration // 连接最大生存时间，超过这个时间，连接会进行关闭

	// 创建连接时重试次数
	MaxRetries      int
	MinRetryBackoff time.Duration // 最小重试间隔
	MaxRetryBackoff time.Duration // 最大重试间隔

	ContextTimeoutEnabled bool

	// 连接管理方法
	PoolFIFO     bool
	PoolSize     int
	MinIdleConns int
	MaxIdleConns int

	// ------------ 待定部分 --------------

	// Protocol 2 or 3. Use the version to negotiate RESP version with redis-server.
	// Default is 3.
	Protocol int

	// CredentialsProvider allows the username and password to be updated
	// before reconnecting. It should return the current username and password.
	CredentialsProvider func() (username string, password string)

	// TLS Config to use. When set, TLS will be negotiated.
	TLSConfig *tls.Config

	// Limiter interface used to implement circuit breaker or rate limiter.
	// Limiter Limiter

	// Enables read only queries on slave/follower nodes.
	readOnly bool
}

func (opt *Options) Init() {
	if opt.Addr == "" {
		opt.Addr = "localhost:6379"
	}
	if opt.Network == "" {
		if strings.HasPrefix(opt.Addr, "/") {
			opt.Network = "unix"
		} else {
			opt.Network = "tcp"
		}
	}
	if opt.DialTimeout == 0 {
		opt.DialTimeout = 5 * time.Second
	}
	if opt.Dialer == nil {
		// opt.Dialer = NewDialer(opt)
	}
	if opt.PoolSize == 0 {
		opt.PoolSize = 10 * runtime.GOMAXPROCS(0)
	}
	switch opt.ReadTimeout {
	case -2:
		opt.ReadTimeout = -1
	case -1:
		opt.ReadTimeout = 0
	case 0:
		opt.ReadTimeout = 3 * time.Second
	}
	switch opt.WriteTimeout {
	case -2:
		opt.WriteTimeout = -1
	case -1:
		opt.WriteTimeout = 0
	case 0:
		opt.WriteTimeout = opt.ReadTimeout
	}
	if opt.PoolTimeout == 0 {
		if opt.ReadTimeout > 0 {
			opt.PoolTimeout = opt.ReadTimeout + time.Second
		} else {
			opt.PoolTimeout = 30 * time.Second
		}
	}
	if opt.ConnMaxIdleTime == 0 {
		opt.ConnMaxIdleTime = 30 * time.Minute
	}

	if opt.MaxRetries == -1 {
		opt.MaxRetries = 0
	} else if opt.MaxRetries == 0 {
		opt.MaxRetries = 3
	}
	switch opt.MinRetryBackoff {
	case -1:
		opt.MinRetryBackoff = 0
	case 0:
		opt.MinRetryBackoff = 8 * time.Millisecond
	}
	switch opt.MaxRetryBackoff {
	case -1:
		opt.MaxRetryBackoff = 0
	case 0:
		opt.MaxRetryBackoff = 512 * time.Millisecond
	}
}
