package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	fmt.Println("client start...")

	time.Sleep(1 * time.Second)
	//连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client conn error", err)
		return
	}

	for {
		//连接调用write方法写数据
		_, err := conn.Write([]byte("Hello ZinxV0.4"))
		if err != nil {
			fmt.Println("conn write error:", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error:", err)
			return
		}
		fmt.Printf("server call back: %s,cnt = %d\n", buf, cnt)

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}

}
