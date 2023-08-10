package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/lufeijun/go-tool-wechat/wechat"
	"github.com/lufeijun/go-tool-wechat/wechat/cache"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/offConfig"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/user"
	"github.com/lufeijun/go-tool-wechat/wechat/pay/notify"
	"github.com/lufeijun/go-tool-wechat/wechat/pay/payconfig"
)

func main() {
	returncheck()
}

func returncheck() {
	// data_source: https://studygolang.com/articles/11811
	notifyStrc := &notify.Notify{Config: &payconfig.Config{Key: "ziR0QKsTUfMOuochC9RfCdmfHECorQAP"}}
	info := "YYwp8C48th0wnQzTqeI+41pflB26v+smFj9z6h9RPBgxTyZyxc+4YNEz7QEgZNWj/6rIb2MfyWMZmCc41CfjKSssoSZPXxOhUayb6KvNSZ1p6frOX1PDWzhyruXK7ouNND+gDsG4yZ0XXzsL4/pYNwLLba/71QrnkJ/BHcByk4EXnglju5DLup9pJQSnTxjomI9Rxu57m9jg5lLQFxMWXyeASZJNvof0ulnHlWJswS4OxKOkmW7VEyKyLGV6npoOm03Qsx2wkRxLsSa9gPpg4hdaReeUqh1FMbm7aWjyrVYT/MEZWg98p4GomEIYvz34XfDncTezX4bf/ZiSLXt79aE1/YTZrYfymXeCrGjlbe0rg/T2ezJHAC870u2vsVbY1/KcE2A443N+DEnAziXlBQ1AeWq3Rqk/O6/TMM0lomzgctAOiAMg+bh5+Gu1ubA9O3E+vehULydD5qx2o6i3+qA9ORbH415NyRrQdeFq5vmCiRikp5xYptWiGZA0tkoaLKMPQ4ndE5gWHqiBbGPfULZWokI+QjjhhBmwgbd6J0VqpRorwOuzC/BHdkP72DCdNcm7IDUpggnzBIy0+seWIkcHEryKjge3YDHpJeQCqrAH0CgxXHDt1xtbQbST1VqFyuhPhUjDXMXrknrGPN/oE1t0rLRq+78cI+k8xe5E6seeUXQsEe8r3358mpcDYSmXWSXVZxK6er9EF98APqHwcndyEJD2YyCh/mMVhERuX+7kjlRXSiNUWa/Cv/XAKFQuvUYA5ea2eYWtPRHa4DpyuF1SNsaqVKfgqKXZrJHfAgslVpSVqUpX4zkKszHF4kwMZO3M7J1P94Mxa7Tm9mTOJePOoHPXeEB+m9rX6pSfoi3mJDQ5inJ+Vc4gOkg/Wd/lqiy6TTyP/dHDN6/v+AuJx5AXBo/2NDD3fWhHjkqEKIuARr2ClZt9ZRQO4HkXdZo7CN06sGCHk48Tg8PmxnxKcMZm7Aoquv5yMIM2gWSWIRJhwJ8cUpafIHc+GesDlbF6Zbt+/KXkafJAQq2RklEN+WvZ/zFz113EPgWPjp16TwBoziq96MMekvWKY/vdhjol8VFtGH9F61Oy1Xwf6DJtPw=="
	res, err := notifyStrc.DecryptReqInfo(&notify.RefundedResult{ReqInfo: &info})
	if err != nil {
		fmt.Println("解密失败：", err)
		return
	}

	bytes, err := xml.Marshal(res)
	if err != nil {
		fmt.Println("解密xml失败：", err)
		return
	}
	fmt.Println(string(bytes))
}

// 简单的获取 token 测试
func wechatToken() {
	wc := wechat.NewWechat()

	redisOpts := &cache.RedisOpts{
		Host:            "192.168.0.33:6379",
		Password:        "123456",
		Database:        10,
		MaxActive:       10,
		MaxIdle:         10,
		ConnMaxIdleTime: 60, //second
	}
	redisCache := cache.NewRedis(context.Background(), redisOpts)

	cfg := &offConfig.Config{
		AppID:     "wx82499d23a58480fa",
		AppSecret: "115e73b25c9cc111111111eab35da8565c80b64b7",
		Token:     "motherfucker",
		//EncodingAESKey: "xxxx",
		Cache: redisCache,
	}

	err := redisCache.Set("name", "hello world", time.Minute)
	if err != nil {
		panic(err)
	}

	officialAccount := wc.GetOfficialAccount(cfg)
	// token, err := officialAccount.GetAccessToken()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(token)
	// fmt.Println("=====================")

	// basic := officialAccount.GetBasic()
	// // fmt.Println(basic.GetAPIDomainIP())
	// short, err := basic.Long2ShortURL("https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/qrcode/shorturl.html#%E8%AF%B7%E6%B1%82%E5%9C%B0%E5%9D%80")
	// if err != nil {
	// 	fmt.Println("short err：", err)
	// 	return
	// }
	// fmt.Println(short)

	// material := officialAccount.GetMaterial()
	// rmc, err2 := material.GetMaterialCount()
	// if err2 != nil {
	// 	fmt.Println("short err：", err2)
	// 	return
	// }
	// fmt.Println(rmc)

	us := officialAccount.GetUser()
	users, err3 := us.BatchGetUserInfo(user.BatchGetUserInfoParams{})
	if err3 != nil {
		fmt.Println("short err：", err3)
		return
	}
	fmt.Println(users)
}
