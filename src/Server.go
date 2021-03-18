package main

import (
	"TankServer/src/apis"
	"TankServer/src/core"
	"TankServer/src/defconst"
	"TankServer/zinx/ziface"
	"TankServer/zinx/znet"
)

func OnConnectionAdd(conn ziface.IConnection)  {
	// 创建一个player对象
	player := core.NewPlayer(conn)

	// 将新创建的玩家添加到worldmanager中
	core.WorldMgrObj.AddPlayer(player)

	// 将该连接绑定玩家pid 方便后续取用
	conn.SetProperty("pid", player.Pid)

	// TODO 可以做一个全局的广播 通知服务器中的全部玩家有玩家上线
}


func main() {
	s := znet.NewServer("TankServer")

	// 连接创建和销毁的钩子函数
	s.SetOnConnStart(OnConnectionAdd)

	// 注册路由 有多少从客户端接收的消息就有多少路由
	s.AddRouter(defconst.LISTROOM, &apis.ListRoom{})			// 获取房间列表
	s.AddRouter(defconst.CREATEROOM, &apis.CreateRoom{})		// 创建房间
	s.AddRouter(defconst.JOINROOM, &apis.JoinRoom{})			// 加入房间
	s.AddRouter(defconst.JOINROOM, &apis.Chat{})				// 聊天
	s.AddRouter(defconst.PLAY, &apis.Play{})					// 开始游戏
	s.AddRouter(defconst.POSITION, &apis.Move{})				// 移动
	s.AddRouter(defconst.ATTACK, &apis.Attack{})				// 开火

	// 启动服务
	s.Serve()
}
