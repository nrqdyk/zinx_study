package main

import (
	"fmt"
	"zinx_master/ziface"
	"zinx_master/znet"
)

//ping test 自定义路由

type PingRouter struct {
	znet.BaseRouter
}

//Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	//回写数据
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinx V0.5]")

	//给当前的zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	//启动server

	s.Serve()
}
