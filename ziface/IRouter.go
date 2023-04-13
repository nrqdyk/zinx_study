package ziface

/*
	路由抽象接口，
	路由里的数据都是IRuquest
*/

type IRouter interface {
	//在处理conn业务之前的钩子方法Hook
	PreHandle(request IRequest)
	//在处理conn业务的主方法hook
	Handle(request IRequest)
	//在处理conn业务之后的钩子方法hook
	PostHandle(request IRequest)
}
