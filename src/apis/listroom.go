package apis

import (
	"TankServer/src/core"
	"TankServer/zinx/ziface"
	"TankServer/zinx/znet"
	"fmt"
)

type ListRoom struct {
	znet.BaseRouter
}

/*
	listroom路由
	处理玩家在主界面按下展示房间列表时发送的数据
*/

func (sr *ListRoom) Handle(request ziface.IRequest)  {
	// 解析从客户端传递过来的protobuf消息

	// 调用获取房间列表方法 该方法里边应该有发消息的方法 消息内容为房间列表 即LISTROOM协议内容
	rooms := make([]*core.RoomManager, 0)
	for _, room := range core.WorldMgrObj.ListRoom() {
		rooms = append(rooms, room)
	}
	fmt.Println("rooms:", rooms)
	// 通过玩家连接属性获取到玩家pid
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("listRoom GetProperty err ,", err)
		return
	}
	// 调用玩家的同步房间列表方法
	player := core.WorldMgrObj.GetPlayerById(pid.(int32))
	player.SyncRoomList(rooms)
}