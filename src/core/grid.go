package core

import (
	"sync"
)

/*
	AOI地图中的格子类型
 */

type Grid struct {
	GID       int          //格子ID
	MinX      int          //格子左边界坐标
	MaxX      int          //格子右边界坐标
	MinY      int          //格子上边界坐标
	MaxY      int          //格子下边界坐标
	playerIDs map[int]bool //当前格子内的玩家或者物体成员ID
	pIDLock   sync.RWMutex //playerIDs的保护map的锁
}

// 初始化一个格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID: gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// 向当前格子添加一个玩家
func (g *Grid) Add(playerID int)  {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

// 从格子中删除一个玩家
func (g *Grid) Remove(playerID int)  {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// 得到当前格子中的所有玩家
func (g *Grid) GetPlayerIDs() (playerIDs []int)  {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()
	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return
}