# 客户端


```go
// Client 属于最外层的对外结构体
type Client struct {
	*baseClient // 基础客户端字段
	cmdable     // 命令接口，type cmdable func(ctx context.Context, cmd Cmder) error
	hooksMixin  // 钩子接口，
}
```


## 方法

```go
func (c *Client) init() // 客户端初始化
func (c *Client) Options() *Options // 获取配置
func (c *Client) WithTimeout(timeout time.Duration) *Client // 设置超时时间
func (c *Client) Conn() *Conn  // 获取连接
func (c *Client) Do(ctx context.Context, args ...interface{}) *Cmd // 执行命令。通用的
func (c *Client) Process(ctx context.Context, cmd Cmder) error // 所有执行命令的入口

```



## baseClient 部分


结构体部分

```go
type baseClient struct {
	opt      *Options // 配置
	connPool pool.Pooler // 连接池

	onClose func() error // hook called when client is closed
}
```


```go
func (c *baseClient) withTimeout(timeout time.Duration) *baseClient // 设置超时时间，给 client 提供更下层逻辑


```


