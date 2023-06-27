package pool

import (
	"context"
)

// 所有的接口定义

/**
 * @Description: 连接池接口
 * 连接池设计，一个连接池需要具有的基础功能
 */
type Pooler interface {
	NewConn(context.Context) (*Conn, error)
	CloseConn(*Conn) error

	Get(context.Context) (*Conn, error)
	Put(context.Context, *Conn)
	Remove(context.Context, *Conn, error)

	Len() int
	IdleLen() int
	Stats() *Stats

	Close() error
}
