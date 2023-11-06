package main

import (
	"fmt"

	"github.com/lufeijun/go-tool-aliyun/common"
	// cdn "github.com/lufeijun/go-tool-aliyun/cdn"
	// "github.com/lufeijun/go-tool-aliyun/openapi"
	// "github.com/lufeijun/go-tool-aliyun/tool"
)

type Person struct {
	Name string
	Age  int
	Like []string
}

func main() {
	// client, _err := CreateClient(
	// 	tool.String("LTAI5tCi7UZEr4FKouV4pWCS"),
	// 	tool.String("ImxIl5N53xPFY4Vt2wlvByeiGxTD2A"),
	// )
	// if _err != nil {
	// 	fmt.Println("客户端连接失败：", _err)
	// }

	// request := &cdn.RefreshObjectCachesRequest{
	// 	ObjectPath: tool.String("https://mallassets.zhufaner.com/6e141848-9abc-11ec-bac9-00163e352f67.jpg"),
	// }

	// _, err := client.RefreshObjectCaches(request)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	p := Person{
		Name: "张三",
		Age:  123,
		Like: []string{
			"篮球",
			"足球",
		},
	}

	fmt.Println(p)
	fmt.Println(common.Prettify(p))

}

// func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *cdn.Client, _err error) {
// 	config := &openapi.Config{
// 		// 必填，您的 AccessKey ID
// 		AccessKeyId: accessKeyId,
// 		// 必填，您的 AccessKey Secret
// 		AccessKeySecret: accessKeySecret,
// 	}
// 	// Endpoint 请参考 https://api.aliyun.com/product/Cdn
// 	config.Endpoint = tool.String("cdn.aliyuncs.com")
// 	// config.RegionId = tea.String("cn-beijing")
// 	_result = &cdn.Client{}
// 	_result, _err = cdn.NewClient(config)
// 	return _result, _err
// }
