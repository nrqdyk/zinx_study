package znet

import (
	"fmt"
	"net"
	"zinx_master/ziface"
)

//iserver的接口实现，定义一个server的服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP :%s , Port %d , is starting \n", s.IP, s.Port)

	go func() {
		//获取一个tcp的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}

		//监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Listen tcp error:", err)
			return
		}
		fmt.Println("start zinx server success", s.Name, "successs,Listenning...")

		//阻塞等待客户端连接，处理客户端连接业务
		for {
			//如果有客户端连接过来阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept tcp error:", err)
				continue
			}

			//已经和客户端建立连接，做一些业务，做一个最基本的回显业务

			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("conn read error:", err)
						continue
					}

					fmt.Printf("recv client buf %s , cnt %d \n", buf, cnt)
					//回显功能
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()
		}
	}()

}

func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态或者一些u已经开辟的连接信息，进行停止或者回收
}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//应该处于阻塞状态

	select {}
}

//初始化Server
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}

	return s
}
