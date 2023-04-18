package znet

import (
	"fmt"
	"strconv"
	"zinx_master/ziface"
)

/*
	消息处理模块的实现
*/

type MsgHandle struct {

	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

//创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

//调度或者执行对应的router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//从request找到msgID
	handle, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgID(), "is not found! Need Register!")
	}
	//根据msgID调度对应的router业务
	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//判断当前的Msg绑定的API处理方法是否已经存在
	//添加msg与API的绑定关系
	if _, ok := mh.Apis[msgID]; ok {
		//id已经注册
		panic("repeat api , msgID=" + strconv.Itoa(int(msgID)))
	}

	//添加msg与api的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID=", msgID, "succ!")
}
