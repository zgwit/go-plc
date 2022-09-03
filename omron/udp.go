package omron

import (
	"errors"
	"fmt"
	"github.com/zgwit/go-plc/helper"
	"github.com/zgwit/go-plc/protocol"
	"io"
)

type UdpFrame struct {
	// 信息控制字段，默认0x80
	ICF byte // 0x80

	// 系统使用的内部信息
	RSV byte // 0x00

	// 网络层信息，默认0x02，如果有八层消息，就设置为0x07
	GCT byte // 0x02

	// PLC的网络号地址，默认0x00
	DNA byte // 0x00

	// PLC的节点地址，这个值在配置了ip地址之后是默认赋值的，默认为Ip地址的最后一位
	DA1 byte // 0x13

	// PLC的单元号地址
	DA2 byte // 0x00

	// 上位机的网络号地址
	SNA byte // 0x00

	// 上位机的节点地址，假如你的电脑的Ip地址为192.168.0.13，那么这个值就是13
	SA1 byte

	// 上位机的单元号地址
	SA2 byte

	// 设备的标识号
	SID byte // 0x00
}

type FinsUdp struct {
	frame UdpFrame
	link  protocol.Messenger
	buf   []byte
}

func NewFinsUDP(link io.ReadWriter, opts string) protocol.Protocol {
	fins := &FinsUdp{
		link: protocol.Messenger{Conn: link},
		buf:  make([]byte, 256)}
	return fins
}

func (f *FinsUdp) execute(cmd []byte) ([]byte, error) {
	//发送请求
	l, err := f.link.Ask(cmd, f.buf)
	if err != nil {
		return nil, err
	}

	//解析数据
	if l < 10 {
		return nil, errors.New("长度不够")
	}

	//[UDP 10字节]

	//记录响应的SID
	f.frame.SID = f.buf[9]

	return f.buf[10:l], nil
}

func (f *FinsUdp) Read(address protocol.Addr, size int) ([]byte, error) {

	//构建读命令
	buf, e := buildReadCommand(address, size)
	if e != nil {
		return nil, e
	}

	//打包命令
	cmd := packUDPCommand(&f.frame, buf)

	//发送请求
	recv, err := f.execute(cmd)
	if err != nil {
		return nil, err
	}

	//[命令码 1 1] [结束码 0 0] , data
	code := helper.ParseUint16(recv[2:])
	if code != 0 {
		return nil, fmt.Errorf("错误码: %d", code)
	}

	return recv[4:], nil
}

func (f *FinsUdp) Write(address protocol.Addr, values []byte) error {
	//构建写命令
	buf, e := buildWriteCommand(address, values)
	if e != nil {
		return e
	}

	//打包命令
	cmd := packUDPCommand(&f.frame, buf)

	//发送请求
	recv, err := f.execute(cmd)
	if err != nil {
		return err
	}
	//[命令码 1 1] [结束码 0 0]
	code := helper.ParseUint16(recv[2:])
	if code != 0 {
		return fmt.Errorf("错误码: %d", code)
	}

	return nil
}
