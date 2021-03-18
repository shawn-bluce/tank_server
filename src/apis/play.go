package apis

import (
	"TankServer/src/core"
	"TankServer/src/pb"
	"TankServer/zinx/ziface"
	"TankServer/zinx/znet"
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
)
/*
	play路由
	处理在房间准备完毕后 房主点击play按钮开始游戏发送的数据
 */

type Play struct {
	znet.BaseRouter
}

func (p *Play) Handle(request ziface.IRequest)  {
	//解析协议信息
	protoMsg := &pb.Play{}
	if err := proto.Unmarshal(request.GetData(), protoMsg); err != nil {
		fmt.Println("Play Handle unmarshal err ", err)
		return
	}

	// 生成玩家列表数量的坐标 并将玩家放入aoi地图中
	for _, pid := range protoMsg.Pid {
		player := core.WorldMgrObj.GetPlayerById(pid)
		aoi := core.WorldMgrObj.AoiMgr[protoMsg.RoomName]
		// 随机玩家坐标
		player.X = float32(rand.Intn(aoi.MaxX))
		player.Z = float32(rand.Intn(aoi.MaxY))
		// 将玩家添加到格子中
		aoi.AddToGridByPos(int(pid), player.X, player.Z)
		// 广播玩家位置信息
		//player.BroadCastStartPosition()
		// 同步周边玩家位置信息
		player.SyncSurrounding()
	}
}
