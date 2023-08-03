package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lufeijun/go-tool-wechat/wechat"
	"github.com/lufeijun/go-tool-wechat/wechat/cache"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/offConfig"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/user"
)

func main() {
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
