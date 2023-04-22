package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	// 链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err : ", err)
		return nil
	}

	client.conn = conn

	// 返回对象
	return client
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 更改用户名")
	fmt.Println("0. 退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>> 请输入合法数字 <<<<<")
		return false
	}
}

// DealResponse 用于异步处理服务端消息的协程，防止程序串行
func (client *Client) DealResponse() {
	// 一旦client.conn有数据，就直接拷贝到stdout标准输出上，永久阻塞监听，用于输出打印来自服务端的消息
	io.Copy(os.Stdout, client.conn) // 它是永久阻塞的
	// 等价于下面这段代码
	//for {
	//	buf := make([]byte, 1024)
	//	client.conn.Read(buf)
	//	fmt.Println(buf)
	//}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err : ", err)
		return false
	}

	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			// 公聊模式
			fmt.Println("公聊模式已选择...")
		case 2:
			// 私聊模式
			fmt.Println("私聊模式已选择...")
		case 3:
			// 更改用户名
			fmt.Println("更改用户名已选择...")
			client.UpdateName()
		case 0:
			// 退出客户端
			// 选0的时候client.flag被赋值0，不满足上面for循环的条件，所以退出
			fmt.Println("退出...")

		}
	}
}

// 定义两个全局变量
var serverIp string
var serverPort int

func init() {
	// init函数在main函数之前执行
	// flag命令行工具包
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口（默认是8888）")
}

func main() {

	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> connect failed...")
		return
	}

	// 创建成功后，单独开启一个协程用于处理server的回执消息
	go client.DealResponse()

	fmt.Println(">>>>> connect success...")

	// 阻塞，然后启动客户端业务
	// select {}
	client.Run()
}
