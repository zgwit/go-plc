package protocol

import (
	"io"
	"time"
)

type CreateFactory func(conn io.ReadWriter, opts string) Protocol

// Protocol 协议接口
type Protocol interface {

	//Read 读数据
	Read(add Addr, size int) ([]byte, error)

	//Write 写数据
	Write(add Addr, data []byte) error

	//Attach(Conn io.ReadWriter) error
}

type transport interface {
	Send(request []byte, timeout time.Duration) (response []byte, err error)
}
