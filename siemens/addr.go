package siemens

import (
	"fmt"
	"github.com/zgwit/go-plc/protocol"
	"regexp"
	"strconv"
	"strings"
)

type Address struct {
	Code   byte
	DB     uint16
	Offset uint32
	HasBit bool
	Bits   uint8
}

func (a *Address) String() string {
	code := ""
	switch a.Code {
	case 0x81:
		code = "I"
	case 0x82:
		code = "Q"
	case 0x83:
		code = "M"
	case 0x84:
		code = "DB"
	case 0x85:
		code = "DI"
	case 0x86:
		code = "L"
	case 0x87:
		code = "V"
	case 0x1D:
		code = "T"
	case 0x1C:
		code = "C"
	}
	return code + strconv.Itoa(int(a.Offset))
}

func (a *Address) Diff(from protocol.Addr) (int, bool) {
	base := from.(*Address)
	if base.Code != a.Code || base.DB != a.DB {
		return 0, false
	}
	if base.HasBit {
		if a.HasBit {
			cursor := int(a.Offset-base.Offset)*8 + int(a.Bits) - int(base.Bits)
			if cursor < 0 {
				return 0, false
			}
			return cursor, true
		} else {
			return 0, false
		}
	} else {
		cursor := int(a.Offset - base.Offset)
		if cursor < 0 {
			return 0, false
		}
		if a.HasBit {
			//return data[cursor]&(0x01<<a.Bits) > 0, true
			return 0, false
		} else {
			return cursor, true
		}
	}
}

var addrRegexp1 *regexp.Regexp
var addrRegexp2 *regexp.Regexp
var addrRegexp3 *regexp.Regexp

func init() {
	addrRegexp1 = regexp.MustCompile(`^\d+$`)
	addrRegexp2 = regexp.MustCompile(`^\d+.\d+$`)
	addrRegexp3 = regexp.MustCompile(`^\d+(.\d+)?$`)
}

func ParseAddress(name, addr string) (protocol.Addr, error) {
	var code uint8 = 1
	switch name {
	case "I":
		code = 0x81
	case "Q":
		code = 0x82
	case "M":
		code = 0x83
	case "DB":
		code = 0x84
	case "DI":
		code = 0x85
	case "L":
		code = 0x86
	case "V":
		code = 0x87
	case "C":
		code = 0x1C
	case "T":
		code = 0x1D
	}
	//addrRegexp.MatchString(addr)

	if name == "DB" {
		if !addrRegexp2.MatchString(addr) {
			return nil, fmt.Errorf("%s 不是有效的DB格式", addr)
		}
		ss := strings.Split(addr, ".")
		db, _ := strconv.ParseUint(ss[0], 10, 16)
		offset, _ := strconv.ParseUint(ss[1], 10, 24)
		return &Address{
			Code:   0x84,
			DB:     uint16(db),
			Offset: uint32(offset),
		}, nil
	}
	if name == "T" || name == "C" {
		if !addrRegexp1.MatchString(addr) {
			return nil, fmt.Errorf("%s 不是有效的整数格式", addr)
		}
		offset, _ := strconv.ParseUint(addr, 10, 24)
		return &Address{
			Code:   code,
			Offset: uint32(offset),
		}, nil
	}

	if !addrRegexp3.MatchString(addr) {
		return nil, fmt.Errorf("%s 不是有效的地址格式", addr)
	}

	//i q m v
	address := &Address{Code: code}
	ss := strings.Split(addr, ".")
	offset, _ := strconv.ParseUint(ss[0], 10, 24)
	address.Offset = uint32(offset)
	if len(ss) > 1 {
		address.HasBit = true
		bits, _ := strconv.ParseUint(ss[1], 10, 16)
		address.Bits = uint8(bits)
	}
	return address, nil
}
