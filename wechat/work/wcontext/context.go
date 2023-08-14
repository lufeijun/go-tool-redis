package wcontext

import (
	"github.com/lufeijun/go-tool-wechat/wechat/credential"
	"github.com/lufeijun/go-tool-wechat/wechat/work/workconfig"
)

// Context struct
type Context struct {
	*workconfig.Config
	credential.AccessTokenHandle
}
