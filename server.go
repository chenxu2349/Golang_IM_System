package main

import (
	"fmt"
	"io"
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

// NewServer 创建一个server的接口
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
	user := NewUser(conn, this)

	//// 用户上线，将用户添加到OnlineMap中
	//this.mapLock.Lock()
	//this.OnlineMap[user.Name] = user
	//this.mapLock.Unlock()
	user.Online()

	// 广播当前用户上线信息，那就需要写一个广播方法
	this.BroadCast(user, "**********is online...**********")

	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				//this.BroadCast(user, "**********log out...**********")
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn read err : ", err)
				return
			}

			// 提取用户的消息（去除\n）
			msg := string(buf[:n-1])

			//// 将得到的消息进行广播
			//this.BroadCast(user, msg)
			user.DoMessage(msg)
		}
	}()

	// 当前Handler阻塞
	select {}

}

// BroadCast 广播消息方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

// MessageListener 启用一个协程监听Message消息，一旦有消息就发送给全部在线User
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

// Start 启动服务器的接口
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
