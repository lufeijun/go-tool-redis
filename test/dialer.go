package test

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

func aaa(conn net.Conn) {
	fmt.Println("defer")
	conn.Close()
}

func Out() {
	conn, err := NewDialer()
	if err != nil {
		fmt.Println("连接失败", err)
		return
	}
	defer aaa(conn)

	// 绑定 io 读取|写入
	rd := NewReader(conn)
	bw := bufio.NewWriter(conn)
	wr := NewWriter(bw)

	// 构造参数
	protocol := 3
	pwd := "123456"

	// 组织参数
	args := make([]interface{}, 0, 7)
	args = append(args, "hello", protocol)
	args = append(args, "auth", "default", pwd)

	// 发送命令
	wr.WriteArgs(args)
	bw.Flush()

	// 读取
	i, err := rd.ReadReply()
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "WRONGPASS") {
			panic("密码错误")
		}
	}

	flag := false
	switch i.(type) {
	case string:
		flag = true
		fmt.Println("读取数据", i.(string))
		if strings.Contains(i.(string), "WRONGPASS") {
			panic("密码错误")
		}
		break
	}

	if flag {
		m := i.(map[interface{}]interface{})
		fmt.Println("读取数据", m)
		version := m["version"]
		fmt.Println("读取数据", version)
	}

	// time.Sleep(15 * time.Second)

	// wr.WriteArgs(args)
	// bw.Flush()

	fmt.Println("out over")
}

// tcp 连接
func NewDialer() (net.Conn, error) {

	ctx := context.Background()

	network := "tcp"
	addr := "192.168.0.87:6379"

	// net.Dialer go 官方的方法
	netDialer := &net.Dialer{
		Timeout: 5 * time.Second,
		// KeepAlive: 10 * time.Second,
	}
	return netDialer.DialContext(ctx, network, addr)
}

func Hello(p *net.Conn) {

	protocol := 3
	pwd := "123456"

	// 组织参数
	args := make([]interface{}, 0, 7)
	args = append(args, "hello", protocol)
	args = append(args, "auth", "default", pwd)

	bw := bufio.NewWriter(*p)

	// p.Write(args)

	bw.Write([]byte("hello"))

	bw.Flush()

	// 发送命令
	fmt.Println(args)

}
