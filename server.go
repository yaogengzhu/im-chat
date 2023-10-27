package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线user
func (this *Server) ListenMessager() {
	for {
		// 尝试读取Message的数据
		msg := <-this.Message

		// 将msg发送给全部的在线user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()

	}
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg

}

func (this *Server) Handler(conn net.Conn) {
	// do something
	// fmt.Println("connect success")

	// 用户上线，将用户加入到onlinemap中
	user := NewUser(conn, this)

	user.UserOnline()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 { // 对端断开，或者出问题
				user.UserOffline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read err:", err)
				return
			}
			// 自定义消息广播格式：ip:port:msg
			msg := string(buf[:n-1]) // windows上多一个换行符
			// 将得到的消息进行广播
			user.DoMessage(msg)
			isLive <- true
		}
	}()
	// 当前handler阻塞

	for {

		select {
		case <-isLive:
			// 当前用户是活跃的，应该重置定时器
			// 不做任何事情，为了激活select，更新下面的定时器
		case <-time.After(time.Second * 10):
			// fmt.Println("超时")
			user.SendMsg("你被踢了\n")
			// 销毁用的资源
			close(user.C)
			conn.Close()
		}
	}
}

func (this *Server) Start() {
	// do something
	listern, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	defer listern.Close()

	// 启动监听Message的goroutine
	go this.ListenMessager()

	for {
		conn, err := listern.Accept()
		if err != nil {
			fmt.Println("listern.Accept err:", err)
			continue
		}

		// handle
		go this.Handler(conn)

	}
}
