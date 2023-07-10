# 连接池


# 连接池技术


## 资源的基础管理

1. 资源的初始化
  tcp 链接的建立
  用户认证逻辑
  资源基础设置

2. 资源的回收
  资源的释放：关闭 tcp 链接

3. 发送请求|接收响应
	通过 tcp 链接发送请求
	通过 tcp 链接接收响应


## 连接池功能

1. 连接池本身的初始化
	连接池的大小
	最大空闲时间
	最小|大空闲连接数
	连接池超时时间
	......

2. 对资源的控制
  资源数量的控制
		1、最大使用量控制
  资源有效性的检查
		1、心跳检查，tcp 是否仍然有效
		2、是否超过生命周期
			最大空闲时间
			最大生存时间


3. 资源的使用
  分配：从资源池中获取一个资源
		1、申请资源时，需要先加锁，在进行资源分配，线性安全问题
		2、如果超过资源限制，需要阻塞等待
  回收：将资源放回资源池中
		1、是否需要放入资源池
		2、关闭连接



# go-redis 设计

## 结构体设计


1、对连接的抽象

```go
type Conn struct {
	usedAt  int64 // 时间戳，用于周期管理
	netConn net.Conn // tcp 连接信息

  // 绑定在 tcp 上的 IO 读写流
	rd *proto.Reader // 读取流，用于读取 redis 数据
	bw *bufio.Writer
	wr *proto.Writer

	Inited    bool // 标志位，链接是否被初始化，ex：用户密码认证，客户端名称、DB 库选择等
	pooled    bool // 标志位，链接是否被放回连接池
	createdAt time.Time // 链接创建时间
}

```

2、连接池配置信息
  
```go
type Options struct {
	Dialer func(context.Context) (net.Conn, error) // 创建连接的方法

	PoolFIFO        bool // 获取连接策略
	PoolSize        int  // 连接池大小
	PoolTimeout     time.Duration // 连接池超时时间
	MinIdleConns    int // 最小空闲连接数
	MaxIdleConns    int // 最大空闲连接数
	ConnMaxIdleTime time.Duration // 最大空闲时间
	ConnMaxLifetime time.Duration // 最大生存时间
}

```

3、对连接池的抽象

```go
type ConnPool struct {
	cfg *Options // 配置信息

	dialErrorsNum uint32 // atomic
	lastDialError atomic.Value

	// 管道字段，管道用于多线程信息同步
  // 当从连接池获取一个连接时，执行 queue <- struct{}{} ，表示获取一个资源，如果当前管道已经写满了，说明已经达到上限，则需要阻塞等待
  // 资源使用完成后，执行 <-queue，表示释放一个资源，如果当前管道已经没有资源了，则需要阻塞等待
	queue chan struct{}

	connsMu   sync.Mutex // 锁控制

  // 连接池当前连接数和空闲连接数
  conns     []*Conn
	idleConns []*Conn
	poolSize     int
	idleConnsLen int

	stats Stats // 统计信息

	_closed uint32 // atomic 一个开关标志位，表示连接池是否已经关闭
}
```

