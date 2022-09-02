package mitsubishi

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/zgwit/go-plc/helper"
	"github.com/zgwit/go-plc/protocol"
	"strings"
)

// FxProgram FX协议
type FxProgram struct {
	link protocol.Messenger
}

func (t *FxProgram) Write(station int, area string, addr string, data []byte) error {
	a, err := ParseFxProgramAddress(area, addr)
	if err != nil {
		return err
	}
	return t.write(a, data)
}

func (t *FxProgram) Read(station int, area string, addr string, size int) ([]byte, error) {
	a, err := ParseFxProgramAddress(area, addr)
	if err != nil {
		return nil, err
	}
	return t.read(a, size)
}

// NewFxSerial 新建
func NewFxSerial() *FxProgram {
	return &FxProgram{}
}

// Read 解析
func (t *FxProgram) read(addr *FxProgramAddress, length int) ([]byte, error) {
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
func (t *FxProgram) write(addr *FxProgramAddress, values []byte) error {

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
