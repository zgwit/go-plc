package omron

import (
	"errors"
	"fmt"
	"iot-master/connect"
	"iot-master/helper"
	"iot-master/protocols/protocol"
	"time"
)

type Fins struct {
	frame UdpFrame
	link  connect.Tunnel
}

func NewFinsTCP(link connect.Tunnel, opts protocol.Options) protocol.Protocol {
	fins := &Fins{link: link}
	link.On("data", func(data []byte) {
		//fins.OnData(data)
	})
	link.On("close", func() {
		//close(fins.queue)
	})
	return fins
}

func (f *Fins) Desc() *protocol.Desc {
	return &DescTCP
}

func (f *Fins) execute(cmd []byte) ([]byte, error) {
	//发送请求
	buf, err := f.link.Ask(cmd, time.Second*5)
	if err != nil {
		return nil, err
	}

	//解析数据
	l := len(buf)
	if l < 16 {
		return nil, errors.New("长度不够")
	}

	//头16字节：FINS + 长度 + 命令 + 错误码
	status := helper.ParseUint32(buf[12:])
	if status != 0 {
		return nil, fmt.Errorf("TCP状态错误: %d", status)
	}

	length := helper.ParseUint32(buf[4:])
	//判断剩余长度
	if int(length)+8 < len(buf) {
		return nil, fmt.Errorf("长度错误: %d", length)
	}

	return buf[16:], nil
}

func (f *Fins) Handshake() error {

	// 节点号
	handshake := []byte{0x00, 0x00, 0x00, 0x01}

	cmd := packTCPCommand(0, handshake)

	//发送请求
	buf, e := f.execute(cmd)
	if e != nil {
		return e
	}

	//0x00, 0x00, 0x00, 0x01, // 客户端节点号
	//0x00, 0x00, 0x00, 0x01, // PLC端节点号

	//客户端节点
	//f.SA1 = buf[3]
	//服务端节点
	f.frame.DA1 = buf[7]

	return nil
}

func (f *Fins) Read(station int, address protocol.Addr, size int) ([]byte, error) {

	//构建读命令
	buf, e := buildReadCommand(address, size)
	if e != nil {
		return nil, e
	}

	//打包命令
	cmd := packTCPCommand(2, packUDPCommand(&f.frame, buf))

	//发送请求
	recv, err := f.execute(cmd)
	if err != nil {
		return nil, err
	}

	//[UDP 10字节] [命令码 1 1] [结束码 0 0] , data

	code := helper.ParseUint16(recv[12:])
	if code != 0 {
		return nil, fmt.Errorf("错误码: %d", code)
	}

	//记录响应的SID
	f.frame.SID = recv[9]

	return recv[14:], nil
}

func (f *Fins) Poll(station int, addr protocol.Addr, size int) ([]byte, error) {
	return f.Read(station, addr, size)
}

func (f *Fins) Write(station int, address protocol.Addr, values []byte) error {
	//构建写命令
	buf, e := buildWriteCommand(address, values)
	if e != nil {
		return e
	}

	//打包命令
	cmd := packTCPCommand(2, packUDPCommand(&f.frame, buf))

	//发送请求
	recv, err := f.execute(cmd)
	if err != nil {
		return err
	}

	//[UDP 10字节] [命令码 1 1] [结束码 0 0]
	code := helper.ParseUint32(recv[12:])
	if code != 0 {
		return fmt.Errorf("错误码: %d", code)
	}

	//记录响应的SID
	f.frame.SID = recv[9]

	return nil
}
