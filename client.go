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
	flag       int // 当前client的模式
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
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn
	return client
}

var serverIp string
var serverPort int

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>> 请输入合法范围内的数字")
		return false
	}
}

// 更新用户名
func (Client *Client) UpdateName() bool {
	fmt.Println(">>>>> 请输入用户名:")
	fmt.Scanln(&Client.Name)
	sendMsg := "rename|" + Client.Name + "\n"
	_, err := Client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

// 处理server回应的消息，直接显示在标准输出即可
func (client *Client) DealResponse() {
	// 一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			fmt.Println(">>>>> 公聊模式")
			break
		case 2:
			fmt.Println(">>>>> 私聊模式")
			break
		case 3:
			// 更新用户名
			client.UpdateName()
			break
		}
	}
}

// 命令行解析
func init() {
	// flg type name default usage
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.01)")
	flag.IntVar(&serverPort, "port", 7777, "设置服务器端口(默认是8888)")
}

func main() {
	// 命令行解析
	flag.Parse()
	// 创建客户端
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> server not start")
		return
	}
	// 开启一个goroutine处理server回执的消息
	go client.DealResponse()
	fmt.Println(">>>>> server start success")
	// select {}
	client.Run()
}
