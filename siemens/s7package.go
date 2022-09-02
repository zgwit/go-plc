package siemens

import (
	"fmt"
	"github.com/zgwit/go-plc/helper"
)

const (
	VariableTypeBit     = 1
	VariableTypeWord    = 2
	VariableTypeInt     = 3
	VariableTypeReal    = 4
	VariableTypeCounter = 5

	ParameterTypeRead  = 4
	ParameterTypeWrite = 5

	MessageTypeJobRequest = 1
	MessageTypeAck        = 2
	MessageTypeAckData    = 3
	MessageTypeUserData   = 4

	TransportSizeBit   = 3
	TransportSizeByte  = 4
	TransportSizeInt   = 5
	TransportSizeReal  = 7
	TransportSizeOctet = 9
)

//S7Parameter 固定14字节长
type S7Parameter struct {
	Code  uint8 //04 read 05 write
	Count uint8
	Type  uint8 //数据类型 1 bit 2 word 3 dint 4 real 5 counter

	Areas []S7ParameterArea
}

type S7ParameterArea struct {
	Code   uint8
	DB     uint16
	Size   uint16
	Offset uint32
}

func (p *S7Parameter) encode() []byte {
	blocks := uint8(len(p.Areas))
	buf := make([]byte, blocks*8+6)

	buf[0] = p.Code       // 4读 5写
	buf[1] = blocks       // 读取块数
	buf[2] = 0x12         //specification type 指定有效值类型
	buf[3] = blocks*8 + 2 //length 接下来本次地址访问长度
	buf[4] = 0x10         //syntax id 语法标记，ANY
	buf[5] = p.Type       //variable type 1 bit 2 word 3 dint 4 real 5 counter???

	cursor := 6
	for i, area := range p.Areas {
		helper.WriteUint16(buf[cursor:], area.Size)     // 访问数据的个数
		helper.WriteUint16(buf[cursor+2:], area.DB)     //db number DB块编号，如果访问的是DB块的话
		buf[cursor+4] = area.Code                       //area 访问数据类型
		helper.WriteUint24(buf[cursor+5:], area.Offset) //address 偏移位置
		cursor += i * 8
	}
	return buf
}

func (p *S7Parameter) decode(buf []byte) error {
	p.Code = buf[0]
	p.Count = buf[1]
	//返回内容 只有以上肉个字节

	return nil
}

type S7Data struct {
	Code  uint8 //0xff代表成功
	Type  uint8
	Count uint16
	Data  []byte
}

func (p *S7Data) encode() []byte {
	data := p.Data
	if p.Type == VariableTypeBit {
		data = helper.ShrinkBool(data)
		p.Count = uint16(len(data))
	}

	l := uint8(len(data))
	buf := make([]byte, l+4)

	buf[0] = p.Code
	buf[1] = p.Type
	helper.WriteUint16(buf[2:], p.Count)
	copy(buf[4:], data)
	return buf
}

func (p *S7Data) decode(buf []byte) (int, error) {
	p.Code = buf[0]
	//写入返回
	if len(buf) < 4 {
		return 0, fmt.Errorf("长度太短")
	}

	p.Type = buf[1]
	p.Count = helper.ParseUint16(buf[2:])
	length := int(p.Count)
	p.Data = buf[4 : length+4]
	if p.Type == VariableTypeBit {
		p.Data = helper.ExpandBool(buf, length*8)
	}
	return length + 4, nil
}

type S7Package struct {
	Type      uint8  // 1 Job Request 2 Ack 3 Ack-Data 7 Userdata
	Reference uint16 //序列号

	param S7Parameter
	data  []S7Data
	//ErrorClass uint8
	//ErrorCode  uint8
}

func (p *S7Package) encode() []byte {

	parameter := p.param.encode()
	paramLength := len(parameter)

	dataLength := 0
	datum := make([][]byte, 0)
	for _, p := range p.data {
		buf := p.encode()
		datum = append(datum, buf)
		dataLength += len(buf)
	}

	size := paramLength + dataLength + 17

	buf := make([]byte, size)
	//TPKT
	buf[0] = 0x03
	buf[1] = 0x00
	helper.WriteUint16(buf[2:], uint16(size)) // 长度
	//ISO-COTP
	buf[4] = 0x02 // 固定
	buf[5] = 0xF0
	buf[6] = 0x80
	//S7 communication
	buf[7] = 0x32   //Desc ID 协议ID，固定为32
	buf[8] = p.Type //Message Type(ROSCTR) 1 Job Request 2 Ack 3 Ack-Data 7 Userdata
	buf[9] = 0x0    //Reserved
	buf[10] = 0x0
	helper.WriteUint16LittleEndian(buf[11:], p.Reference) // PDU ref 标识序列号(可以像Modbus TCP一样使用)
	helper.WriteUint16(buf[13:], uint16(paramLength))     // Param length
	helper.WriteUint16(buf[15:], uint16(dataLength))      // Data length

	//复制参数
	cursor := 17
	if paramLength > 0 {
		copy(buf[cursor:], parameter)
		cursor += paramLength
	}

	//复制数据
	if dataLength > 0 {
		for _, p := range datum {
			copy(buf[cursor:], p)
			cursor += len(p)
		}
	}

	return buf
}

func (p *S7Package) decode(buf []byte) error {
	length := helper.ParseUint16(buf[2:])
	if len(buf) < int(length) {
		return fmt.Errorf("长度不够 %d %d", length, len(buf))
	}

	p.Type = buf[8]
	p.Reference = helper.ParseUint16LittleEndian(buf[11:])
	paramLength := helper.ParseUint16(buf[13:])
	dataLength := helper.ParseUint16(buf[15:])

	//仅当 ack-data
	ErrorClass := buf[17]
	ErrorCode := buf[18]
	if ErrorClass != 0 || ErrorCode != 0 {
		return fmt.Errorf("错误码：%d %d", ErrorClass, ErrorCode)
	}

	err := p.param.decode(buf[19:])
	if err != nil {
		return err
	}

	if dataLength == 0 {
		return nil
	}

	p.data = make([]S7Data, 0)

	cursor := int(paramLength) + 19
	remain := int(length) - cursor
	if p.Type == MessageTypeAckData && p.param.Code == ParameterTypeWrite {
		for remain > 0 {
			var d S7Data
			d.Code = buf[cursor]
			p.data = append(p.data, d)
			cursor++
			remain--
		}
	} else {
		for remain > 0 {
			var d S7Data
			cnt, err := d.decode(buf[cursor:])
			if err != nil {
				return err
			}
			p.data = append(p.data, d)
			cursor += cnt
			remain -= cnt
		}
	}

	return nil
}
