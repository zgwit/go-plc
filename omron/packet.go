package omron

import (
	"github.com/zgwit/go-plc/helper"
	"github.com/zgwit/go-plc/protocol"
)

func buildReadCommand(address protocol.Addr, length int) ([]byte, error) {
	//解析地址
	addr := address.(*Address)

	buf := make([]byte, 8)

	//命令
	buf[0] = 0x01 //MRC 读取存储区数据
	buf[1] = 0x01 //SRC
	buf[2] = addr.Code
	// 地址
	helper.WriteUint16(buf[3:], uint16(addr.Offset))
	// 位地址
	buf[5] = addr.Bits
	// 长度
	helper.WriteUint16(buf[6:], uint16(length))

	return buf, nil
}

func buildWriteCommand(address protocol.Addr, values []byte) ([]byte, error) {
	//解析地址
	addr := address.(*Address)

	length := len(values)

	buf := make([]byte, 8+length)

	buf[0] = 0x01 //MRC 读取存储区数据
	buf[1] = 0x02 //SRC
	buf[2] = addr.Code

	// 地址
	helper.WriteUint16(buf[3:], uint16(addr.Offset))
	buf[5] = addr.Bits

	if addr.IsBit {
		length = length / 2 // 一个word是双字节
	}
	// 长度
	helper.WriteUint16(buf[6:], uint16(length))

	//数据
	copy(buf[8:], values)

	return buf, nil
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

func packUDPCommand(uf *UdpFrame, payload []byte) []byte {
	length := len(payload)
	buf := make([]byte, 10+length)

	//UDP头
	buf[0] = uf.ICF
	buf[1] = uf.RSV
	buf[2] = uf.GCT
	buf[3] = uf.DNA
	buf[4] = uf.DA1
	buf[5] = uf.DA2
	buf[6] = uf.SNA
	buf[7] = uf.SA1
	buf[8] = uf.SA2
	buf[9] = uf.SID

	//附加数据
	copy(buf[10:], payload)

	return buf
}

func packAsciiCommand(uf *UdpFrame, payload []byte) []byte {
	cmd := helper.ToHex(payload)

	length := len(cmd)

	buf := make([]byte, 18+length)

	buf[0] = '@'
	helper.WriteByteHex(buf[1:], uf.DA1) //PLC设备号
	buf[3] = 'F'                         //识别码
	buf[4] = 'A'
	buf[5] = 0x30 //响应等待时间 x 15ms
	helper.WriteByteHex(buf[6:], uf.ICF)
	helper.WriteByteHex(buf[8:], uf.DA2)
	helper.WriteByteHex(buf[10:], uf.SA2)
	helper.WriteByteHex(buf[12:], uf.SID)
	copy(buf[14:], cmd)

	//计算FCS
	tmp := buf[0]
	for i := 1; i < length+14; i++ {
		tmp = tmp ^ buf[i]
	}
	helper.WriteByteHex(buf[length+14:], tmp)
	buf[length+16] = '*'
	buf[length+17] = 0x0D //CR

	return buf
}
