package redis

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/lufeijun/go-tool-redis/redis/pool"
	"github.com/lufeijun/go-tool-redis/redis/tool/proto"
	"github.com/lufeijun/go-tool-redis/redis/tool/util"
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

// 设置超时时间
func (c *baseClient) withTimeout(timeout time.Duration) *baseClient {
	opt := c.opt.clone()
	opt.ReadTimeout = timeout
	opt.WriteTimeout = timeout

	clone := c.clone()
	clone.opt = opt

	return clone
}

// 返回 redis 连接信息
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

// 初始化连接池
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
	// 发送简单命令，判断是否支持 HELLO 命令，判断是否可以连接上 redis
	if err := conn.Hello(ctx, protocol, username, password, "").Err(); err == nil {
		auth = true
	} else if !isRedisError(err) {
		return err
	}

	// 管道
	// 1、用户密码认证
	// 2、选择数据库
	// 3、设置只读
	// 4、设置客户端名称
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

	// 如果有链接钩子函数，执行钩子函数
	if c.opt.OnConnect != nil {
		return c.opt.OnConnect(ctx, conn)
	}
	return nil
}

// 包含限流逻辑，暂时忽略
func (c *baseClient) getConn(ctx context.Context) (*pool.Conn, error) {

	// 如果有限流器，先判断是否可以获取连接
	// if c.opt.Limiter != nil {
	// 	err := c.opt.Limiter.Allow()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	cn, err := c._getConn(ctx)
	if err != nil {
		// if c.opt.Limiter != nil {
		// 	c.opt.Limiter.ReportResult(err)
		// }
		return nil, err
	}

	return cn, nil
}

// 实际处理连接的函数，私有
func (c *baseClient) _getConn(ctx context.Context) (*pool.Conn, error) {
	// 从连接池获取，目测应该是阻塞的
	cn, err := c.connPool.Get(ctx)
	if err != nil {
		return nil, err
	}

	// 如果已经初始化过了，直接返回
	if cn.Inited {
		return cn, nil
	}

	// 初始化连接
	if err := c.initConn(ctx, cn); err != nil {
		c.connPool.Remove(ctx, cn, err)
		if err := errors.Unwrap(err); err != nil {
			return nil, err
		}
		return nil, err
	}

	return cn, nil
}

// 释放连接
func (c *baseClient) releaseConn(ctx context.Context, cn *pool.Conn, err error) {
	// if c.opt.Limiter != nil {
	// 	c.opt.Limiter.ReportResult(err)
	// }

	// 释放逻辑
	// 如果是错误连接，直接移除
	// 如果是正常连接，放回连接池
	if isBadConn(err, false, c.opt.Addr) {
		c.connPool.Remove(ctx, cn, err)
	} else {
		c.connPool.Put(ctx, cn)
	}
}

func (c *baseClient) withConn(
	ctx context.Context, fn func(context.Context, *pool.Conn) error,
) error {
	cn, err := c.getConn(ctx)
	if err != nil {
		return err
	}

	var fnErr error
	defer func() {
		c.releaseConn(ctx, cn, fnErr)
	}()

	fnErr = fn(ctx, cn)

	return fnErr
}

// 拨号
func (c *baseClient) dial(ctx context.Context, network, addr string) (net.Conn, error) {
	return c.opt.Dialer(ctx, network, addr)
}

// 执行命令
func (c *baseClient) process(ctx context.Context, cmd Cmder) error {
	var lastErr error
	for attempt := 0; attempt <= c.opt.MaxRetries; attempt++ {
		attempt := attempt

		retry, err := c._process(ctx, cmd, attempt)
		if err == nil || !retry {
			return err
		}

		lastErr = err
	}
	return lastErr
}

// 执行命令函数主体，私有
func (c *baseClient) _process(ctx context.Context, cmd Cmder, attempt int) (bool, error) {
	if attempt > 0 {
		if err := util.Sleep(ctx, c.retryBackoff(attempt)); err != nil {
			return false, err
		}
	}

	retryTimeout := uint32(0)
	if err := c.withConn(ctx, func(ctx context.Context, cn *pool.Conn) error {
		// 具体执行命令逻辑

		// 发送 redis 命令
		if err := cn.WithWriter(c.context(ctx), c.opt.WriteTimeout, func(wr *proto.Writer) error {
			return writeCmd(wr, cmd)
		}); err != nil {
			atomic.StoreUint32(&retryTimeout, 1)
			return err
		}

		// 读取 redis 命令返回值
		if err := cn.WithReader(c.context(ctx), c.cmdTimeout(cmd), cmd.readReply); err != nil {
			if cmd.readTimeout() == nil {
				atomic.StoreUint32(&retryTimeout, 1)
			} else {
				atomic.StoreUint32(&retryTimeout, 0)
			}
			return err
		}

		return nil
	}); err != nil { // if 的 end 部分
		retry := shouldRetry(err, atomic.LoadUint32(&retryTimeout) == 1)
		return retry, err
	}

	return false, nil
}

func (c *baseClient) retryBackoff(attempt int) time.Duration {
	return util.RetryBackoff(attempt, c.opt.MinRetryBackoff, c.opt.MaxRetryBackoff)
}

func (c *baseClient) context(ctx context.Context) context.Context {
	if c.opt.ContextTimeoutEnabled {
		return ctx
	}
	return context.Background()
}

func (c *baseClient) cmdTimeout(cmd Cmder) time.Duration {
	if timeout := cmd.readTimeout(); timeout != nil {
		t := *timeout
		if t == 0 {
			return 0
		}
		return t + 10*time.Second
	}
	return c.opt.ReadTimeout
}

func (c *baseClient) Close() error {
	var firstErr error
	if c.onClose != nil {
		if err := c.onClose(); err != nil {
			firstErr = err
		}
	}
	if err := c.connPool.Close(); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}

// 管道查询
func (c *baseClient) processPipeline(ctx context.Context, cmds []Cmder) error {
	if err := c.generalProcessPipeline(ctx, cmds, c.pipelineProcessCmds); err != nil {
		return err
	}
	return cmdsFirstErr(cmds)
}

func (c *baseClient) processTxPipeline(ctx context.Context, cmds []Cmder) error {
	if err := c.generalProcessPipeline(ctx, cmds, c.txPipelineProcessCmds); err != nil {
		return err
	}
	return cmdsFirstErr(cmds)
}

type pipelineProcessor func(context.Context, *pool.Conn, []Cmder) (bool, error)

func (c *baseClient) generalProcessPipeline(
	ctx context.Context, cmds []Cmder, p pipelineProcessor,
) error {
	var lastErr error
	for attempt := 0; attempt <= c.opt.MaxRetries; attempt++ {
		if attempt > 0 {
			if err := util.Sleep(ctx, c.retryBackoff(attempt)); err != nil {
				setCmdsErr(cmds, err)
				return err
			}
		}

		// Enable retries by default to retry dial errors returned by withConn.
		canRetry := true
		lastErr = c.withConn(ctx, func(ctx context.Context, cn *pool.Conn) error {
			var err error
			canRetry, err = p(ctx, cn, cmds)
			return err
		})
		if lastErr == nil || !canRetry || !shouldRetry(lastErr, true) {
			return lastErr
		}
	}
	return lastErr
}

// 管道命令群
func (c *baseClient) pipelineProcessCmds(
	ctx context.Context, cn *pool.Conn, cmds []Cmder,
) (bool, error) {

	// 发送 redis 命令群
	if err := cn.WithWriter(c.context(ctx), c.opt.WriteTimeout, func(wr *proto.Writer) error {
		return writeCmds(wr, cmds)
	}); err != nil {
		setCmdsErr(cmds, err)
		return true, err
	}

	// 读取 redis 命令群返回值
	if err := cn.WithReader(c.context(ctx), c.opt.ReadTimeout, func(rd *proto.Reader) error {
		return pipelineReadCmds(rd, cmds)
	}); err != nil {
		return true, err
	}

	return false, nil
}

func (c *baseClient) txPipelineProcessCmds(
	ctx context.Context, cn *pool.Conn, cmds []Cmder,
) (bool, error) {
	if err := cn.WithWriter(c.context(ctx), c.opt.WriteTimeout, func(wr *proto.Writer) error {
		return writeCmds(wr, cmds)
	}); err != nil {
		setCmdsErr(cmds, err)
		return true, err
	}

	if err := cn.WithReader(c.context(ctx), c.opt.ReadTimeout, func(rd *proto.Reader) error {
		statusCmd := cmds[0].(*StatusCmd)
		// Trim multi and exec.
		trimmedCmds := cmds[1 : len(cmds)-1]

		if err := txPipelineReadQueued(rd, statusCmd, trimmedCmds); err != nil {
			setCmdsErr(cmds, err)
			return err
		}

		return pipelineReadCmds(rd, trimmedCmds)
	}); err != nil {
		return false, err
	}

	return false, nil
}

func txPipelineReadQueued(rd *proto.Reader, statusCmd *StatusCmd, cmds []Cmder) error {
	// Parse +OK.
	if err := statusCmd.readReply(rd); err != nil {
		return err
	}

	// Parse +QUEUED.
	for range cmds {
		if err := statusCmd.readReply(rd); err != nil && !isRedisError(err) {
			return err
		}
	}

	// Parse number of replies.
	line, err := rd.ReadLine()
	if err != nil {
		if err == Nil {
			err = TxFailedErr
		}
		return err
	}

	if line[0] != proto.RespArray {
		return fmt.Errorf("redis: expected '*', but got line %q", line)
	}

	return nil
}
