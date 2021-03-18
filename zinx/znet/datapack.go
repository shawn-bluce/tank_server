package znet

import (
	"TankServer/zinx/utils"
	"TankServer/zinx/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

// 封包 拆包具体模块
type DataPack struct {}

// 拆包封包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头的长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32 (4byte) ID uint32 (4byte) = 8byte
	return 8
}
// 封包方法
// |datalen|msgID|data|
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error)  {
	// 创建存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	// 将DataLen写进dataBuff中 大小端问题
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将MsgId写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将data数据写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}
// 拆包方法 将包的head信息读取出来 根据head中的data长度读取信息内容
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	// 只解压head信息，得到dataLen和MsgID
	msg := &Message{}
	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	// 判断dataLen是否已经超出了允许最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data receive")
	}
	return msg, nil
}