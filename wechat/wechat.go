package wechat

import (
	"github.com/lufeijun/go-tool-wechat/wechat/cache"
	"github.com/lufeijun/go-tool-wechat/wechat/miniprogram"
	"github.com/lufeijun/go-tool-wechat/wechat/miniprogram/miniconfig"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/offConfig"
	"github.com/lufeijun/go-tool-wechat/wechat/pay"
	"github.com/lufeijun/go-tool-wechat/wechat/pay/payconfig"
)

// Wechat struct
type Wechat struct {
	cache cache.Cache
}

// NewWechat init
func NewWechat() *Wechat {
	return &Wechat{}
}

// SetCache 设置cache
func (wc *Wechat) SetCache(cache cache.Cache) {
	wc.cache = cache
}

// GetOfficialAccount 获取微信公众号实例
func (wc *Wechat) GetOfficialAccount(cfg *offConfig.Config) *officialaccount.OfficialAccount {
	if cfg.Cache == nil {
		cfg.Cache = wc.cache
	}
	return officialaccount.NewOfficialAccount(cfg)
}

// GetMiniProgram 获取小程序的实例
func (wc *Wechat) GetMiniProgram(cfg *miniconfig.Config) *miniprogram.MiniProgram {
	if cfg.Cache == nil {
		cfg.Cache = wc.cache
	}
	return miniprogram.NewMiniProgram(cfg)
}

// GetPay 获取微信支付的实例
func (wc *Wechat) GetPay(cfg *payconfig.Config) *pay.Pay {
	return pay.NewPay(cfg)
}

// // GetOpenPlatform 获取微信开放平台的实例
// func (wc *Wechat) GetOpenPlatform(cfg *openConfig.Config) *openplatform.OpenPlatform {
// 	if cfg.Cache == nil {
// 		cfg.Cache = wc.cache
// 	}
// 	return openplatform.NewOpenPlatform(cfg)
// }

// // GetWork 获取企业微信的实例
// func (wc *Wechat) GetWork(cfg *workConfig.Config) *work.Work {
// 	if cfg.Cache == nil {
// 		cfg.Cache = wc.cache
// 	}
// 	return work.NewWork(cfg)
// }
