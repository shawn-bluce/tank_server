package znet

import (
	"TankServer/zinx/utils"
	"TankServer/zinx/ziface"
	"fmt"
	"net"
)

// iServer的接口实现，定义一个Server的服务模块
type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的ip
	IP string
	// 服务器监听的端口
	Port int
	// 当前server的消息管理模块 用来绑定MsgID和对应的处理业务API关系
	MsgHandle ziface.IMsgHandle

	// 该server的链接管理器
	ConnMgr ziface.IConnManager

	// 该Server创建链接之后自动调用的Hook函数 OnConnStart
	OnConnStart func(conn ziface.IConnection)
	// 该Server销毁链接之前自动调用的Hook函数 OnConnStop
	OnConnStop func(conn ziface.IConnection)

}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgID, router)
	fmt.Println("Add Router Success!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) Start()  {
	fmt.Printf("[Zinx] Server Name : %s, listenner at IP : %s, Port : %d is starting\n",
						utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn : %d, MaxPacketSize : %d\n",
						utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	go func() {

		// 0 开启消息队列及Worker工作池
		s.MsgHandle.StartWorkerPool()

		// 1 获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		// 2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " error ", err)
			return
		}
		fmt.Println("start Zinx server success,", s.Name, "success, Listening...")

		var cid uint32
		cid = 0

		// 3 阻塞等待客户端链接，处理客户端链接业务（读写）
		for {
			// 如果有客户端链接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 设置最大链接数的判断 如果超过最大连接数 则关闭此新的链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超出最大链接的错误包
				fmt.Println("===================================>Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				_ = conn.Close()
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定 得到我们的链接模块
			dealConnection := NewConnection(s, conn, cid, s.MsgHandle)
			cid++

			// 启动当前的链接业务处理
			go dealConnection.Start()
		}
	}()
}
func (s *Server) Stop()  {
	// TODO 将一些服务器的资源、状态或者一些已经开辟的连接信息进行停止或回收

	fmt.Println("[STOP] Zinx server name :", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve()  {
	// 启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {

	}
}
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP: utils.GlobalObject.Host,
		Port: utils.GlobalObject.TcpPort,
		MsgHandle:NewMsgHandle(),
		ConnMgr: NewConnManager(),
	}
	return s
}

// 注册 OnConnStart Hook函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}
// 注册 OnConnStop Hook函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}
// 调用 OnConnStart Hook函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("----> Call OnConnStart() ...")
		s.OnConnStart(conn)
	}
}
// 调用 OnConnStop Hook函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("----> Call OnConnStop() ...")
		s.OnConnStop(conn)
	}
}









