package znet

import (
	"TankServer/zinx/utils"
	"TankServer/zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

/*
	链接模块
 */
type Connection struct {
	// 当前Conn隶属于哪个server
	TcpServer ziface.IServer
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn
	// 链接的ID
	ConnID uint32
	// 当前链接的状态
	isClosed bool

	// 告知当前链接已经退出/停止 channel(由reader告知writer)
	ExitChan chan bool

	// 无缓冲管道 用于读、写goroutine之间的消息通信
	msgChan chan []byte

	// 消息的管理MsgID 和对应的处理业务API关系
	MsgHandle ziface.IMsgHandle

	// 链接属性集合
	property map[string]interface{}
	// 保护链接属性的锁
	propertyLock sync.RWMutex
}

// 初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		ExitChan:  make(chan bool,1),
		MsgHandle:msgHandle,
		msgChan: make(chan []byte),
		property: make(map[string]interface{}),
	}

	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// 链接的读业务方法
func (c *Connection) StartReader()  {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("connID = ", c.ConnID, "[Reader is exit], remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for  {
		// 创建一个拆包解包对象
		dp := NewDataPack()

		// 读取客户端Msg Head 的二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		// 拆包 得到msgID 和msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		// 根据dataLen 再次读取Data 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)
		// 得到当前conn数据的Request请求数据
		req := Request {
			conn: c,
			msg: msg,
		}

		// 判断是否已开启了工作池
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制 将消息发送给工作池即可
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中找到注册帮i的那个的Conn对应的router调用
			// 根据绑定好的MsgID 找到对应处理API业务 执行
			go c.MsgHandle.DoMsgHandler(&req)
		}
	}
}

/*
	写消息goroutine 专门发送给客户端消息的模块
 */
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")

	// 不断阻塞的等待channel的消息 进行写给客户端
	for {
		select {
		case data := <- c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <- c.ExitChan:
			// 代表Reader已经推出 此时Writer也要退出
			return
		}
	}

}

// 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Connection Start() ... ConnID = ", c.ConnID)
	// 启动当前链接的读数据业务
	go c.StartReader()
	// TODO 启动从当前链接写数据的业务
	go c.StartWriter()

	// 按照开发者传递进来的 创建链接之后需要调用的处理业务 执行对应Hook函数
	c.TcpServer.CallOnConnStart(c)
}

// 停止链接 结束当前链接的工作
func (c *Connection) Stop()  {
	fmt.Println("Connection Stop()... ConnID = ", c.ConnID)

	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 调用开发者注册的 销毁链接之前 需要执行的业务Hook函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket链接
	_ = c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前链接从ConnMsg中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前链接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn  {
	return c.Conn
}

// 获取当前连接模块的链接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP PORT
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 提供一个SendMsg 方法 将我们要发送给客户端的数据先进行封包 再发送
// SendMsg发送的消息不再直接是客户端，而是与Writer通信的管道
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	// 将data进行封包 MsgDataLen | MsgID | Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("pack error msg")
	}

	c.msgChan <- binaryMsg
	return nil
}


func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 添加一个链接属性
	c.property[key] = value
}

// 读取属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("value is not exist")
	}
}

// 删除属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
