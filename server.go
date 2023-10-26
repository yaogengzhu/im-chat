package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}

func (this *Server) Handler(conn net.Conn) {
	// do something
	fmt.Println("connect success")
}

func (this *Server) Start() {
	// do something
	listern, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))

	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	defer listern.Close()

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
