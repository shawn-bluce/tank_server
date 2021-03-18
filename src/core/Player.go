package core

import (
	"TankServer/src/defconst"
	"TankServer/src/pb"
	"TankServer/zinx/ziface"
	"fmt"
	"github.com/golang/protobuf/proto"
	"sync"
)

type Player struct {
	// 玩家ID
	Pid int32
	// 当前玩家的连接 用于和客户端的连接 比如需要给玩家发送消息
	Conn ziface.IConnection
	// 平面X坐标
	X float32
	// 高度
	Y float32
	// 平面Y坐标
	Z float32
	// 旋转角度 0~360
	V float32
	// 房间名称
	roomName string
	// TODO:坦克对象，包括坦克的攻击移动视野范围攻击距离等信息
}

/*
	Player ID 生成器
 */
var PidGen int32 = 1
var IdLock sync.Mutex

// 初始化一个玩家 在连接的时候调用
func NewPlayer(conn ziface.IConnection) *Player{
	// 生成一个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()
	p := &Player{
		Pid: id,
		Conn: conn,
		X: 0,
		Y: 0,
		Z: 0,
		V: 0,
		roomName: "",
		// TODO:使用配置文件初始化坦克 / 读取玩家信息获取玩家的坦克
	}
	return p
}


func (p *Player) SendMsg(msgId uint32, data proto.Message)  {
	// 将 proto Message 结构体序列化 转换成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("SendMsg Marshal data err ", err)
		return
	}
	// 将二进制文件通过zinx框架的SendMsg将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("target client is closed or miss")
		return
	}
	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("SendMsg error ,", err)
		return
	}
	return
}

// 广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition()  {
	// 组建proto数据
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp: 1,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 发送消息给客户端
	p.SendMsg(defconst.BROADCAST, protoMsg)
}
// 玩家广播世界聊天消息

// 同步玩家上线的位置消息
func (p *Player) SyncSurrounding()  {
	// 获取玩家周围格子玩家的pid
	pids := WorldMgrObj.AoiMgr[p.roomName].GetPidByPos(p.X, p.Z)
	players := make([]*Player, 0)
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerById(int32(pid)))
	}
	// 将当前玩家位置信息发送给周边玩家 让周围玩家能看到你
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp: 1,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	for _, player := range players{
		player.SendMsg(defconst.BROADCAST, protoMsg)
	}
	// 将周围玩家信息发送给你 让你能看到周围玩家
	playersProtoMsg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersProtoMsg = append(playersProtoMsg, p)
	}
	// 封装proto数据
	SyncPlayers := &pb.SyncPlayers{
		Ps: playersProtoMsg,
	}
	p.SendMsg(defconst.SYNCPLAYERS, SyncPlayers)
}


// 广播当前玩家的位置移动信息
func (p *Player) UpdatePos(x, y, z, v float32)  {
	// 更新当前玩家player对象的坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v
	// 组建广播proto协议 MsgID:7 Tp = 3
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp: 2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: x,
				Y: y,
				Z: z,
				V: v,
			},
		},
	}
	// 遍历aoi中的全部玩家
	players := p.GetSurroundingPlayers()
	for _, player := range players {
		player.SendMsg(defconst.BROADCAST, protoMsg)
	}
}

// 获取当前玩家的周边玩家AOI九宫格之内的玩家
func (p *Player) GetSurroundingPlayers() []*Player {
	pids := WorldMgrObj.AoiMgr[p.roomName].GetPidByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerById(int32(pid)))
	}
	return players
}

// 玩家下线业务

// 玩家同步房间列表
func (p *Player) SyncRoomList(rooms []*RoomManager)  {
	// 组建同步房间列表的proto协议 MsgID:10
	// 先组建出一个协议中的room列表
	roomProtoMsg := make([]*pb.Room, 0)
	for _, room := range rooms {
		r := &pb.Room{
			Pid: room.ownerPid,
			RoomName: room.name,
		}
		roomProtoMsg = append(roomProtoMsg, r)
	}
	SyncRoomListProtoMsg := &pb.SyncRoomList{
		Room: roomProtoMsg,
	}
	// TODO:后续在房间列表信息中添加地图信息 例如大小或者更详细的信息（比如buf多少之类的可扩展内容）等...
	p.SendMsg(defconst.SYNCROOMLIST, SyncRoomListProtoMsg)
}

// 聊天
func (p *Player) Chat(content string, data ...int) {
	// 组建chat的proto协议 MsgID:12
	chatProtoMsg := &pb.BroadCastChat{
		Pid: p.Pid,
		Content: content,
	}
	players := make([]*Player,0)
	switch data[0] {
	case 1:
		players = WorldMgrObj.GetAllPlayersInWorld()
		break
	case 2:
		players = WorldMgrObj.GetAllPlayersInRoom(p.roomName)
		break
	case 3:
		// 根据玩家pid得到玩家对象
		chatTarget := WorldMgrObj.GetPlayerById(int32(data[1]))
		players = append(players, chatTarget)
		break
	}
	for _, player := range players {
		player.SendMsg(defconst.BROADCASTCHAT, chatProtoMsg)
	}
}