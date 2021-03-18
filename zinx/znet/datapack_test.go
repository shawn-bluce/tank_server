package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 负责测试datapack拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	*/
	// 1 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err ", err)
		return
	}

	// 创建go 负责从客户端处理业务
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept error", err)
		}

		// 处理客户端请求
		go func(conn net.Conn) {
			// 拆包过程
			// 定义拆包对象
			dp := NewDataPack()
			for {
				// 第一次读 将head头读出来
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head error", err)
					break
				}
				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("server unpack err", err)
					return
				}
				if msgHead.GetMsgLen() > 0 {
					// msg是有数据的 需要进行第二次读取
					// 第二次读 根据head中datalen读取data内容
					msg := msgHead.(*Message)
					msg.Data = make([]byte, msg.GetMsgLen())

					// 根据datalen的长度再次从io流中读取
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data err", err)
						return
					}
					// 完整的一个消息已经读取完毕
					fmt.Println("--> Recv MsgID:", msg.Id, ", datalen = ", msg.DataLen, ", data = ", string(msg.Data))
				}
			}
		}(conn)

	}()
	// 1 从客户端读取数据 拆包处理

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err ", err)
		return
	}
	// 创建一个封包对象
	dp := NewDataPack()

	// 模拟粘包过程 封装两个msg一同发送
	// 封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client package msg1 err", err)
		return
	}
	// 封装第二个msg2包
	msg2 := &Message{
		Id:      2,
		DataLen: 9,
		Data:    []byte{'z', 'i', 'n', 'x', 'h', 'e', 'l', 'l', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client package msg2 err", err)
		return
	}
	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
