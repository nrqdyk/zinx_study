package znet

import (
	"fmt"
	"net"
	"zinx_master/ziface"
)

//当前连接模块
type Connection struct {
	//当前连接的socket conn
	Conn *net.TCPConn

	//连接的id
	ConnID uint32

	//当前的连接状态
	isClosed bool

	//告知当前连接已经停止的channel
	ExitChan chan bool

	//该链接处理的方法Router
	Router ziface.IRouter
}

//初始化连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}

	return c
}

//连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf error:", err)
			continue
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}
		//执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

//启动连接 ，让当前的连接开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID=", c.ConnID)
	//启动当前连接读数据的业务
	go c.StartReader()

	//TODO 启动从当前连接写数据的业务

}

//停止链接 结束当前连接工作
func (c *Connection) Stop() {
	fmt.Println("Conn stop()... ConnID=", c.ConnID)
	//如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//关闭socket conn
	c.Conn.Close()
	//回收资源
	close(c.ExitChan)
}

//获取当前连接所绑定的socket conn

func (c *Connection) GetTCPConnection() *net.TCPConn {

	return c.Conn
}

//获取当前连接模块的连接id
func (c *Connection) GetConnID() uint32 {

	return c.ConnID
}

//获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
