package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 当前在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

func (this *Server) Handler(conn net.Conn) {
	// ...当前链接的业务
	fmt.Println("链接建立成功！")

	// 链接成功后，就创建一个user
	user := NewUser(conn)

	// 用户上线，将用户添加到OnlineMap中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	// 广播当前用户上线信息，那就需要写一个广播方法
	this.BroadCast(user, "is online...")

	// 当前Handler阻塞
	select {}

}

// 广播消息方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

// 启用一个协程监听Message消息，一旦有消息就发送给全部在线User
func (this *Server) MessageListener() {
	for {
		msg := <-this.Message
		// 将msg发送给全部在线用户
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 启动服务器的接口
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Printf("net.Listen er : %v", err)
		return
	}

	// close listen socket
	defer listener.Close()

	// 启动监听Message的协程
	go this.MessageListener()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err : ", err)
			continue
		}

		// do handler启动协程处理handler
		go this.Handler(conn)
	}

}
