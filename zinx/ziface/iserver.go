package ziface

// 定义一个服务器接口
type IServer interface {
	// 启动服务器方法
	Start()
	// 停止服务器方法
	Stop()
	// 运行服务器方法
	Serve()
	// 路由功能：给当前的服务注册一个路由方法，供客户端的链接处理
	AddRouter(msgID uint32, router IRouter)

	// 获取当前server的链接管理器
	GetConnMgr() IConnManager

	// 注册 OnConnStart Hook函数的方法
	SetOnConnStart(func(connection IConnection))
	// 注册 OnConnStop Hook函数的方法
	SetOnConnStop(func(connection IConnection))
	// 调用 OnConnStart Hook函数的方法
	CallOnConnStart(conn IConnection)
	// 调用 OnConnStop Hook函数的方法
	CallOnConnStop(conn IConnection)

}