package core

import (
	"strconv"
	"strings"
	"testing"
)

func TestWorldManager_ListRoom(t *testing.T) {
	//WorldMgrObj.ListRoom()

}

func TestWorldManager_CreateRo(t *testing.T) {
	//roomMgr := NewRoomManager("testroom")
	//fmt.Println(roomMgr)
	//aoiMgr := NewAOIManager(0,500,5,0,500,5)
	//fmt.Println(aoiMgr)
	////WorldMgrObj.CreateRoom( "testroom", 1)
	//WorldMgrObj.CreateRoom(roomMgr, 1)
	//roomMgr2 := NewRoomManager("testroom2")
	//
	//WorldMgrObj.CreateRoom(roomMgr2, 2)
	//
	//// 从协议中解析得到的信息 aoi地图大小
	//NewAOIManager(0,100,5,0,100,5)
	//fmt.Println("roomlist==========  ",WorldMgrObj.ListRoom())
	//for name, room := range WorldMgrObj.ListRoom() {
	//	fmt.Printf("name:: %v  room:: %v  room.pidlist:: %v\n", name, room, room.pIDList)
	//}

	// 模拟listroom的Hanle
	// 解析protobuf数据 但是listroom没有数据传递过来

	// 组合房间切片 [] *RoomManager
	//roomManagers := make([]*RoomManager, 0)
	//for _, v := range WorldMgrObj.ListRoom() {
	//	roomManagers = append(roomManagers, v)
	//}
	//fmt.Println("rooms:::", roomManagers)


	// 模拟createroom的Handle
	// 解析protobuf数据 即获取地图信息 player对象 房间名称 并通过aoi地图信息、房间信息分别初始化aoi地图和房间 然后调用创建房间方法 最后还需要调用player中的发送消息方法
	// 其中listroom需要将房间信息封装成protobuf数据进行下发
	//roomName := "lambo"
	//aoi := NewAOIManager(0,100,5,0,100,5)
	//room := NewRoomManager(roomName, 1)
	//WorldMgrObj.CreateRoom(room, aoi, 233)
	//
	//roomName2 := "kiko"
	//aoi2 := NewAOIManager(0,100,5,0,100,5)
	//room2 := NewRoomManager(roomName2, 2)
	//WorldMgrObj.CreateRoom(room2, aoi2, 820)
	//
	//
	//roomManagerss := make([]*RoomManager, 0)
	//for _, v := range WorldMgrObj.ListRoom() {
	//	roomManagerss = append(roomManagerss, v)
	//}
	//fmt.Println("rooms:::", roomManagerss)
	//for _, v := range roomManagerss {
	//	fmt.Println("room ===========> name:", v.name, " pidList:", v.pIDList)
	//}
	//
	//WorldMgrObj.RemoveRoom("lambo")
	//for _, v := range roomManagerss {
	//	fmt.Println("after remove room ===========> name:", v.name, " pidList:", v.pIDList)
	//}


}

func BenchmarkWorldManager_CreateRoom(b *testing.B) {
	var str strings.Builder
	str.WriteString("test")
	for i := 0; i < b.N; i++ {
		str.WriteString(strconv.Itoa(i))
		aoi := NewAOIManager(0,100,5,0,100,5)
		room := NewRoomManager(str.String(), int32(i))
		WorldMgrObj.CreateRoom(room, aoi, int32(i))
	}
}
