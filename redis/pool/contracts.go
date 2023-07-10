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
	NewConn(context.Context) (*Conn, error) // 创建新连接，连接不一定在连接池中，这个函数额外调用过多，会导致实际连接数超标
	CloseConn(*Conn) error                  // 关闭连接

	Get(context.Context) (*Conn, error)   // 从连接池获取一个连接
	Put(context.Context, *Conn)           // 将放回连接池
	Remove(context.Context, *Conn, error) // 将连接移除连接池

	Len() int
	IdleLen() int
	Stats() *Stats

	Close() error
}
