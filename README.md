# go-tool-study

1、大量的指针
  ex *string 类型，防止多次值复制？

# 对外 API 接口设计

  网关部分
    用户身份认证机制
    Endpoint
    限速
  接口定义
    接口名称
    输入参数和返回参数设计
    错误码定义
    数据加密逻辑


# 阿里 openApi 设计

1、github.com/alibabacloud-go/darabonba-openapi/v2
  为不同的服务提供统一的客户端，提供发送请求的功能
2、github.com/aliyun/credentials-go
  管理账户密码、token 等的组件

3、github.com/alibabacloud-go/cdn-20180510/v4
  CDN 相关 api 的 sdk ，即具体的阿里云服务

