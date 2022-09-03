package modbus

import "github.com/zgwit/go-plc/protocol"

var Areas = []protocol.Area{
	{"C", "01 线圈"},
	{"D", "02 离散输入"},
	{"H", "03 保持寄存器"},
	{"I", "04 输入寄存器"},
}

var DescRTU = protocol.Manifest{
	Name:     "ModbusRTU",
	Version:  "1.0",
	Label:    "Modbus RTU",
	Areas:    Areas,
	Resolver: ResolveAddress,
	Factory:  NewRTU,
}

var DescTCP = protocol.Manifest{
	Name:     "ModbusTCP",
	Version:  "1.1",
	Label:    "Modbus TCP",
	Areas:    Areas,
	Resolver: ResolveAddress,
	Factory:  NewTCP,
}
