package wcontext

import (
	"github.com/lufeijun/go-tool-wechat/wechat/credential"
	"github.com/lufeijun/go-tool-wechat/wechat/miniprogram/miniconfig"
)

// Context struct
type Context struct {
	*miniconfig.Config
	credential.AccessTokenHandle
}
