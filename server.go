package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// Server 创建server对象
type Server struct {
	Ip string
	Port int

	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// NewServer 创建一个server的接口
func NewServer(ip string,port int) *Server  {
	server := &Server{
		Ip:ip,
		Port:port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}

	return server
}

// ListenMessage 监听Message广播消息的goroutine,一旦有消息就发送给全部在user
func (this *Server) ListenMessage()  {
	for {
		msg := <-this.Message

		//将msg发送给全部的在线用户
		this.mapLock.Lock()
		for _,cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}
// BroadCast 广播消息的方法
func (this *Server) BroadCast(user *User,msg string)  {
	sendMsg := "{" + user.Addr + "}" + user.Name + ":" + msg
	this.Message <- sendMsg // 这里会阻塞
}


func (this *Server) Handler(conn net.Conn)  {
	//....当前链接的业务
	fmt.Println("链接建立成功",conn)

	user := NewUser(conn,this)

	user.Online()
	/*
	//用户上线，将用户加入OnlineMap中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播当前用户上线的消息
	this.BroadCast(user, "已上线")
	 */

	// 接受客户端发送的消息
	go func() {
		for {
			buf := make([]byte, 4096) // byte 为uint
			cnt, err := conn.Read(buf)
			if cnt == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				// != io.EOF,是因为读操作末端都存在io.EOF标识，如果没有，则用户进行了非法操作
				fmt.Printf("Conn Read err:", err)
				return
			}

			// 提取用户的消息
			msg := string(buf[:cnt-1])
			user.DoMessage(msg)
		}
	}()

}

// Start 启动服务的接口
func (this *Server) Start() {
	// socket listen
	listener,err := net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Port))

	if err != nil {
		fmt.Println("net.Listen err:",err)
		return
	}

	// close listen socket
	defer listener.Close()

	// 启动监听Message的goroutine
	go this.ListenMessage()

	//accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:",err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}
}

