package modbus

import (
	"iot-master/protocols/protocol"
)

var Codes = []protocol.Code{
	{"C", "01 线圈"},
	{"D", "02 离散输入"},
	{"H", "03 保持寄存器"},
	{"I", "04 输入寄存器"},
}

var DescRTU = protocol.Desc{
	Name:    "ModbusRTU",
	Version: "1.0",
	Label:   "Modbus RTU",
	Codes:   Codes,
	Parser:  ParseAddress,
	Factory: NewRTU,
}

var DescTCP = protocol.Desc{
	Name:    "ModbusTCP",
	Version: "1.1",
	Label:   "Modbus TCP",
	Codes:   Codes,
	Parser:  ParseAddress,
	Factory: NewTCP,
}
