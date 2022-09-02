package modbus

import (
	"errors"
	"fmt"
	"github.com/zgwit/go-plc/helper"
	"github.com/zgwit/go-plc/protocol"

	"strconv"
)

// TCP Modbus-TCP协议
type TCP struct {
	link protocol.Messenger
	buf  []byte
}

func NewTCP(link io.ReadWriter, opts string) protocol.Protocol {
	tcp := &TCP{
		link: protocol.Messenger{Conn: link},
		buf:  make([]byte, 260),
	}
	return tcp
}

func (m *TCP) execute(cmd []byte) ([]byte, error) {
	helper.WriteUint16(cmd, 0x0A0A) //写入事务ID

	//下发指令
	l, err := m.link.AskAtLeast(cmd, m.buf, 10)
	if err != nil {
		return nil, err
	}
	buf := m.buf[:l]

	length := helper.ParseUint16(buf[4:])
	packLen := int(length) + 6
	if packLen > l {
		return nil, errors.New("长度不够")
	}

	//slave := buf[6]
	fc := buf[7]
	//解析错误码
	if fc&0x80 > 0 {
		return nil, fmt.Errorf("错误码：%d", buf[2])
	}

	//解析数据
	//length := 4
	count := int(buf[8])
	switch fc {
	case FuncCodeReadDiscreteInputs,
		FuncCodeReadCoils:
		//数组解压
		bb := helper.ExpandBool(buf[9:], count)
		return bb, nil
	case FuncCodeReadInputRegisters,
		FuncCodeReadHoldingRegisters,
		FuncCodeReadWriteMultipleRegisters:
		return helper.Dup(buf[9:]), nil
	case FuncCodeWriteSingleCoil, FuncCodeWriteMultipleCoils,
		FuncCodeWriteSingleRegister, FuncCodeWriteMultipleRegisters:
		//写指令不处理
		return nil, nil
	default:
		return nil, fmt.Errorf("错误功能码：%d", fc)
	}
}

func (m *TCP) Read(station int, area string, addr string, size int) ([]byte, error) {
	code := parseCode(area)
	offset, err := strconv.ParseUint(addr, 10, 16)
	if err != nil {
		return nil, err
	}

	b := make([]byte, 12)
	//helper.WriteUint16(b, id)
	helper.WriteUint16(b[2:], 0) //协议版本
	helper.WriteUint16(b[4:], 6) //剩余长度
	b[6] = uint8(station)
	b[7] = code
	helper.WriteUint16(b[8:], uint16(offset))
	helper.WriteUint16(b[10:], uint16(size))

	return m.execute(b)
}

func (m *TCP) Write(station int, area string, addr string, buf []byte) error {

	code := parseCode(area)
	offset, err := strconv.ParseUint(addr, 10, 16)
	if err != nil {
		return err
	}

	length := len(buf)
	switch code {
	case FuncCodeReadCoils:
		if length == 1 {
			code = 5
			//数据 转成 0x0000 0xFF00
			if buf[0] > 0 {
				buf = []byte{0xFF, 0}
			} else {
				buf = []byte{0, 0}
			}
		} else {
			code = 15 //0x0F
			//数组压缩
			b := helper.ShrinkBool(buf)
			count := len(b)
			buf = make([]byte, 3+count)
			helper.WriteUint16(buf, uint16(length))
			buf[2] = uint8(count)
			copy(buf[3:], b)
		}
	case FuncCodeReadHoldingRegisters:
		if length == 2 {
			code = 6
		} else {
			code = 16 //0x10
			b := make([]byte, 3+length)
			helper.WriteUint16(b, uint16(length/2))
			b[2] = uint8(length)
			copy(b[3:], buf)
			buf = b
		}
	default:
		return errors.New("功能码不支持")
	}

	l := 10 + len(buf)
	b := make([]byte, l)
	//helper.WriteUint16(b, id)
	helper.WriteUint16(b[2:], 0) //协议版本
	helper.WriteUint16(b[4:], 6) //剩余长度
	b[6] = uint8(station)
	b[7] = code
	helper.WriteUint16(b[8:], uint16(offset))
	copy(b[10:], buf)

	_, err = m.execute(b)
	return err
}
