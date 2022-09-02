package siemens

import (
	"github.com/zgwit/go-plc/helper"
	"github.com/zgwit/go-plc/protocol"
)

type Simatic struct {
	handshake1 []byte
	handshake2 []byte

	link io.ReadWriter
	desc *protocol.Desc
}

func (s *Simatic) Init() {
	s.link.On("online", func() {
		_ = s.HandShake()
	})
	return
}

func (s *Simatic) Desc() *protocol.Desc {
	return &DescS7_200_Smart
}

func (s *Simatic) HandShake() error {
	_, err := s.link.Ask(s.handshake1, 5)
	if err != nil {
		return err
	}
	//TODO 检查结果
	_, err = s.link.Ask(s.handshake2, 5)
	if err != nil {
		return err
	}
	//TODO 检查结果
	return nil
}

func (s *Simatic) Read(station int, addr protocol.Addr, size int) ([]byte, error) {
	address := addr.(*Address)

	buf := make([]byte, 14)
	buf[0] = 0x04                                // 4读 5写
	buf[1] = 1                                   // 读取块数
	buf[2] = 0x12                                //specification type 指定有效值类型
	buf[3] = 0x0A                                //length 接下来本次地址访问长度
	buf[4] = 0x10                                //syntax id 语法标记，ANY
	buf[5] = 0x02                                //variable type 1 bit 2 word 3 dint 4 real 5 counter???
	helper.WriteUint16(buf[6:], uint16(size))    // 访问数据的个数
	helper.WriteUint16(buf[8:], address.DB)      //db number DB块编号，如果访问的是DB块的话
	buf[10] = address.Code                       //area 访问数据类型
	helper.WriteUint24(buf[11:], address.Offset) //address 偏移位置

	cmd := packCommand(buf)

	resp, err := s.link.Ask(cmd, 5)
	if err != nil {
		return nil, err
	}

	//TODO 解析数据

	return resp, nil
}

func (s *Simatic) Poll(station int, addr protocol.Addr, size int) ([]byte, error) {
	return s.Read(station, addr, size)
}

func (s *Simatic) Write(station int, addr protocol.Addr, data []byte) error {
	address := addr.(*Address)
	length := len(data)

	buf := make([]byte, 14)
	buf[0] = 0x05                                // 4读 5写
	buf[1] = 1                                   // 读取块数
	buf[2] = 0x12                                // 指定有效值类型
	buf[3] = 0x0A                                // 接下来本次地址访问长度
	buf[4] = 0x10                                // 语法标记，ANY
	buf[5] = 0x02                                // 按字为单位，1 位 2 字
	helper.WriteUint16(buf[6:], uint16(length))  // 访问数据的个数
	helper.WriteUint16(buf[8:], address.DB)      // DB块编号，如果访问的是DB块的话
	buf[10] = address.Code                       // 访问数据类型
	helper.WriteUint24(buf[11:], address.Offset) // 偏移位置
	// 按字写入
	buf[14] = 0x00
	buf[15] = 0x04
	helper.WriteUint16(buf[16:], uint16(length*8)) // 按位计算的长度

	//添加数据
	copy(buf[18:], data)

	cmd := packCommand(buf)

	_, err := s.link.Ask(cmd, 5)
	if err != nil {
		return err
	}

	//TODO 解析结果

	return nil
}

//packCommand 打包命令
func packCommand(cmd []byte) []byte {
	length := len(cmd)

	buf := make([]byte, length+17)
	//TPKT
	buf[0] = 0x03
	buf[1] = 0x00
	helper.WriteUint16(buf[2:], uint16(length+17)) // 长度
	//ISO-COTP
	buf[4] = 0x02 // 固定
	buf[5] = 0xF0
	buf[6] = 0x80

	//S7 communication
	buf[7] = 0x32 //Desc ID 协议ID，固定为32
	buf[8] = 0x01 //Message Type(ROSCTR) 1 Job Request 2 Ack 3 Ack-Data 7 Userdata
	buf[9] = 0x0  //Reserved
	buf[10] = 0x0
	helper.WriteUint16LittleEndian(buf[11:], 0)  // PDU ref 标识序列号(可以像Modbus TCP一样使用)
	helper.WriteUint16(buf[13:], uint16(length)) // Param length
	helper.WriteUint16(buf[15:], 0)              // Data length

	//发送请求Request job
	//buf[17]功能码 04读 05写
	//buf[18]块数

	//仅出现在Ack-Data消息中
	//buf[17] Error class 0x00无 0x81应用程序关系错误 0x82对象定义错误 0x83无效资源可用错误 0x84服务处理错误 0x85请求错误 0x87访问错误
	//buf[18] Error Code

	//Parameter区
	//读取或写入
	//buf[19]0x12 读写是固定的
	//buf[20]此项剩余长度
	//buf[21]常量0x10 语法标记
	//buf[22]数据类型 1 2
	//buf[23 24]读取长度
	//buf[25 26]DB号
	//buf[27]Code
	//buf[28 29 30]地址

	//Data数据区
	//读取的结果
	//buf[19] Return Code 0xff 代表成功
	//buf[20] Variable type 数据类型
	//buf[21 22] Count
	//buf[23+]数据，长度是len(variable)*count

	//写入的内容 和 结果
	//buf[31] Return Code 0x00 固定
	//buf[32] Variable type 数据类型
	//buf[33 34] Count
	//buf[35+]数据，长度是len(variable)*count

	copy(buf[17:], cmd)

	return buf
}
