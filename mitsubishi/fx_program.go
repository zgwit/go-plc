package mitsubishi

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/zgwit/go-plc/helper"
	"github.com/zgwit/go-plc/protocol"
	"strconv"
	"strings"
)

type FxProgramCommand struct {
	Code  uint16
	IsBit bool
	Base  int
}

type FxProgramAddress struct {
	Code string
	Addr uint16
}

func (a *FxProgramAddress) String() string {
	return fmt.Sprintf("%s %d", a.Code, a.Addr)
}

func (a *FxProgramAddress) Diff(from protocol.Addr) (int, bool) {

	return 0, false
}

var fxProgramCommands = map[string]FxProgramCommand{
	"X":  {0x0080, true, 8},   //X输入继电器
	"Y":  {0x00A0, true, 8},   //Y输出继电器
	"M":  {0x0100, true, 10},  //M中间继电器
	"D":  {0x1000, false, 10}, //D数据寄存器
	"S":  {0x0000, true, 10},  //S步进继电器
	"TS": {0x00C0, true, 10},  //定时器的触点
	"TC": {0x02C0, true, 10},  //定时器的线圈
	"TN": {0x0800, false, 10}, //定时器的当前值 ?
	"CS": {0x01C0, true, 10},  //计数器的触点
	"CC": {0x03C0, true, 10},  //计数器的线圈
	"CN": {0x0A00, false, 10}, //计数器的当前值 ?
}

func ParseFxProgramAddress(code string, address string) (protocol.Addr, error) {
	var addr FxProgramAddress

	cmd, ok := fxProgramCommands[code]
	if !ok {
		return nil, fmt.Errorf("不支持的区域 %s", code)
	}
	addr.Code = code
	v, err := strconv.ParseUint(address[2:], cmd.Base, 16)
	if cmd.IsBit {
		addr.Addr = cmd.Code + uint16(int(v)/8)
	} else {
		addr.Addr = cmd.Code + uint16(v)*2
	}
	return &addr, err
}

// FxProgram FX协议
type FxProgram struct {
	link protocol.Messenger
}

// NewFxProgram 新建
func NewFxProgram() *FxProgram {
	return &FxProgram{}
}

// Read 解析
func (t *FxProgram) Read(address protocol.Addr, length int) ([]byte, error) {
	addr := address.(*FxProgramAddress)

	buf := make([]byte, 11)
	buf[0] = 0x02                                // STX
	buf[1] = 0x30                                // 命令 Read
	helper.WriteUint16Hex(buf[2:], addr.Addr)    // 偏移地址
	helper.WriteUint8Hex(buf[6:], uint8(length)) // 读取长度
	buf[8] = 0x03                                // ETX

	// 计算和校验
	var sum uint8 = 0
	for i := 1; i < len(buf)-2; i++ {
		sum += buf[i]
	}

	//最后两位是校验
	helper.WriteUint8Hex(buf[len(buf)-2:], sum)

	fmt.Println("FxProgram read buff = ", hex.EncodeToString(buf))
	var recv [256]byte
	n, err := t.link.Ask(buf, recv[:])
	fmt.Println("FxProgram recv buff", hex.EncodeToString(recv[:n]))
	if err != nil {
		return nil, err
	}

	if recv[0] == 0x15 {
		return nil, errors.New("返回错误")
	}

	ret, err := hex.DecodeString(string(recv[1 : n-3]))

	if err != nil {
		return nil, err
	}
	//ret := helper.FromHex(recv[1 : len(recv)-3])

	return ret, nil
}

// Write 写
func (t *FxProgram) Write(address protocol.Addr, values []byte) error {
	addr := address.(*FxProgramAddress)

	//先转成十六进制
	length := len(values)

	values = []byte(strings.ToUpper(hex.EncodeToString(values)))

	buf := make([]byte, 11+(length*2))
	buf[0] = 0x02                                // STX
	buf[1] = 0x31                                // 命令 Write
	helper.WriteUint16Hex(buf[2:], addr.Addr)    // 偏移地址
	helper.WriteUint8Hex(buf[6:], uint8(length)) // 写入长度
	copy(buf[8:], values)                        // 写入内容
	buf[len(buf)-3] = 0x03                       // ETX

	// 计算和校验
	var sum uint8 = 0
	for i := 1; i < len(buf)-2; i++ {
		sum += buf[i]
	}
	//最后两位是校验
	helper.WriteUint8Hex(buf[len(buf)-2:], sum)

	fmt.Println("FxProgram write buff = ", hex.EncodeToString(buf))
	var recv [256]byte
	n, err := t.link.Ask(buf, recv[:])
	fmt.Println("FxProgram recv buff", hex.EncodeToString(recv[:n]))
	if err != nil {
		return err
	}
	if recv[0] == 0x15 {
		return errors.New("错误")
	} else {
		return nil
	}
}
