package omron

import (
	"fmt"
	"iot-master/model"
	"iot-master/protocols/protocol"
	"regexp"
	"strconv"
)

type Command struct {
	BitCode  byte
	WordCode byte
}

var commands = map[string]Command{
	//DM Area
	"D": {0x02, 0x82},
	//CIO Area
	"C": {0x30, 0xB0},
	//Work Area
	"W": {0x31, 0xB1},
	//Holding Bits Area
	"H": {0x32, 0xB2},
	//Auxiliary Bits Area
	"A": {0x33, 0xB3},
}

type Address struct {
	Code   byte
	Offset uint16
	Bits   uint8
	IsBit  bool
}

func (a *Address) String() string {
	code := ""
	switch a.Code {
	case 0x02, 0x82:
		code = "D"
	case 0x30, 0xB0:
		code = "C"
	case 0x31, 0xB1:
		code = "W"
	case 0x32, 0xB2:
		code = "H"
	case 0x33, 0xB3:
		code = "A"
	}
	return code + strconv.Itoa(int(a.Offset))
}

func (a *Address) Lookup(data []byte, from protocol.Addr, tp model.DataType, le bool, precision int) (interface{}, bool) {
	base := from.(*Address)
	if base.Code != a.Code {
		return nil, false
	}
	if base.IsBit {
		cursor := int(a.Offset)*16 + int(a.Bits) - int(base.Offset)*16 - int(base.Bits)
		if cursor < 0 || cursor > len(data) {
			return nil, false
		}
		//TODO 此处应该明确数据格式
		return data[cursor] > 0, true
	} else {
		cursor := int(a.Offset-base.Offset) * 2
		if cursor < 0 || cursor > len(data) {
			return nil, false
		}
		val, err := tp.Decode(data[cursor:], le, precision)
		if err != nil {
			return nil, false
		}
		return val, true
	}
}

var addrRegexp *regexp.Regexp

func init() {
	addrRegexp = regexp.MustCompile(`^(D|C|W|H|A)(\d+)(.\d+)?$`)
}

func ParseAddress(name string, addr string) (protocol.Addr, error) {
	ss := addrRegexp.FindStringSubmatch(addr)
	if ss == nil || len(ss) < 3 {
		return nil, fmt.Errorf("欧姆龙地址格式错误: %s", addr)
	}
	var code uint8 = 1
	switch ss[1] {
	case "D":
		code = 0x82
	case "C":
		code = 0xB0
	case "W":
		code = 0xB1
	case "H":
		code = 0xB2
	case "A":
		code = 0xB3
	}
	offset, _ := strconv.ParseUint(ss[2], 10, 16)
	//offset, _ := strconv.Atoi(ss[2])
	address := &Address{
		Code:   code,
		Offset: uint16(offset),
	}
	if len(ss) == 4 {
		//TODO 判断bit，只能在 0~15之间
		bits, _ := strconv.ParseUint(ss[3][1:], 10, 8)
		address.Bits = uint8(bits)
		address.IsBit = true
		address.Code -= 0x80
	}
	return address, nil
}
