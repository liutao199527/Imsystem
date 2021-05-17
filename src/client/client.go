package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // 当前client的模式
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// 链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
	}

	client.conn = conn

	// 返回对象
	return client
}

//解析命令行
var serverIp string
var serverPort int

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.public")
	fmt.Println("2.private")
	fmt.Println("3.update")
	fmt.Println("0.exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>>>请输入合法范围内的数字<<<<<<")
		return false
	}
}

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "do IP,default:127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "do port.default:8888")
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>>>>>>>connect server failed........")
		return
	}

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>connect server success.......")

	// 启动客户端业务
	select {}
}
