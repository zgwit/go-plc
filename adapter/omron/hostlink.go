package omron

import (
	"errors"
	"fmt"
	"github.com/zgwit/go-plc/helper"
	"github.com/zgwit/go-plc/link"
)

type HostLink struct {
	// 信息控制字段，默认0x80
	ICF byte // 0x80

	// PLC的节点地址，这个值在配置了ip地址之后是默认赋值的，默认为Ip地址的最后一位
	DA1 byte // 0x13

	// PLC的单元号地址
	DA2 byte // 0x00

	// 上位机的节点地址，假如你的电脑的Ip地址为192.168.0.13，那么这个值就是13
	SA1 byte

	// 上位机的单元号地址
	SA2 byte

	// 设备的标识号
	SID byte // 0x00

	link link.Link
}

func NewHostLink(link link.Link) *HostLink {
	return &HostLink{
		ICF:  0x80,
		link: link,
	}
}

func (adapter *HostLink) packCommand(payload []byte) []byte {
	cmd := helper.ToHex(payload)

	length := len(cmd)

	buf := make([]byte, 18+length)

	buf[0] = '@'
	helper.WriteByteHex(buf[1:], adapter.DA1) //PLC设备号
	helper.WriteByteHex(buf[3:], 0xFA)        //识别码
	buf[5] = 0x30                             //响应等待时间 x 15ms
	helper.WriteByteHex(buf[6:], adapter.ICF)
	helper.WriteByteHex(buf[8:], adapter.DA2)
	helper.WriteByteHex(buf[10:], adapter.SA2)
	helper.WriteByteHex(buf[12:], adapter.SID)
	copy(buf[14:], cmd)

	//计算FCS，异或校验
	helper.WriteByteHex(buf[length+14:], helper.Xor(buf[:length+14]))
	buf[length+16] = '*'
	buf[length+17] = 0x0D //CR

	return buf
}

func (adapter *HostLink) read(cmd []byte, expect int) ([]byte, error) {

	//@ [0 0 :单元号] [F A] [0 0] [4 0 :ICF][0 0 :DA2][0 0 :SA2][0 0 :SID]
	//[命令码 4字节] [状态码 4字节] [ ...data... ]
	//[FCS][* CR]

	buf, err := adapter.link.Request(cmd)
	if err != nil {
		return nil, err
	}

	//adapter.SID = FromHex(buf[13:15])[0]

	v := helper.FromHex(buf[15 : len(buf)-4])

	//[命令码 1 1] [结束码 0 0] , data
	code := helper.ParseUint16(v[2:])
	if code != 0 {
		return nil, errors.New(fmt.Sprintf("错误码: %d", code))
	}

	return v[4:], nil
}
func (adapter *HostLink) write(cmd []byte) error {
	buf, err := adapter.link.Request(cmd)
	if err != nil {
		return err
	}

	v := helper.FromHex(buf[15 : len(buf)-4])
	code := helper.ParseUint16(v[2:])
	if code != 0 {
		return errors.New(fmt.Sprintf("错误码: %d", code))
	}
	return nil
}

func (adapter *HostLink) ReadBit(code Code, addr uint16, bit uint8, length uint16) ([]bool, error) {
	cmd := buildReadBitCommand(code, addr, bit, length)
	cmd = adapter.packCommand(cmd)

	buf, err := adapter.read(cmd, int(length))
	if err != nil {
		return nil, err
	}

	return helper.ByteToBool(buf), nil
}

func (adapter *HostLink) ReadWord(code Code, addr uint16, length uint16) ([]byte, error) {
	cmd := buildReadWordCommand(code, addr, length)
	cmd = adapter.packCommand(cmd)
	return adapter.read(cmd, int(length))
}

func (adapter *HostLink) WriteBit(code Code, addr uint16, bit uint8, values []bool) error {
	v := helper.BoolToByte(values)
	cmd := buildWriteBitCommand(code, addr, bit, v)
	cmd = adapter.packCommand(cmd)
	return adapter.write(cmd)
}

func (adapter *HostLink) WriteWord(code Code, addr uint16, values []byte) error {
	cmd := buildWriteWordCommand(code, addr, values)
	cmd = adapter.packCommand(cmd)
	return adapter.write(cmd)
}
