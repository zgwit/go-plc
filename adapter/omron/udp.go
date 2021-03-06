package omron

import (
	"errors"
	"fmt"
	"github.com/zgwit/go-plc/helper"
	"github.com/zgwit/go-plc/link"
)

type FinsUDP struct {
	Fins

	link link.Link
}

func NewFinsUDP(link link.Link) *FinsUDP {
	return &FinsUDP{
		Fins: Fins{
			ICF: 0x80,
			GCT: 0x02,
		},
		link: link,
	}
}

func (adapter *FinsUDP) request(cmd []byte) ([]byte, error) {
	buf, err := adapter.link.Request(cmd)
	if err != nil {
		return nil, err
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

func (adapter *FinsUDP) ReadBit(code Code, addr uint16, bit uint8, length uint16) ([]bool, error) {
	cmd := buildReadBitCommand(code, addr, bit, length)
	cmd = adapter.packUDPCommand(cmd)

	buf, err := adapter.request(cmd)
	if err != nil {
		return nil, err
	}

	return helper.ByteToBool(buf), nil
}

func (adapter *FinsUDP) ReadWord(code Code, addr uint16, length uint16) ([]byte, error) {
	cmd := buildReadWordCommand(code, addr, length)
	cmd = adapter.packUDPCommand(cmd)
	return adapter.request(cmd)
}

func (adapter *FinsUDP) WriteBit(code Code, addr uint16, bit uint8, values []bool) error {
	v := helper.BoolToByte(values)
	cmd := buildWriteBitCommand(code, addr, bit, v)
	cmd = adapter.packUDPCommand(cmd)
	_, err := adapter.request(cmd)
	return err
}

func (adapter *FinsUDP) WriteWord(code Code, addr uint16, values []byte) error {
	cmd := buildWriteWordCommand(code, addr, values)
	cmd = adapter.packUDPCommand(cmd)
	_, err := adapter.request(cmd)
	return err
}
