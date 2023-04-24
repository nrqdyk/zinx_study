package znet

import (
	"errors"
	"fmt"
	"io"
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

	//用于无缓冲的管道，用于读写Goroutine之间的通信
	msgChan chan []byte

	//消息的管理MsgID和对应的处理API关系
	MsgHandler ziface.IMsgHandle
}

//初始化连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
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
		/* buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf error:", err)
			continue
		} */

		//创建一个拆包解包的对象
		dp := NewDataPack()

		//读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error :", err)
			break
		}

		//拆包，得到msgID和msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error :", err)
			break
		}
		//根据dataLen再次读取data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}

		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		//从路由中，找到注册绑定的Conn对应的router调用
		//根据绑定好的MsgID找到对应吃力api业务执行
		go c.MsgHandler.DoMsgHandler(&req)
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

//提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包在发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {

	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	//将data进行封包
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msgid=", msgId)
		return errors.New("Pack error msg")
	}

	//将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id ", msgId, "error:", err)
		return errors.New("Pack error msg")
	}
	return nil
}
