package pool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// 实现 Pooler 接口

type ConnPool struct {
	cfg *Options

	dialErrorsNum uint32 // atomic
	lastDialError atomic.Value

	queue chan struct{}

	connsMu   sync.Mutex
	conns     []*Conn
	idleConns []*Conn

	poolSize     int
	idleConnsLen int

	stats Stats

	_closed uint32 // atomic
}

// 检测 ConnPool 是否实现了 Pooler 接口
var _ Pooler = (*ConnPool)(nil)

// 创建连接
func (p *ConnPool) NewConn(ctx context.Context) (*Conn, error) {
	return p.newConn(ctx, false)
}

func (p *ConnPool) newConn(ctx context.Context, pooled bool) (*Conn, error) {
	cn, err := p.dialConn(ctx, pooled)
	if err != nil {
		return nil, err
	}

	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	// It is not allowed to add new connections to the closed connection pool.
	if p.closed() {
		_ = cn.Close()
		return nil, ErrClosed
	}

	p.conns = append(p.conns, cn)
	if pooled {
		// If pool is full remove the cn on next Put.
		if p.poolSize >= p.cfg.PoolSize {
			cn.pooled = false
		} else {
			p.poolSize++
		}
	}

	return cn, nil
}

// 拨号连接
// --------------------当前位置
func (p *ConnPool) dialConn(ctx context.Context, pooled bool) (*Conn, error) {
	// 判断连接池是否关闭
	if p.closed() {
		return nil, ErrClosed
	}
	if atomic.LoadUint32(&p.dialErrorsNum) >= uint32(p.cfg.PoolSize) {
		return nil, p.getLastDialError()
	}

	netConn, err := p.cfg.Dialer(ctx)
	if err != nil {
		p.setLastDialError(err)
		// 重试逻辑
		if atomic.AddUint32(&p.dialErrorsNum, 1) == uint32(p.cfg.PoolSize) {
			go p.tryDial()
		}
		return nil, err
	}

	cn := NewConn(netConn)
	cn.pooled = pooled
	return cn, nil
}

// 重试逻辑
func (p *ConnPool) tryDial() {
	for {
		if p.closed() {
			return
		}

		conn, err := p.cfg.Dialer(context.Background())
		if err != nil {
			p.setLastDialError(err)
			time.Sleep(time.Second)
			continue
		}

		atomic.StoreUint32(&p.dialErrorsNum, 0)
		_ = conn.Close()
		return
	}
}

func (p *ConnPool) closed() bool {
	// atomic.LoadUint32 原子性获取某个地址的值
	return atomic.LoadUint32(&p._closed) == 1
}

// 额外自有方法

// 错误信息存储、获取
func (p *ConnPool) getLastDialError() error {
	err, _ := p.lastDialError.Load().(*lastDialErrorWrap)
	if err != nil {
		return err.err
	}
	return nil
}
func (p *ConnPool) setLastDialError(err error) {
	p.lastDialError.Store(&lastDialErrorWrap{err: err})
}
