package mitsubishi

import (
	"errors"
	"fmt"
	helper2 "github.com/zgwit/go-plc/helper"
	"io"
	"strconv"
	"strings"
)

type fxSpecialCommand struct {
	IsBit bool
	Base  int
}

type fxSpecialAddr struct {
	fxSpecialCommand
	Address string
}

var fxSpecialCommands = map[string]fxSpecialCommand{
	"X":  {true, 8},   //X输入继电器
	"Y":  {true, 8},   //Y输出继电器
	"M":  {true, 10},  //M中间继电器
	"D":  {false, 10}, //D数据寄存器
	"R":  {false, 10}, //R文件寄存器
	"S":  {true, 10},  //S步进继电器
	"TS": {true, 10},  //定时器的触点
	"TN": {false, 10}, //定时器的当前值
	"CS": {true, 10},  //计数器的触点
	"CN": {false, 10}, //计数器的当前值
}

func parseFxSpecialAddress(address string) (addr fxSpecialAddr, err error) {
	var v uint64

	//先检查两字节
	k := strings.ToUpper(address[:2])
	if cmd, ok := fxSpecialCommands[k]; ok {
		v, err = strconv.ParseUint(address[2:], cmd.Base, 16)
		if err == nil {
			addr.fxSpecialCommand = cmd
			addr.Address = k + fmt.Sprintf("%03d", v)
		}
		return
	}

	//检测单字节
	k = strings.ToUpper(address[:1])
	if cmd, ok := fxSpecialCommands[k]; ok {
		v, err = strconv.ParseUint(address[1:], cmd.Base, 16)
		if err == nil {
			addr.fxSpecialCommand = cmd
			addr.Address = k + fmt.Sprintf("%04d", v)
		}
		return
	}

	err = errors.New("未知消息")
	return
}

// FxSpecial FX协议
type FxSpecial struct {
	StationNumber byte //站编号 00
	PlcNumber     byte //PLC编号 FF

	CheckSum bool //默认true
	Delay    uint8

	link io.ReadWriter
}

func NewFxSpecial() *FxSpecial {
	return &FxSpecial{}
}

func (t *FxSpecial) Read(address string, length int) ([]byte, error) {

	recvLength := length

	addr, err := parseFxSpecialAddress(address)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 23)
	buf[0] = 0x05                                   //ENQ
	helper2.WriteUint8Hex(buf[1:], t.StationNumber) //站号
	helper2.WriteByteHex(buf[3:], t.PlcNumber)      //PLC号
	if addr.IsBit {
		buf[10] = 'B' //位
	} else {
		buf[10] = 'W'                //字
		recvLength = recvLength << 1 // recvLength*2 字是双字节
	}
	buf[11] = 'R'                               //读取
	helper2.WriteUint8Hex(buf[12:], t.Delay)    //延迟
	copy(buf[14:], addr.Address)                // 偏移地址
	copy(buf[19:], fmt.Sprintf("%02d", length)) //读取长度

	if t.CheckSum {
		//最后两位是和校验
		helper2.WriteUint8Hex(buf[len(buf)-2:], helper2.Sum(buf[1:len(buf)-2]))
	} else {
		buf = buf[:21]
	}

	//发送请求

	recv := make([]byte, recvLength+8)
	//_, err = t.link.Read(recv)
	//if err != nil {
	//	return nil, err
	//}

	//STX 站号 PLC号 数据 ETX 和校验
	ret := recv[5 : len(recv)-3]
	if addr.IsBit {
		//布尔数组
		ret = helper2.AsciiToBool(ret)
	} else {
		//转十六进制
		ret = helper2.FromHex(ret)
	}
	return ret, nil
}

func (t *FxSpecial) Write(address string, values []byte) error {

	addr, err := parseFxSpecialAddress(address)
	if err != nil {
		return err
	}

	if addr.IsBit {
		//布尔数组
		values = helper2.BoolToAscii(values)
	} else {
		//转十六进制
		values = helper2.ToHex(values)
	}

	length := len(values)

	buf := make([]byte, 23+length)
	buf[0] = 0x05                                   //ENQ
	helper2.WriteUint8Hex(buf[1:], t.StationNumber) //站号
	helper2.WriteByteHex(buf[3:], t.PlcNumber)      //PLC号
	if addr.IsBit {
		buf[10] = 'B' //位
	} else {
		buf[10] = 'W'        //字
		length = length >> 1 // length/2 字是双字节
	}
	buf[11] = 'W'                               //写入
	helper2.WriteUint8Hex(buf[12:], t.Delay)    //延迟
	copy(buf[14:], addr.Address)                // 偏移地址
	copy(buf[19:], fmt.Sprintf("%02d", length)) //读取长度
	copy(buf[21:], values)                      //写入数据

	if t.CheckSum {
		//最后两位是和校验
		helper2.WriteUint8Hex(buf[len(buf)-2:], helper2.Sum(buf[1:len(buf)-2]))
	} else {
		buf = buf[:21+length]
	}

	//发送请求

	//recv := make([]byte, 5)
	//_, err = t.link.Read(recv)
	//if err != nil {
	//	return err
	//}

	//ACK 站号 PLC号
	return nil
}
