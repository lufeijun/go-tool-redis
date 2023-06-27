package redis

import (
	"context"
	"crypto/tls"
	"net"
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
