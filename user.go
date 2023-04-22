package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// NewUser 创建一个用户的接口
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// 启动监听当前User channel的协程
	go user.ListenMessage()

	return user
}

// Online 用户上线业务
func (this *User) Online() {

	// 用户上线，将用户添加到OnlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播当前用户上线信息，那就需要写一个广播方法
	this.server.BroadCast(this, ">>>>> is Online... ")
}

// Offline 用户下线业务
func (this *User) Offline() {

	// 用户下线，将用户从OnlineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 广播当前用户上线信息，那就需要写一个广播方法
	this.server.BroadCast(this, ">>>>> is Offline... ")
}

// SendMsg 向当前用户的客户端发消息
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// DoMessage 用户处理消息业务
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "is online...\n"
			fmt.Println("上线！")
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 消息格式: to|张三|hello

		// 1、获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("消息格式不正确，请使用'to|张三|hello '格式\n")
			return
		}

		// 2、根据用户名获取对方User对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("该用户名不存在\n")
			return
		}

		// 3、获取消息内容，并用过对方User对象发送过去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("消息为空，请重发！\n")
			return
		}
		remoteUser.SendMsg(this.Name + " : " + content)

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 更改用户名消息格式为：rename|张三
		newName := msg[8:]

		// 判断新名字是否已经存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名已存在！")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg("您已更新用户名：" + this.Name + "\n")
		}

	} else {
		this.server.BroadCast(this, msg)
	}
}

// ListenMessage 监听当前User channel的方法，一旦有消息就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
