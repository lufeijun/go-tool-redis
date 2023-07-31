# go-tool-study

https://github.com/joho/godotenv


## godotenv 流程

1、读取文件
2、解析文件
3、将值写入到 os.env 中


## 注意

1、在 load 之前开启的协程，是读不到配置信息的