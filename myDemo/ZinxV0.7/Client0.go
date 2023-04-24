package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx_master/znet"
)

func main() {

	fmt.Println("client0 start...")

	time.Sleep(1 * time.Second)
	//连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client conn error", err)
		return
	}

	for {
		//发送封包的message消息 MsgId :0
		dp := znet.NewDataPack()
		binaryMsg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx client0 Test Message")))

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("Write error", err)
			return
		}

		//服务器就应该给我们回复一个msg数据，MsgID:1 pingpingping

		//先读取流中的head部分，得到ID和dataLen

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error:", err)
			break
		}

		//将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error:", err)
		}

		if msgHead.GetMsgLen() > 0 {
			//再根据dataLen进行第二次读取，将data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {

				fmt.Println("read msg data error:", err)
				return
			}

			fmt.Println("---->Recv Server Msg : ID=", msg.Id, ",Len=", msg.DataLen, ",data=", string(msg.Data))

		}

		time.Sleep(1 * time.Second)
	}

}
