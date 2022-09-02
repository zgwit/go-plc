package protocol

import (
	"time"
)

type NewFunc func(conn io.ReadWriter, opts string) Protocol

// Protocol 协议接口
type Protocol interface {

	//Read 读数据
	Read(station int, area string, addr string, size int) ([]byte, error)

	//Write 写数据
	Write(station int, area string, addr string, data []byte) error

	//Attach(Conn io.ReadWriter) error
}

type transport interface {
	Send(request []byte, timeout time.Duration) (response []byte, err error)
}

func ask(conn io.ReadWriter, request []byte, response []byte) (int, error) {
	_, err := conn.Write(request)
	if err != nil {
		return 0, err
	}
	return conn.Read(response)
}
