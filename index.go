package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lufeijun/go-tool-wechat/wechat"
	"github.com/lufeijun/go-tool-wechat/wechat/cache"
	"github.com/lufeijun/go-tool-wechat/wechat/officialaccount/offConfig"
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
		AppSecret: "115e73b25c9cceab35da8565c80b64b7",
		Token:     "motherfucker",
		//EncodingAESKey: "xxxx",
		Cache: redisCache,
	}

	err := redisCache.Set("name", "hello world", time.Minute)
	if err != nil {
		panic(err)
	}

	officialAccount := wc.GetOfficialAccount(cfg)

	token, err := officialAccount.GetAccessToken()
	if err != nil {
		panic(err)
	}

	fmt.Println(token)

}
