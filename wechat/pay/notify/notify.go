package notify

import "github.com/lufeijun/go-tool-wechat/wechat/pay/payconfig"

// Notify 回调
type Notify struct {
	*payconfig.Config
}

// NewNotify new
func NewNotify(cfg *payconfig.Config) *Notify {
	return &Notify{cfg}
}
