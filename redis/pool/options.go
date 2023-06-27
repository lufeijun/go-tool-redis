package pool

import (
	"context"
	"net"
	"time"
)

type Options struct {
	Dialer func(context.Context) (net.Conn, error)

	PoolFIFO        bool
	PoolSize        int
	PoolTimeout     time.Duration
	MinIdleConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}
