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

func (adapter *FinsTCP) request(cmd []byte) ([]byte, error) {
	buf, e := adapter.link.Request(cmd)
	if e != nil {
		return nil, e
	}

	if len(buf) < 16+14 {
		//TODO error
	}

	//[UDP 10字节] [命令码 1 1] [结束码 0 0] , data
	//记录响应的SID
	adapter.SID = buf[25]
	//判断错误码
	code := helper.ParseUint16(buf[28:])
	if code != 0 {
		return nil, errors.New(fmt.Sprintf("错误码: %d", code))
	}
	return buf[30:], nil
}

func (adapter *FinsTCP) HandShake() error {
	//上位机节点号
	cmd := []byte{0, 0, 0, adapter.SA1}
	cmd = packTCPCommand(0, cmd)
	buf, e := adapter.link.Request(cmd)
	if e != nil {
		return e
	}
	if len(buf) < 24 {
		//TODO return error
	}

	//头16字节：FINS + 长度 + 命令 + 错误码

	//0x00, 0x00, 0x00, 0x01, // 客户端节点号
	//0x00, 0x00, 0x00, 0x01, // PLC端节点号
	//adapter.SA1 = buf[16+3]
	adapter.DA1 = buf[16+7]

	return nil
}

func (adapter *FinsTCP) ReadBit(code Code, addr uint16, bit uint8, length uint16) ([]bool, error) {
	cmd := buildReadBitCommand(code, addr, bit, length)
	cmd = adapter.packUDPCommand(cmd)
	cmd = packTCPCommand(2, cmd)

	buf, err := adapter.request(cmd)
	if err != nil {
		return nil, err
	}

	return helper.ByteToBool(buf), nil
}

func (adapter *FinsTCP) ReadWord(code Code, addr uint16, length uint16) ([]byte, error) {
	cmd := buildReadWordCommand(code, addr, length)
	cmd = adapter.packUDPCommand(cmd)
	cmd = packTCPCommand(2, cmd)
	return adapter.request(cmd)
}

func (adapter *FinsTCP) WriteBit(code Code, addr uint16, bit uint8, values []bool) error {
	v := helper.BoolToByte(values)

	cmd := buildWriteBitCommand(code, addr, bit, v)
	cmd = adapter.packUDPCommand(cmd)
	cmd = packTCPCommand(2, cmd)
	_, err := adapter.request(cmd)
	return err
}

func (adapter *FinsTCP) WriteWord(code Code, addr uint16, values []byte) error {
	cmd := buildWriteWordCommand(code, addr, values)
	cmd = adapter.packUDPCommand(cmd)
	cmd = packTCPCommand(2, cmd)
	_, err := adapter.request(cmd)
	return err
}
