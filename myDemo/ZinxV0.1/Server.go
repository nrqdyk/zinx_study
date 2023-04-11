package main

import "zinx_master/znet"

func main() {
	//创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinx V0.1]")

	//启动server

	s.Serve()
}
