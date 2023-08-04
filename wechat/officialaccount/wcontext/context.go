package wcontext

import (
	"github.com/lufeijun/go-tool-wechat/wechat/credential"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/offConfig"
)

// Context struct
type Context struct {
	*offConfig.Config
	credential.AccessTokenHandle
}
