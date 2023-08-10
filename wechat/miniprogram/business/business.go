package business

import "github.com/lufeijun/go-tool-wechat/wechat/miniprogram/wcontext"

// Business 业务
type Business struct {
	*wcontext.Context
}

// NewBusiness init
func NewBusiness(ctx *wcontext.Context) *Business {
	return &Business{ctx}
}
