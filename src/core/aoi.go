package core

/*
	AOI管理模块
*/

type AOIManager struct {
	MinX    int           //区域左边界坐标
	MaxX    int           //区域右边界坐标
	CopiesX int           //x方向格子的数量
	MinY    int           //区域上边界坐标
	MaxY    int           //区域下边界坐标
	CopiesY int           //y方向的格子数量
	grids   map[int]*Grid //当前区域中都有哪些格子，key=格子ID， value=格子对象
}

// 初始化aoi区域 在createRoom路由的handle中调用 只有在调用该方法之后才能调用worldManager中的createRoom方法
func NewAOIManager(minX, maxX, copiesX, minY, maxY, copiesY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:    minX,
		MaxX:    maxX,
		CopiesX: copiesX,
		MinY:    minY,
		MaxY:    maxY,
		CopiesY: copiesY,
		grids:   make(map[int]*Grid),
	}
	//给AOI初始化区域中所有的格子
	for y := 0; y < copiesY; y++ {
		for x := 0; x < copiesX; x++ {
			//计算格子ID
			//格子编号：id = idy *nx + idx  (利用格子坐标得到格子编号)
			gid := y*copiesX + x

			//初始化一个格子放在AOI中的map里，key是当前格子的ID
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.GetGridWith(),
				aoiMgr.MinX+(x+1)*aoiMgr.GetGridWith(),
				aoiMgr.MinY+y*aoiMgr.GetGridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.GetGridLength())
		}
	}
	return aoiMgr
}

// 得到每个格子在X轴的宽度
func (am *AOIManager) GetGridWith() int {
	return (am.MaxX - am.MinX) / am.CopiesX
}

// 得到每个格子在Y轴的长度
func (am *AOIManager) GetGridLength() int {
	return (am.MaxY - am.MinY) / am.CopiesY
}

// 根据格子gID得到视野范围内的九宫格格子的id集合
func (am *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	// 查看当前gid是否在格子中
	if _, ok := am.grids[gID]; !ok {
		return
	}
	// 将gid放入返回值列表
	grids = append(grids, am.grids[gID])
	// 根据gID得到对应坐标
	x, y := gID%am.CopiesX, gID/am.CopiesX
	// 设置偏移量数组 包含了周边格子对于gID坐标的偏移量
	offSetX := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	offSetY := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	// 遍历目标格子九宫格视野范围内的格子 得到对应坐标
	for i := 0; i < 8; i++ {
		newX := x + offSetX[i]
		newY := y + offSetY[i]
		// 校验新xy是否满足存在条件 即是否超出边界 未超出边界则将其放入返回值列表中
		if newX >= 0 && newX < am.CopiesX && newY >= 0 && newY < am.CopiesY {
			grids = append(grids, am.grids[newY*am.CopiesX+newX])
		}
	}
	return
}

// 通过横纵坐标得到当前gid格子编号
func (am *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - am.MinX) / am.GetGridWith()
	idy := (int(y) - am.MinY) / am.GetGridLength()
	return idy*am.CopiesX + idx
}

// 通过横纵坐标得到周边九宫格全部的PlayerIDs
func (am *AOIManager) GetPidByPos(x, y float32) (playerIDs []int) {
	// 得到当前玩家的GID格子id
	gID := am.GetGidByPos(x, y)
	// 通过gid得到视野范围内的格子信息
	grids := am.GetSurroundGridsByGid(gID)
	// 将视野范围内的格子里边全部player的id 累加到返回值列表中
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
	}
	return
}

// 添加一个player到格子中
func (am *AOIManager) AddPidToGrid(pID, gID int) {
	am.grids[gID].Add(pID)
}

// 移除一个格子中的playerid
func (am *AOIManager) RemovePidFromGrid(pID, gID int) {
	am.grids[gID].Remove(pID)
}

// 通过GID获取全部的PlayerID
func (am *AOIManager) GetPidByGid(gID int) (playerIDs []int) {
	playerIDs = am.grids[gID].GetPlayerIDs()
	return
}

// 通过坐标将Player添加到一个格子中
func (am *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := am.GetGidByPos(x,y)
	grid := am.grids[gID]
	grid.Add(pID)
}
// 通过坐标将Player从一个格子中移除
func (am *AOIManager) RemoveToGridByPos(pID int, x, y float32) {
	gID := am.GetGidByPos(x,y)
	grid := am.grids[gID]
	grid.Remove(pID)
}
