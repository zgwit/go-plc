package omron

import (
	"regexp"
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

var addrRegexp *regexp.Regexp

func init() {
	addrRegexp = regexp.MustCompile(`^(D|C|W|H|A)(\d+)(.\d+)?$`)
}

func parseCode(area string) uint8 {
	var code uint8 = 1
	switch area {
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
	return code
}

type Addr struct {
}

func parseAddr(area string, addr string) (*Addr, error) {

	return nil, nil
}
