package mitsubishi

import (
	"errors"
	"fmt"
	"github.com/zgwit/go-plc/helper"
	"io"
)

// A3EBinaryAdapter 协议
type A3EBinaryAdapter struct {
	StationNumber byte //站编号
	NetworkNumber byte //网络编号
	PlcNumber     byte //PLC编号
	IoNumber      byte //IO编号

	link io.ReadWriter
}

func NewA3EBinaryAdapter() *A3EBinaryAdapter {
	a := A3EBinaryAdapter{}
	a.StationNumber = 0
	a.NetworkNumber = 0
	a.PlcNumber = 0xFF
	a.IoNumber = 0xFF
	return &a
}

func (t *A3EBinaryAdapter) request(cmd []byte) ([]byte, error) {

	// 副标题 D0 00 网络号 00 PLC号 FF IO编号 FF 03 站号 00 应答长度 L H 结束代码 L H
	//
	buf := make([]byte, 16)
	//_, e := t.link.Read(buf)
	//if e != nil {
	//	return nil, e
	//}

	status := helper.ParseUint32(buf[12:])
	if status != 0 {
		return nil, errors.New(fmt.Sprintf("TCP状态错误: %d", status))
	}

	length := helper.ParseUint32(buf[4:]) - 8

	payload := make([]byte, length)
	//t.link.Read(payload)

	return payload, nil
}

func (t *A3EBinaryAdapter) BuildCommand(cmd []byte) []byte {
	length := len(cmd)

	buf := make([]byte, 11+length)
	buf[0] = 0x50 //副标题
	buf[1] = 0x00
	buf[2] = t.NetworkNumber                                //网络号
	buf[3] = t.PlcNumber                                    //PLC编号
	buf[4] = t.IoNumber                                     //目标IO编号 L
	buf[5] = 0x03                                           //目标IO编号 H
	buf[6] = t.StationNumber                                //站编号
	helper.WriteUint16LittleEndian(buf[7:], uint16(length)) //请求数据长度
	helper.WriteUint16LittleEndian(buf[9:], 10)             //CPU监视定时器
	//命令
	copy(buf[11:], cmd)

	return buf
}

// Read 读取数据
func (t *A3EBinaryAdapter) Read(address string, length int) ([]byte, error) {

	//解析地址
	addr, e := ParseAddress(address)
	if e != nil {
		return nil, e
	}

	//构建命令
	buf := make([]byte, 10)
	buf[0] = 0x01 // 批量读取数据命令
	buf[1] = 0x04
	// 以点为单位还是字为单位成批读取
	if addr.IsBit {
		buf[2] = 0x01
	} else {
		buf[2] = 0x00
	}
	buf[3] = 0x00
	helper.WriteUint24LittleEndian(buf[4:], uint32(addr.Addr)) // 起始地址的地位
	buf[7] = addr.Code                                         // 指明读取的数据
	helper.WriteUint16LittleEndian(buf[8:], uint16(length))    //软元件的长度

	//构建命令
	cmd := t.BuildCommand(buf)

	recv, err := t.request(cmd)
	if err != nil {
		return nil, err
	}

	//TODO 解压位， 0x10 10  =>  0x01 00 01 00

	return recv, nil

}

func (t *A3EBinaryAdapter) Write(address string, values []byte) error {
	//解析地址
	addr, e := ParseAddress(address)
	if e != nil {
		return e
	}

	length := len(values)

	//TODO 压缩位， 0x01 00 01 00 => 0x10 10，读取的时候，也要解压

	//构建命令
	buf := make([]byte, 10+length)
	buf[0] = 0x01 // 批量写入数据命令
	buf[1] = 0x14
	// 以点为单位还是字为单位成批读取
	if addr.IsBit {
		buf[2] = 0x01
	} else {
		buf[2] = 0x00
	}
	buf[3] = 0x00
	helper.WriteUint24LittleEndian(buf[4:], uint32(addr.Addr)) // 起始地址的地位
	buf[7] = addr.Code                                         // 指明写入的数据
	helper.WriteUint16LittleEndian(buf[8:], uint16(length))    //软元件的长度

	copy(buf[10:], values)

	//构建命令
	cmd := t.BuildCommand(buf)

	//发送命令
	_, err := t.request(cmd)
	return err
}
