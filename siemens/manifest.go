package siemens

import (
	"github.com/zgwit/go-plc/protocol"
	"io"
)

var Codes = []protocol.Area{
	{"I", "输入"},
	{"Q", "输出"},
	{"M", "内部标志"},
	{"DB", "数据块"},
	{"DI", "背景数据块"},
	{"L", "局部变量"},
	{"V", "全局变量"},
	{"C", "计数器"},
	{"T", "定时器"},
}

var DescS7_200 = protocol.Manifest{
	Name:     "S7-S7-200",
	Version:  "1.0",
	Label:    "S7 S7-200",
	Areas:    Codes,
	Resolver: ParseAddress,
	Factory: func(link io.ReadWriter, opts string) protocol.Protocol {
		return &S7{
			handshake1: parseHex(handshake1_200),
			handshake2: parseHex(handshake2_200),
			link:       protocol.Messenger{Conn: link},
		}
	},
}

var DescS7_200_Smart = protocol.Manifest{
	Name:     "S7-S7-200-Smart",
	Version:  "1.0",
	Label:    "S7 S7-200 Smart",
	Areas:    Codes,
	Resolver: ParseAddress,
	Factory: func(link io.ReadWriter, opts string) protocol.Protocol {
		s7 := &S7{
			handshake1: parseHex(handshake1_200_smart),
			handshake2: parseHex(handshake2_200_smart),
			link:       protocol.Messenger{Conn: link},
		}
		return s7
	},
}

var DescS7_300 = protocol.Manifest{
	Name:     "S7-S7-300",
	Version:  "1.0",
	Label:    "S7 S7-300",
	Areas:    Codes,
	Resolver: ParseAddress,
	Factory: func(link io.ReadWriter, opts string) protocol.Protocol {
		return &S7{
			handshake1: parseHex(handshake1_300),
			handshake2: parseHex(handshake2_300),
			link:       protocol.Messenger{Conn: link},
		}
	},
}

var DescS7_400 = protocol.Manifest{
	Name:     "S7-S7-400",
	Version:  "1.0",
	Label:    "S7 S7-400",
	Areas:    Codes,
	Resolver: ParseAddress,
	Factory: func(link io.ReadWriter, opts string) protocol.Protocol {
		//TODO 设置机架和槽号
		//setRackSlot()
		return &S7{
			handshake1: parseHex(handshake1_400),
			handshake2: parseHex(handshake2_400),
			link:       protocol.Messenger{Conn: link},
		}
	},
}

var DescS7_1200 = protocol.Manifest{
	Name:     "S7-S7-1200",
	Version:  "1.0",
	Label:    "S7 S7-1200",
	Areas:    Codes,
	Resolver: ParseAddress,
	Factory: func(link io.ReadWriter, opts string) protocol.Protocol {
		return &S7{
			handshake1: parseHex(handshake1_1200),
			handshake2: parseHex(handshake2_1200),
			link:       protocol.Messenger{Conn: link},
		}
	},
}

var DescS7_1500 = protocol.Manifest{
	Name:     "S7-S7-1500",
	Version:  "1.0",
	Label:    "S7 S7-1500",
	Areas:    Codes,
	Resolver: ParseAddress,
	Factory: func(link io.ReadWriter, opts string) protocol.Protocol {
		return &S7{
			handshake1: parseHex(handshake1_1500),
			handshake2: parseHex(handshake2_1500),
			link:       protocol.Messenger{Conn: link},
		}
	},
}
