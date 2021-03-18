package znet

import (
	"TankServer/zinx/ziface"
	"errors"
	"fmt"
	"sync"
)

/*
	链接管理实现模块
 */

type ConnManager struct {
	connections map[uint32]ziface.IConnection	// 管理链接集合
	connLock	sync.RWMutex					// 保护链接集合的读写锁
}

// 创建当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32] ziface.IConnection),
	}
}


func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将conn加入map中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connID = ",conn.GetConnID(), " connection add to ConnManager successfully : conn num = ", cm.Len())
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源map 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除链接信息
	delete(cm.connections, conn.GetConnID())

	fmt.Println("connID = ", conn.GetConnID(), "remove from ConnManager successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源map 加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection is not FOUND")
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	// 保护共享资源map 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}
	fmt.Println("Clear all connections success! conn num = ", cm.Len())
}
