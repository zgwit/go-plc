package modbus

import (
	"github.com/zgwit/go-plc/protocol"
	"strconv"
)

type Address struct {
	Slave  uint8  `json:"slave"`
	Code   uint8  `json:"code"`
	Offset uint16 `json:"offset"`
}

func (a *Address) String() string {
	code := ""
	switch a.Code {
	case 1:
		code = "C"
	case 2:
		code = "D"
	case 3:
		code = "H"
	case 4:
		code = "I"
	}
	return code + strconv.Itoa(int(a.Offset))
}

func (a *Address) Diff(from protocol.Addr) (int, bool) {
	base := from.(*Address)
	if base.Code != a.Code {
		return 0, false
	}
	cursor := int(a.Offset - base.Offset)
	//Modbus是以双字
	if a.Code == FuncCodeReadHoldingRegisters || a.Code == FuncCodeReadInputRegisters {
		cursor *= 2
	}

	if cursor < 0 {
		return 0, false
	}
	return cursor, true
}

func ResolveAddress(name string, addr string) (protocol.Addr, error) {
	var code uint8 = 1
	switch name {
	case "C":
		code = 1
	case "D", "DI":
		code = 2
	case "H":
		code = 3
	case "I":
		code = 4
	}
	offset, err := strconv.ParseUint(addr, 10, 16)
	if err != nil {
		return nil, err
	}
	//offset, _ := strconv.Atoi(ss[2])

	return &Address{
		Code:   code,
		Offset: uint16(offset),
	}, nil
}
