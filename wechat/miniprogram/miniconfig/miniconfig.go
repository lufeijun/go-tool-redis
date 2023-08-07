package miniconfig

import "github.com/lufeijun/go-tool-wechat/wechat/cache"

// Config .config for 小程序
type Config struct {
	AppID     string `json:"app_id"`     // appid
	AppSecret string `json:"app_secret"` // appSecret
	Cache     cache.Cache
}
