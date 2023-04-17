package utils

import (
	"encoding/json"
	"os"
	"zinx_master/ziface"
)

/*
	存储一切有关zinx框架的全局参数，供其他模块使用
	一些参数是可以通过zinx.json由用户进行配置

*/

type GlobalObj struct {

	//当前Zinx全局的server
	TcpServer ziface.IServer
	Host      string //当前主机监听的ip
	TcpPort   int    //当前主机监听的端口号
	Name      string //当前服务器名称

	Version        string //当前zinx的版本号
	MaxConn        int    //当前服务器主机允许的最大连接数
	MaxPackageSize uint32 //当前zinx框架数据包的最大值
}

//定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

//从zinx.json去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

//提供一个init方法，初始化当前的GlobalObject
func init() {
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.5",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	//应该尝试从conf/zinx.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}
