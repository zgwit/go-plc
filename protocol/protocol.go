package protocol

//Protocol 协议接口
type Protocol interface {
	Desc() *Desc

	//HandShake 握手
	//HandShake() error

	//Write 写数据
	Write(station int, addr Addr, data []byte) error

	//Read 读数据
	Read(station int, addr Addr, size int) ([]byte, error)

	//Poll 轮询读
	Poll(station int, addr Addr, size int) ([]byte, error)
}
