package znet

import (
	"fmt"
	"net"
	"zinx_master/utils"
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
	//当前server消息管理模块，用来绑定msgid和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
}

func (s *Server) Start() {

	fmt.Printf("[Zinx] Server Name : %s,listenner at IP:%s , Port:%d is starting",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

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
			var cid uint32
			cid = 0
			//将处理新连接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			//启动当前的连接业务处理
			go dealConn.Start()
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

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!!!")
}

//初始化Server
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}

	return s
}
