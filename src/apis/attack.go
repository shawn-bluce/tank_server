package apis

import (
	"TankServer/zinx/ziface"
	"TankServer/zinx/znet"
)

/*
	attack路由
	处理在房间中玩家点击开火按钮发送的数据
*/

type Attack struct {
	znet.BaseRouter
}

func (f *Attack) Handle(request ziface.IRequest)  {
	// 开始战斗会遍历房间中的pid 然后将玩家添加到aoi地图中
	// 即 开始战斗时才会分配玩家之间的位置信息 所以移除时只需要将aoi room 管理器中对应房间移除 并且将玩家身上保存的房间信息移除即可
	// TODO 开始战斗时需要修改房间内玩家的坐标信息 并广播每个玩家出生地点

}
