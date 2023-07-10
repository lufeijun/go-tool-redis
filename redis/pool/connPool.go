package pool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lufeijun/go-tool-redis/redis/tool/logger"
)

// 实现 Pooler 接口

type ConnPool struct {
	cfg *Options

	dialErrorsNum uint32 // atomic
	lastDialError atomic.Value

	// 一个通道，用于控制连接池的并发数
	queue chan struct{}

	connsMu   sync.Mutex
	conns     []*Conn
	idleConns []*Conn

	poolSize     int
	idleConnsLen int

	stats Stats

	_closed uint32 // atomic 一个开关标志位，表示连接池是否已经关闭
}

// 检测 ConnPool 是否实现了 Pooler 接口
var _ Pooler = (*ConnPool)(nil)

// 创建连接 NewConn 方法重写
func (p *ConnPool) NewConn(ctx context.Context) (*Conn, error) {
	return p.newConn(ctx, false)
}

// 创建新连接的具体方法
func (p *ConnPool) newConn(ctx context.Context, pooled bool) (*Conn, error) {
	cn, err := p.dialConn(ctx, pooled)
	if err != nil {
		return nil, err
	}

	// 上锁
	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	// 连接池的状态，是否关闭与否
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
// pooled 是否将连接放入连接池 true 放入，false 不放入
func (p *ConnPool) dialConn(ctx context.Context, pooled bool) (*Conn, error) {
	// 判断连接池是否关闭
	if p.closed() {
		return nil, ErrClosed
	}

	// 这个判断逻辑有点无语
	if atomic.LoadUint32(&p.dialErrorsNum) >= uint32(p.cfg.PoolSize) {
		return nil, p.getLastDialError()
	}

	// 进行连接
	netConn, err := p.cfg.Dialer(ctx)
	if err != nil {
		p.setLastDialError(err)
		// 重试逻辑，应该判断重试次数吧
		if atomic.AddUint32(&p.dialErrorsNum, 1) == uint32(p.cfg.PoolSize) {
			go p.tryDial()
		}
		return nil, err
	}

	// 构建 conn 结构体，这个结构体是对 net.Conn 的封装，也是 connpool 管理的对象
	cn := NewConn(netConn)
	cn.pooled = pooled // 此字段用于处理回收逻辑
	return cn, nil
}

// 每隔 1 s，不断的尝试进行连接，如果成功后，将 dialErrorsNum 置为 0
func (p *ConnPool) tryDial() {
	for {
		if p.closed() {
			return
		}

		conn, err := p.cfg.Dialer(context.Background())
		if err != nil {
			p.setLastDialError(err)
			// 每次重试之间间隔 1 秒
			time.Sleep(time.Second)
			continue
		}

		// 设置连接错误次数为 0
		atomic.StoreUint32(&p.dialErrorsNum, 0)

		// 这不就关闭连接了么
		_ = conn.Close()
		return
	}
}

// 重写 CloseConn 方法
func (p *ConnPool) CloseConn(cn *Conn) error {
	p.removeConnWithLock(cn)
	return p.closeConn(cn)
}

// tcp 连接的关闭
func (p *ConnPool) closeConn(cn *Conn) error {
	return cn.Close()
}

// 加锁，从 connpool 中删除一个连接
func (p *ConnPool) removeConnWithLock(cn *Conn) {
	p.connsMu.Lock()         // 加锁
	defer p.connsMu.Unlock() // 延迟解锁

	// 加锁后搞事情
	p.removeConn(cn)
}

// 从连接池中删除一个连接
func (p *ConnPool) removeConn(cn *Conn) {
	// 遍历连接池中的连接
	for i, c := range p.conns {
		if c == cn {
			// 这个删除操作，好魔幻
			p.conns = append(p.conns[:i], p.conns[i+1:]...)

			// 如果此连接在连接池中管理
			if cn.pooled {
				p.poolSize--
				p.checkMinIdleConns() // 检查最小空闲连接数
			}
			break
		}
	}
	atomic.AddUint32(&p.stats.StaleConns, 1)
}

// 判断连接池是否关闭，方法名改为 isClosed 更好
func (p *ConnPool) closed() bool {
	// atomic.LoadUint32 原子性获取某个地址的值
	return atomic.LoadUint32(&p._closed) == 1
}

// 从连接池中获取一个连接
func (p *ConnPool) Get(ctx context.Context) (*Conn, error) {
	if p.closed() {
		return nil, ErrClosed
	}

	// 等待一个资源，有可能超时了
	if err := p.waitTurn(ctx); err != nil {
		return nil, err
	}

	// 循环，加锁，从连接池中获取一个连接
	for {
		p.connsMu.Lock()
		cn, err := p.popIdle()
		p.connsMu.Unlock()

		if err != nil {
			return nil, err
		}

		// 表示没有空闲链接了，需要新建一个连接
		if cn == nil {
			break
		}

		// 链接监控检测，同时也设置了 usedAt 时间
		if !p.isHealthyConn(cn) {
			_ = p.CloseConn(cn)
			continue
		}

		// 连接池命中数 +1
		atomic.AddUint32(&p.stats.Hits, 1)
		return cn, nil
	}

	// 连接池未命中数 +1
	atomic.AddUint32(&p.stats.Misses, 1)

	// 连接池没有，新建一个连接
	newcn, err := p.newConn(ctx, true)
	if err != nil {
		p.freeTurn()
		return nil, err
	}

	return newcn, nil
}

// 阻塞等待，从连接池获取资源
func (p *ConnPool) waitTurn(ctx context.Context) error {
	// 上下文判断

	// 如果 ctx.Done() 有值，直接返回
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 如果管道可以直接写入，说明资源是够用的，直接返回
	select {
	case p.queue <- struct{}{}:
		return nil
	default:
	}

	// 计时器 sync.Pool 用于减少内存分配
	timer := timers.Get().(*time.Timer)
	timer.Reset(p.cfg.PoolTimeout)

	// 进行阻塞等待
	select {
	case <-ctx.Done(): // 上下文超时
		if !timer.Stop() {
			<-timer.C
		}
		timers.Put(timer)
		return ctx.Err()
	case p.queue <- struct{}{}: // 有资源可用
		if !timer.Stop() {
			<-timer.C
		}
		timers.Put(timer)
		return nil
	case <-timer.C: // 配置的时间超时
		timers.Put(timer)
		atomic.AddUint32(&p.stats.Timeouts, 1)
		return ErrPoolTimeout
	}
}

// 获取空闲链接
func (p *ConnPool) popIdle() (*Conn, error) {
	if p.closed() {
		return nil, ErrClosed
	}
	n := len(p.idleConns)
	// 如果没有空闲连接，直接返回
	if n == 0 {
		return nil, nil
	}

	var cn *Conn
	// 按照配置方式来取连接，FIFO 或者 LIFO
	if p.cfg.PoolFIFO {
		cn = p.idleConns[0]
		copy(p.idleConns, p.idleConns[1:])
		p.idleConns = p.idleConns[:n-1]
	} else {
		idx := n - 1
		cn = p.idleConns[idx]
		p.idleConns = p.idleConns[:idx]
	}
	p.idleConnsLen--
	p.checkMinIdleConns() // 这个感觉可以新开协程进行处理
	return cn, nil
}

// 使用完毕之后，放回连接池
func (p *ConnPool) Put(ctx context.Context, cn *Conn) {
	// 判断是否还有未处理的数据
	if cn.rd.Buffered() > 0 {
		logger.Logger.Printf(ctx, "Conn has unread data")
		p.Remove(ctx, cn, BadConnError{})
		return
	}

	// 表示这个 tcp 连接不是从连接池中获取的，直接关闭
	if !cn.pooled {
		p.Remove(ctx, cn, nil)
		return
	}

	var shouldCloseConn bool

	p.connsMu.Lock()

	// 是放回连接池，还是关闭连接
	if p.cfg.MaxIdleConns == 0 || p.idleConnsLen < p.cfg.MaxIdleConns {
		p.idleConns = append(p.idleConns, cn)
		p.idleConnsLen++
	} else {
		p.removeConn(cn)
		shouldCloseConn = true
	}

	p.connsMu.Unlock()

	p.freeTurn()

	if shouldCloseConn {
		_ = p.closeConn(cn)
	}
}

// 关闭连接
func (p *ConnPool) Remove(_ context.Context, cn *Conn, reason error) {
	p.removeConnWithLock(cn) // 从连接池中移除
	p.freeTurn()             // 释放资源，表示可以从连接池中获取资源了
	_ = p.closeConn(cn)      // 关闭连接
}

// 获取当前连接池中的连接数
func (p *ConnPool) Len() int {
	p.connsMu.Lock()
	n := len(p.conns)
	p.connsMu.Unlock()
	return n
}

// 获取当前连接池中的空闲连接数
func (p *ConnPool) IdleLen() int {
	p.connsMu.Lock()
	n := p.idleConnsLen
	p.connsMu.Unlock()
	return n
}

// 当前情况下的状态信息
func (p *ConnPool) Stats() *Stats {
	return &Stats{
		Hits:     atomic.LoadUint32(&p.stats.Hits),
		Misses:   atomic.LoadUint32(&p.stats.Misses),
		Timeouts: atomic.LoadUint32(&p.stats.Timeouts),

		TotalConns: uint32(p.Len()),
		IdleConns:  uint32(p.IdleLen()),
		StaleConns: atomic.LoadUint32(&p.stats.StaleConns),
	}
}

// 关闭连接池
func (p *ConnPool) Close() error {

	// 防止重复关闭
	if !atomic.CompareAndSwapUint32(&p._closed, 0, 1) {
		return ErrClosed
	}

	var firstErr error
	p.connsMu.Lock()
	// 关闭所有连接
	for _, cn := range p.conns {
		if err := p.closeConn(cn); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	p.conns = nil
	p.poolSize = 0
	p.idleConns = nil
	p.idleConnsLen = 0
	p.connsMu.Unlock()

	return firstErr
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

// 检测最小空闲连接数是否满足要求
func (p *ConnPool) checkMinIdleConns() {

	if p.cfg.MinIdleConns == 0 {
		return
	}

	for p.poolSize < p.cfg.PoolSize && p.idleConnsLen < p.cfg.MinIdleConns {
		select {
		case p.queue <- struct{}{}: // 占用一个资源
			p.poolSize++
			p.idleConnsLen++

			go func() {
				err := p.addIdleConn()
				if err != nil && err != ErrClosed {
					p.connsMu.Lock()
					p.poolSize--
					p.idleConnsLen--
					p.connsMu.Unlock()
				}

				p.freeTurn() // 释放资源
			}()
		default:
			return
		}
	}
}

// 添加一个空闲连接
func (p *ConnPool) addIdleConn() error {
	cn, err := p.dialConn(context.TODO(), true)
	if err != nil {
		return err
	}

	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	// It is not allowed to add new connections to the closed connection pool.
	if p.closed() {
		_ = cn.Close()
		return ErrClosed
	}

	p.conns = append(p.conns, cn)
	p.idleConns = append(p.idleConns, cn)
	return nil
}

// 出管道 <-p.queue，表示释放一个资源
func (p *ConnPool) freeTurn() {
	<-p.queue
}

// 连接的健康检查
func (p *ConnPool) isHealthyConn(cn *Conn) bool {
	now := time.Now()

	// 判断链接存在时间是否超时
	if p.cfg.ConnMaxLifetime > 0 && now.Sub(cn.createdAt) >= p.cfg.ConnMaxLifetime {
		return false
	}

	// 判断链接空闲时间是否超时
	if p.cfg.ConnMaxIdleTime > 0 && now.Sub(cn.UsedAt()) >= p.cfg.ConnMaxIdleTime {
		return false
	}

	if connCheck(cn.netConn) != nil {
		return false
	}

	// 设置连接最后一次使用时间
	cn.SetUsedAt(now)
	return true
}
