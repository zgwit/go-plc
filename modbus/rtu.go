package modbus

import (
	"errors"
	"fmt"
	"github.com/zgwit/go-plc/helper"
	"github.com/zgwit/go-plc/protocol"
	"io"
	"strconv"
)

type response struct {
	buf []byte
	err error
}

type request struct {
	cmd  []byte
	resp chan response //out
}

// RTU Modbus-RTU协议
type RTU struct {
	link io.ReadWriter
	buf  []byte
}

func NewRTU(link io.ReadWriter, opts string) protocol.Protocol {
	//TODO parse opts(yaml)
	rtu := &RTU{
		link: link,
		//slave: opts["slave"].(uint8),
		buf: make([]byte, 256),
	}

	return rtu
}

func (m *RTU) execute(cmd []byte) ([]byte, error) {

	_, err := m.link.Write(cmd)
	if err != nil {
		return nil, err
	}

	l, err := m.link.Read(m.buf)
	if l < 5 {
		return nil, errors.New("长度不足")
	}

	crc := helper.ParseUint16LittleEndian(m.buf[l-2:])

	if crc != CRC16(m.buf[:l-2]) {
		//检验错误
		return nil, errors.New("校验错误")
	}

	//slave := buf[0]
	fc := m.buf[1]

	//解析错误码
	if fc&0x80 > 0 {
		return nil, fmt.Errorf("错误码：%d", m.buf[2])
	}

	//解析数据
	length := 4
	count := int(m.buf[2])
	switch fc {
	case FuncCodeReadDiscreteInputs,
		FuncCodeReadCoils:
		length += 1 + count/8
		if count%8 != 0 {
			length++
		}

		if l < length {
			//长度不够
			return nil, errors.New("长度不足")
		}
		b := m.buf[3 : l-2]
		//数组解压
		bb := helper.ExpandBool(b, count)
		return bb, nil
	case FuncCodeReadInputRegisters,
		FuncCodeReadHoldingRegisters,
		FuncCodeReadWriteMultipleRegisters:
		length += 1 + count
		if l < length {
			//长度不够
			return nil, errors.New("长度不足")
		}
		b := m.buf[3 : l-2]
		return helper.Dup(b), nil
	case FuncCodeWriteSingleCoil, FuncCodeWriteMultipleCoils,
		FuncCodeWriteSingleRegister, FuncCodeWriteMultipleRegisters:
		//写指令不处理
		return nil, nil
	default:
		return nil, errors.New("不支持的指令")
	}
}

func (m *RTU) Read(station int, area string, addr string, size int) ([]byte, error) {
	code := parseCode(area)
	offset, err := strconv.ParseUint(addr, 10, 16)
	if err != nil {
		return nil, err
	}

	b := make([]byte, 8)
	b[0] = uint8(station)
	b[1] = code
	helper.WriteUint16(b[2:], uint16(offset))
	helper.WriteUint16(b[4:], uint16(size))
	helper.WriteUint16LittleEndian(b[6:], CRC16(b[:6]))

	return m.execute(b)
}

func (m *RTU) Write(station int, area string, addr string, buf []byte) error {
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

	l := 6 + len(buf)
	b := make([]byte, l)
	b[0] = uint8(station)
	b[1] = code
	helper.WriteUint16(b[2:], uint16(offset))
	copy(b[4:], buf)
	helper.WriteUint16LittleEndian(b[l-2:], CRC16(b[:l-2]))

	_, err = m.execute(b)
	return err
}
