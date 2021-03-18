package core

import "sync"

type RoomManager struct {
	// 房间ID
	rID int32
	// 房间名
	name string
	// 玩家列表id
	pIDList []int32
	// 房主pid
	ownerPid int32

}

/*
	roomID生成器
 */
var rIDGen int32 = 1
var rIDLock sync.Mutex

func NewRoomManager(name string, ownerPid int32) *RoomManager {
	// 生成一个房间ID
	rIDLock.Lock()
	id := rIDGen
	rIDGen ++
	rIDLock.Unlock()
	return &RoomManager{
		rID: id,
		name: name,
		pIDList: make([]int32, 0),
		ownerPid: ownerPid,
	}
}

