package core

import (
	"fmt"
	"sync"
)

/*
	游戏的世界管理模块
 */

// 属性
type WorldManager struct {
	// aoi管理模块map key:房间名 value:aoi管理模块
	AoiMgr map[string]*AOIManager
	// 玩家map key:pid value:玩家对象
	Players map[int32]*Player
	// 保护玩家的锁
	pLock sync.RWMutex
	// 房间管理模块map key:房间名 value:房间管理模块
	RoomMgr map[string]*RoomManager
}

// 提供对外的世界管理模块句柄
var WorldMgrObj *WorldManager


// aoi管理模块 玩家 房间管理模块 集合的初始化
func init()  {
	WorldMgrObj = &WorldManager{
		AoiMgr: make(map[string]*AOIManager),
		Players: make(map[int32]*Player),
		RoomMgr: make(map[string]*RoomManager),
	}
}

// 创建房间
// 创建房间之前必须先对aoi管理模块和房间管理模块进行初始化 即调用New方法
func (wm *WorldManager) CreateRoom(roomManager *RoomManager,aoiManager *AOIManager, pID int32) {
	// 给玩家设置地图名称信息
	// 加入房间列表
	wm.RoomMgr[roomManager.name] = roomManager
	wm.RoomMgr[roomManager.name].pIDList = append(wm.RoomMgr[roomManager.name].pIDList, pID)

	// 加入aoi列表
	wm.AoiMgr[roomManager.name] = aoiManager
}

// 加入房间
func (wm *WorldManager) JoinRoom(pID int32, roomName string) {
	// 将玩家pID加入房间列表中的玩家队列
	wm.RoomMgr[roomName].pIDList = append(wm.RoomMgr[roomName].pIDList,pID)
	// 设置玩家的房间名称
	wm.Players[pID].roomName = roomName
}

// 列出房间列表
func (wm *WorldManager) ListRoom() (roomList map[string]*RoomManager) {
	if len(wm.RoomMgr) == 0 {
		fmt.Println("roomList is empty, please create room!")
		return
	}
	return WorldMgrObj.RoomMgr
}

// 根据房间名称得到房间对象
func (wm *WorldManager) GetRoomMgrByName(name string) *RoomManager  {
	return wm.RoomMgr[name]
}

// 删除房间
// 由于房间和aoi是一一对应的 不同房间对应着不同aoi
// 此方法在战斗结束之后调用 移除内容包括 aoi地图 房间 player对象中的房间信息
func (wm *WorldManager) RemoveRoom(roomName string)  {
	// roomName 房间名称
	// 移除玩家身上的房间信息 遍历房间信息中的玩家列表
	room := wm.GetRoomMgrByName(roomName)
	for _, v := range room.pIDList {
		player, ok := wm.Players[v]
		if !ok {
			fmt.Println("player is not exist on this room!")
			return
		}
		// 将房间名称设置为空
		player.roomName = ""
	}
	// 移除aoi地图信息 根据房间名称
	delete(wm.AoiMgr, roomName)
	// 移除房间信息 根据房间名称
	delete(wm.RoomMgr, roomName)

	// TODO 移除房间之后可能需要将战绩保留 即 将房间信息 地图信息 以及获胜pid/玩家设置名称 记录到库中
}

// 此处为玩家连接创建时调用 将玩家加入到玩家列表中 战斗中维护的player列表为房间信息中的pid切片
func (wm *WorldManager) AddPlayer(player *Player)  {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	defer wm.pLock.Unlock()
}

// 删除玩家 为玩家断开连接时调用 将玩家从世界管理器中移除
func (wm *WorldManager) RemovePlayerInWorld(pid int32){
	wm.pLock.Lock()
	defer wm.pLock.Unlock()
	delete(wm.Players, pid)
}

// 获取世界中的所有玩家
func (wm *WorldManager) GetAllPlayersInWorld() []*Player  {
	players := make([]*Player, 0)
	for _, p := range wm.Players {
		players = append(players, p)
	}
	return players
}

// 获取房间内的所有玩家
func (wm *WorldManager) GetAllPlayersInRoom(roomName string) []*Player{
	wm.pLock.Lock()
	defer wm.pLock.Unlock()
	players := make([]*Player, 0)
	for _, id := range WorldMgrObj.RoomMgr[roomName].pIDList {
		players = append(players, wm.GetPlayerById(id))
	}
	return players
}

// 通过玩家id查询player对象
func (wm *WorldManager) GetPlayerById(pid int32) *Player {
	wm.pLock.Lock()
	defer wm.pLock.Unlock()
	return wm.Players[pid]
}

