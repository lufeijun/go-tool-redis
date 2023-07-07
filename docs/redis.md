# REDIS 部分说明

# RESP 协议

目前新版的是  3

参考：
* https://blog.csdn.net/LZH984294471/article/details/114233835
* https://www.jianshu.com/p/f670dfc9409b

# redis 操作基础逻辑

1、进行 tcp 连接
2、进行用户密码验证
  resp 中，通过 RESP "hello" "3" "auth" "default" "123456" 进行验证
3、按照 resp 协议，组织要发送的命令
  ex：set name tom  == *3\r\n$3\r\nSET\r\n$4\r\nname\r\n$3\r\ntom\r\n
4、通过网络 IO 流 发送 redis 命令
5、通过网络 IO 流 接受 redis 的返回值，并按照 resp 协议进行解析
6、最后，关闭 tcp 连接

```
18:05:01.710606 IP 192.168.0.92.55743 > jipeng87.6379: Flags [P.], seq 1:58, ack 1, win 2058, options [nop,nop,TS val 498569807 ecr 2798473668], length 57: RESP "hello" "3" "auth" "default" "123456"
18:05:01.711070 IP jipeng87.6379 > 192.168.0.92.55743: Flags [P.], seq 1:149, ack 58, win 227, options [nop,nop,TS val 2798473669 ecr 498569807], length 148: RESP "%7" "server" "redis" "version" "7.0.9" "proto" "3" "id" "27644" "mode" "standalone" "role" "master" "modules" empty
```


# 说明

1、所有的封装，都是基于 redis 基础连接操作，进行抽象
  连接池：基于 tcp 连接的管理
  钩子函数：基于向 tcp 发送命令时，分别执行相关函数
    tcp 连接部分：创建新连接时的 hook
    发送 command：发送命令时的 hook