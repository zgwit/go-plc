package omron

import (
	"errors"
	"fmt"
	helper2 "github.com/zgwit/go-plc/helper"
	"github.com/zgwit/go-plc/protocol"
	"io"
)

type FinsHostLink struct {
	frame UdpFrame
	link  protocol.Messenger
	buf   []byte
}

func NewHostLink(link io.ReadWriter, opts string) protocol.Protocol {
	fins := &Fins{
		link: protocol.Messenger{Conn: link},
		buf:  make([]byte, 256)}
	return fins
}

func (f *FinsHostLink) execute(cmd []byte) ([]byte, error) {
	//发送请求
	l, err := f.link.Ask(cmd, f.buf)
	if err != nil {
		return nil, err
	}

	//解析数据
	if l < 23 {
		return nil, errors.New("长度不够")
	}

	//@ [单元号] [F A] [0 0] [4 0 ICF][0 0 DA2][0 0 SA2][ SID ]
	//[命令码 4字节] [状态码 4字节] [ ...data... ]
	//[FCS][* CR]
	recv := helper2.FromHex(f.buf[15 : l-4])

	//记录响应的SID
	//t.frame.SID = FromHex(payload[13:15])[0]

	return recv, nil
}

func (f *FinsHostLink) Read(station int, address protocol.Addr, size int) ([]byte, error) {

	//构建读命令
	buf, e := buildReadCommand(address, size)
	if e != nil {
		return nil, e
	}

	//打包命令
	cmd := packAsciiCommand(&f.frame, buf)

	//发送请求
	recv, err := f.execute(cmd)
	if err != nil {
		return nil, err
	}

	//[命令码 1 1] [结束码 0 0] , data
	code := helper2.ParseUint16(recv[2:])
	if code != 0 {
		return nil, fmt.Errorf("错误码: %d", code)
	}

	return recv[4:], nil
}

func (f *FinsHostLink) Write(station int, address protocol.Addr, values []byte) error {
	//构建写命令
	buf, e := buildWriteCommand(address, values)
	if e != nil {
		return e
	}

	//打包命令
	cmd := packAsciiCommand(&f.frame, buf)

	//发送请求
	recv, err := f.execute(cmd)
	if err != nil {
		return err
	}

	//[命令码 1 1] [结束码 0 0]
	code := helper2.ParseUint16(recv[2:])
	if code != 0 {
		return fmt.Errorf("错误码: %d", code)
	}

	return nil
}
