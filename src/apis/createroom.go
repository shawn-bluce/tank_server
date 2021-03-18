package apis

import (
	"TankServer/src/core"
	"TankServer/src/pb"
	"TankServer/zinx/ziface"
	"TankServer/zinx/znet"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type CreateRoom struct {
	znet.BaseRouter
}
/*
	createroom路由
	处理玩家点击创建房间后发送的数据
*/


func (cr *CreateRoom) Handle(request ziface.IRequest)  {
	// 解析客户端传递过来的proto协议
	protoMsg := &pb.CreatRoom{}
	err := proto.Unmarshal(request.GetData(), protoMsg)
	if err != nil {
		fmt.Println("create room unmarshal err ,", err)
		return
	}
	fmt.Printf("pid:%v, x:%v, copiesX:%v, y:%v, copiesY:%v, roomName:%v\n",
				protoMsg.Pid, protoMsg.MaxX, protoMsg.CopiesX, protoMsg.MaxY, protoMsg.CopiesY, protoMsg.RoomName)
	aoi := core.NewAOIManager(0,int(protoMsg.MaxX),int(protoMsg.CopiesX),0,int(protoMsg.MaxY),int(protoMsg.CopiesY))
	room := core.NewRoomManager(protoMsg.RoomName, protoMsg.Pid)
	core.WorldMgrObj.CreateRoom(room, aoi, protoMsg.Pid)
}