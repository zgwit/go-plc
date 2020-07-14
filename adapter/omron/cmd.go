package omron

import "github.com/zgwit/go-plc/helper"

func buildReadBitCommand(code Code, addr uint16, bit uint8, length uint16) []byte {
	buf := make([]byte, 8)
	//命令
	buf[0] = 0x01 //MRC 读取存储区数据
	buf[1] = 0x01 //SRC
	buf[2] = byte(code) + 0x80
	helper.WriteUint16(buf[3:], addr)   // 地址
	buf[5] = bit                        // 位地址
	helper.WriteUint16(buf[6:], length) // 长度

	return buf
}

func buildReadWordCommand(code Code, addr uint16, length uint16) []byte {
	buf := make([]byte, 8)
	//命令
	buf[0] = 0x01 //MRC 读取存储区数据
	buf[1] = 0x01 //SRC
	buf[2] = byte(code)
	helper.WriteUint16(buf[3:], addr)   // 地址
	buf[5] = 0                          // 位地址
	helper.WriteUint16(buf[6:], length) // 长度

	return buf
}

func buildWriteBitCommand(code Code, addr uint16, bit byte, values []byte) []byte {
	length := len(values)

	buf := make([]byte, 8+length)
	buf[0] = 0x01 //MRC 读取存储区数据
	buf[1] = 0x02 //SRC
	buf[2] = byte(code)
	helper.WriteUint16(buf[3:], addr) // 地址
	buf[5] = bit
	helper.WriteUint16(buf[6:], uint16(length)) // 长度
	copy(buf[8:], values)                       //数据
	return buf
}

func buildWriteWordCommand(code Code, addr uint16, values []byte) []byte {
	length := len(values)

	buf := make([]byte, 8+length)
	buf[0] = 0x01 //MRC 读取存储区数据
	buf[1] = 0x02 //SRC
	buf[2] = byte(code) + 0x80
	helper.WriteUint16(buf[3:], addr) // 地址
	buf[5] = 0
	helper.WriteUint16(buf[6:], uint16(length / 2)) // 长度 一个word是双字节
	copy(buf[8:], values)                       //数据
	return buf
}

func packTCPCommand(cmd uint32, payload []byte) []byte {
	length := len(payload)
	buf := make([]byte, 16+length)

	//copy(buf, "FINS")
	buf[0] = 0x46 //FINS
	buf[1] = 0x49
	buf[2] = 0x4e
	buf[3] = 0x53

	//长度
	helper.WriteUint32(buf[4:], uint32(length))

	//命令码 读写时为2
	helper.WriteUint32(buf[8:], uint32(cmd))

	//错误码
	helper.WriteUint32(buf[12:], 0)

	//附加数据
	copy(buf[16:], payload)

	return buf
}
