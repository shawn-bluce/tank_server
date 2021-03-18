package apis

import (
	"TankServer/src/core"
	"TankServer/src/pb"
	"TankServer/zinx/ziface"
	"TankServer/zinx/znet"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type JoinRoom struct {
	znet.BaseRouter
}

/*
	joinroom路由
	处理在玩家点击加入房间后发送过来的数据
*/


func (jr *JoinRoom) Handle(request ziface.IRequest)  {
	// 获取传递过来的信息
	protoMsg := &pb.JoinRoom{}
	err := proto.Unmarshal(request.GetData(), protoMsg)
	if err != nil {
		fmt.Println("JoinRoom unmarshal protoMsg err ,", err)
		return
	}
	// 调用世界管理器中加入房间方法
	core.WorldMgrObj.JoinRoom(protoMsg.Pid, protoMsg.RoomName)
}
