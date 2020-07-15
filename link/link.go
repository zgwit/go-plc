package link

//Link 通讯链路
type Link interface {
	//关闭
	Close() error

	//发送请求
	Request(req []byte) ([]byte, error)
}
