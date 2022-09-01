package omron

import (
	"iot-master/protocols/protocol"
)

var Codes = []protocol.Code{
	{"D", "DM 区"},
	{"C", "CIO 区"},
	{"W", "Work 区"},
	{"H", "Holding Bits Area"},
	{"A", "Auxiliary Bits Area"},
}

var DescHostlink = protocol.Desc{
	Name:    "OmronHostlink",
	Version: "1.0",
	Label:   "Omron Hostlink",
	Codes:   Codes,
	Parser:  ParseAddress,
	Factory: NewHostLink,
}

var DescUDP = protocol.Desc{
	Name:    "OmronFinsUDP",
	Version: "1.0",
	Label:   "Omron FINS UDP",
	Codes:   Codes,
	Parser:  ParseAddress,
	Factory: NewFinsUDP,
}

var DescTCP = protocol.Desc{
	Name:    "OmronFinsTCP",
	Version: "1.0",
	Label:   "Omron FINS TCP",
	Codes:   Codes,
	Parser:  ParseAddress,
	Factory: NewFinsTCP,
}
