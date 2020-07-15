package omron

import (
	"errors"
	"fmt"
	"github.com/zgwit/go-plc/helper"
	"github.com/zgwit/go-plc/link"
)

type FinsTCP struct {
	Fins

	link link.Link
}

func NewFinsTCP(link link.Link) *FinsTCP {
	return &FinsTCP{
		Fins: Fins{
			ICF: 0x80,
			GCT: 0x02,
		},
		link: link,
	}
}

func (adapter *FinsTCP) read(cmd []byte, expect int) ([]byte, error) {
	buf, e := adapter.link.Request(cmd)
	if e != nil {
		return nil, e
	}

	if len(buf) < 14 {
		//TODO error
	}

	//[UDP 10字节] [命令码 1 1] [结束码 0 0] , data
	//记录响应的SID
	adapter.SID = buf[9]
	//判断错误码
	code := helper.ParseUint16(buf[12:])
	if code != 0 {
		return nil, errors.New(fmt.Sprintf("错误码: %d", code))
	}
	return buf[14:], nil
}

func (adapter *FinsTCP) write(cmd []byte) error {
	_, e := adapter.read(cmd, 0)
	return e
}

func (adapter *FinsTCP) request(cmd []byte) ([]byte, error) {
	return nil, nil
}

func (adapter *FinsTCP) HandShake() error {
	packTCPCommand(0, []byte{0,0,0,1})
	return nil
}

func (adapter *FinsTCP) ReadBit(code Code, addr uint16, bit uint8, length uint16) ([]bool, error) {
	cmd := buildReadBitCommand(code, addr, bit, length)
	cmd = adapter.packCommand(cmd)
	cmd = packTCPCommand(2, cmd)

	buf, err := adapter.read(cmd, int(length))
	if err != nil {
		return nil, err
	}

	return helper.ByteToBool(buf), nil
}

func (adapter *FinsTCP) ReadWord(code Code, addr uint16, length uint16) ([]byte, error) {
	cmd := buildReadWordCommand(code, addr, length)
	cmd = adapter.packCommand(cmd)
	cmd = packTCPCommand(2, cmd)
	return adapter.read(cmd, int(length))
}

func (adapter *FinsTCP) WriteBit(code Code, addr uint16, bit uint8, values []bool) error {
	v := helper.BoolToByte(values)
	cmd := buildWriteBitCommand(code, addr, bit, v)
	cmd = adapter.packCommand(cmd)
	cmd = packTCPCommand(2, cmd)
	return adapter.write(cmd)
}

func (adapter *FinsTCP) WriteWord(code Code, addr uint16, values []byte) error {
	cmd := buildWriteWordCommand(code, addr, values)
	cmd = adapter.packCommand(cmd)
	cmd = packTCPCommand(2, cmd)
	return adapter.write(cmd)
}
