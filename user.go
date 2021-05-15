package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C chan string
	conn net.Conn
	server *Server
}

// NewUser 创建一个用户的API
func NewUser(conn net.Conn,server *Server) *User{
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:make(chan string),
		conn: conn,
		server:server,
	}

	//启动当前监听user channel的goroutine
	go user.ListenMessage()

	return user
}

// Online 用户上线的业务
func (this *User) Online()  {

	//用户上线，将用户加入OnlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播当前用户上线的消息
	this.server.BroadCast(this, "已上线")

}

// Offline 用户下线的业务
func (this *User) Offline()  {
	//用户下线，将用户从OnlineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap,this.Name)
	this.server.mapLock.Unlock()

	// 广播当前用户上线的消息
	this.server.BroadCast(this, "下线")
}

// SendMsg 将当前用户查询的结果发送给客户端
func (this *User) SendMsg(msg string)  {
	this.conn.Write([]byte(msg))
}

// DoMessage 用户处理消息的业务
func (this *User) DoMessage(msg string)  {
	if msg == "who"{
		// 查询当前用户都有哪些
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap{
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线....\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	}else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式 rename|s三
		newName := strings.Split(msg,"|")[1]

		_,ok := this.server.OnlineMap[newName]

		if ok {
			this.SendMsg("当前的用户名已经被使用")
		}else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap,this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName

			this.SendMsg("您已更新用户名")
		}
	}else {
		this.server.BroadCast(this,msg)
	}

}

// ListenMessage 监听当前user channel的方法，一旦有消息，就直接发送给客户端
func (this *User) ListenMessage()  {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
