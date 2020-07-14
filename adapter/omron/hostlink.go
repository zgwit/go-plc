package omron

type HostLink struct {
	// 信息控制字段，默认0x80
	ICF byte // 0x80

	// PLC的节点地址，这个值在配置了ip地址之后是默认赋值的，默认为Ip地址的最后一位
	DA1 byte // 0x13

	// PLC的单元号地址
	DA2 byte // 0x00

	// 上位机的节点地址，假如你的电脑的Ip地址为192.168.0.13，那么这个值就是13
	SA1 byte

	// 上位机的单元号地址
	SA2 byte

	// 设备的标识号
	SID byte // 0x00
}

