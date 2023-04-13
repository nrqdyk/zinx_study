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

// Test PreHandle
func (this *PingRouter) PreHandle(requset ziface.IRequest) {
	fmt.Println("Call Router PreHandle....")
	_, err := requset.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("Call back before ping error")
	}

}

//Test Handle
func (this *PingRouter) Handle(requset ziface.IRequest) {
	fmt.Println("Call Router Handle....")
	_, err := requset.GetConnection().GetTCPConnection().Write([]byte("ping ...ping... ping..."))
	if err != nil {
		fmt.Println("Call back  ping error")
	}
}

//Test PostHandle
func (this *PingRouter) PostHandle(requset ziface.IRequest) {
	fmt.Println("Call Router PostHandle....")
	_, err := requset.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
	if err != nil {
		fmt.Println("Call back after ping error")
	}
}

func main() {
	//创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinx V0.3]")

	//给当前的zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	//启动server

	s.Serve()
}
