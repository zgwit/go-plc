package mitsubishi

import (
	"errors"
	"strconv"
	"strings"
)

// Command 命令
type Command struct {
	Code  byte
	IsBit bool
	Base  int
	//Code uint64
}

type Address struct {
	Command
	Name string
	Addr uint64
}

var commands = map[string]Command{
	"X":  {0x9C, true, 16},  //X输入继电器
	"Y":  {0x9D, true, 16},  //Y输出继电器
	"M":  {0x90, true, 10},  //M中间继电器
	"D":  {0xA8, false, 10}, //D数据寄存器
	"W":  {0xB4, false, 16}, //D数据寄存器
	"L":  {0x92, true, 10},  //L锁存继电器
	"F":  {0x93, true, 10},  //F报警器
	"V":  {0x94, false, 10}, //V边沿继电器
	"B":  {0xA0, true, 16},  //B链接继电器
	"R":  {0xAF, false, 10}, //R文件寄存器
	"S":  {0x98, true, 10},  //S步进继电器
	"Z":  {0xCC, false, 10}, //变址寄存器
	"ZR": {0xB0, false, 16}, //文件寄存器ZR区
	"TC": {0xC0, true, 10},  //定时器的线圈
	"TS": {0xC1, true, 10},  //定时器的触点
	"TN": {0xC2, false, 10}, //定时器的当前值
	"CC": {0xC3, true, 16},  //计数器的线圈
	"CS": {0xC4, true, 10},  //计数器的触点
	"CN": {0xC5, false, 10}, //计数器的当前值
	"SC": {0xC6, true, 10},  //累计定时器的线圈
	"SS": {0xC7, false, 10}, //累计定时器的触点
	"SN": {0xC8, false, 10}, //累计定时器的当前值
}

var a1eCommands = map[string]Command{
	"X": {'X', true, 8},   //X输入继电器
	"Y": {'Y', true, 8},   //Y输出继电器
	"M": {'M', true, 10},  //M中间继电器
	"D": {'D', false, 10}, //D数据寄存器
	"R": {'R', false, 10}, //R文件寄存器
	"S": {'S', true, 10},  //S步进继电器
}

func ParseAddress(address string) (addr Address, err error) {

	//先检查两字节
	k := strings.ToUpper(address[:2])
	if cmd, ok := commands[k]; ok {
		addr.Command = cmd
		addr.Name = k
		addr.Addr, err = strconv.ParseUint(address[2:], cmd.Base, 16)
		return
	}

	//检测单字节
	k = strings.ToUpper(address[:1])
	if cmd, ok := commands[k]; ok {
		addr.Command = cmd
		addr.Name = k + "*"
		addr.Addr, err = strconv.ParseUint(address[1:], cmd.Base, 16)
		return
	}

	err = errors.New("未知消息")
	return
}

func ParseA1EAddress(address string) (addr Address, err error) {
	//检测单字节
	k := strings.ToUpper(address[:1])
	if cmd, ok := a1eCommands[k]; ok {
		addr.Command = cmd
		addr.Name = k + "*"
		addr.Addr, err = strconv.ParseUint(address[1:], cmd.Base, 16)
		return
	}

	err = errors.New("未知消息")
	return
}
