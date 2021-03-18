package apis

import (
	"TankServer/src/core"
	"TankServer/src/pb"
	"TankServer/zinx/ziface"
	"TankServer/zinx/znet"
	"fmt"
	"github.com/golang/protobuf/proto"
)

/*
	move路由
	处理玩家移动时发送的数据
*/

type Move struct {
	znet.BaseRouter
}

func (m *Move) Handle(request ziface.IRequest)  {
	// 解析proto数据
	protoMsg := &pb.Position{}
	if err := proto.Unmarshal(request.GetData(), protoMsg); err != nil {
		fmt.Println("Move Handle unmarshal err ", err)
		return
	}
	// 通过连接对象获取到pid
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid err ", err)
		return
	}
	player := core.WorldMgrObj.GetPlayerById(pid.(int32))
	player.UpdatePos(protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)

}
