package GoChat

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

//NewServer creates a server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

//Handler handler connection
func (s *Server) Handler(conn net.Conn) {
	//do something
	fmt.Println("连接成功！")
}

//StartServer starts a server
func (s *Server) StartServer() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen error: ", err)
		return
	}
	//close socket listener
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("net.listener error2: ", err)
		}
	}(listener)

	for {
		//accept connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("net.Accept error: ", err)
			continue
		}
		go s.Handler(conn)
	}
}
