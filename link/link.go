package link

import "io"

//Link 通讯链路
type Link interface {
	//继承部分io接口，读、写、关闭
	io.ReadWriteCloser
}
