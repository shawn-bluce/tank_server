package apis

import (
	"TankServer/src/core"
	"TankServer/src/pb"
	"TankServer/zinx/ziface"
	"TankServer/zinx/znet"
	"fmt"
	"github.com/golang/protobuf/proto"
)

type Chat struct {
	znet.BaseRouter
}
/*
	chat路由
	处理处理世界中玩家发送聊天信息的数据
	聊天类型有 1 世界 2 房间 3 私聊
*/


func (rc *Chat) Handle(request ziface.IRequest){
	// 解析客户端传递过来的protobuf数据
	protoMsg := &pb.Chat{}
	err := proto.Unmarshal(request.GetData(), protoMsg)
	if err != nil {
		fmt.Println("Chat unmarshal err ,", err)
		return
	}
	// 根据pid得到玩家对象player
	pid, err  := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("Chat GetProperty pid err, ", err)
		return
	}
	player := core.WorldMgrObj.GetPlayerById(pid.(int32))
	if protoMsg.Tp == 3 {
		player.Chat(protoMsg.Content, int(protoMsg.Tp), int(protoMsg.Pid))
	} else {
		player.Chat(protoMsg.Content, int(protoMsg.Tp))
	}
}
