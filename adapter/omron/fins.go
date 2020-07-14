package omron

type Code byte

const (
	//DM Area
	DM Code = 0x02
	//CIO Area
	CIO = 0x30
	//Work Area
	WR = 0x31
	//Holding Bit Area
	HR = 0x32
	//Aux Bit Area
	AR = 0x33
)

type Fins struct {
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
