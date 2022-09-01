package mitsubishi

import (
	"fmt"
	"iot-master/connect"
	"iot-master/helper"
	"strconv"
)

// A3C1 A3C1协议
type A3C1 struct {
	StationNumber byte //站编号
	NetworkNumber byte //网络编号
	PlcNumber     byte //PLC编号
	UpperNumber   byte //上位机编号

	link connect.Tunnel
}

func NewA3C1() *A3C1 {
	a := A3C1{}
	a.StationNumber = 0
	a.NetworkNumber = 0
	a.PlcNumber = 0xFF
	a.UpperNumber = 0
	return &a

}

func (t *A3C1) BuildCommand(cmd []byte) []byte {
	length := len(cmd)

	buf := make([]byte, 13+length)
	buf[0] = 0x05                                 //ENQ
	buf[1] = 0x46                                 //F 帧识别号H
	buf[2] = 0x39                                 //9 帧识别号L
	helper.WriteByteHex(buf[3:], t.StationNumber) //站编号
	helper.WriteByteHex(buf[5:], t.NetworkNumber) //网络编号
	helper.WriteByteHex(buf[7:], t.PlcNumber)     //PLC编号
	helper.WriteByteHex(buf[9:], t.UpperNumber)   //上位机编号

	//命令
	copy(buf[11:], cmd)

	// 计算和校验
	var sum byte = 0
	for i := 1; i < len(buf)-3; i++ {
		sum += buf[i]
	}

	//最后两位是校验
	helper.WriteByteHex(buf[len(buf)-2:], sum)

	return buf

}

func (t *A3C1) Read(address string, length int) ([]byte, error) {

	//解析地址
	addr, e := ParseAddress(address)
	if e != nil {
		return nil, e
	}

	buf := make([]byte, 20)
	copy(buf, []byte("0401")) //命令 读取
	if addr.IsBit {
		copy(buf[4:], []byte("0001")) //子命令 按位
	} else {
		copy(buf[4:], []byte("0000")) //子命令 按字
	}
	copy(buf[8:], addr.Name) // 软元件类型
	if addr.Base == 10 {
		copy(buf[10:], fmt.Sprintf("%d6", addr.Addr)) // 起始地址的地位
	} else {
		copy(buf[10:], fmt.Sprintf("%X6", addr.Addr)) // 起始地址的地位
	}
	copy(buf[16:], fmt.Sprintf("%X4", length)) // 软元件点数

	//构建命令
	cmd := t.BuildCommand(buf)

	//发送命令
	if err := t.link.Write(cmd); err != nil {
		return nil, err
	}

	//如果不是位，需要纠正长度
	if !addr.IsBit {
		length = length * 4
	}

	//接收响应
	recv := make([]byte, 11+length)
	//if _, err := t.link.Read(recv); err != nil {
	//	return nil, err
	//}

	// 正确
	// STX  帧识别号 站编号 网络编号 PLC编号 上位站编号 ---内容--- ETX 和校验
	// 0x20  F 9  0 0  0 0  F F  0 0  ....  0x03 0 0
	data := recv[11:length]

	if addr.IsBit {
		//0x31 => 0x01, 0x30 => 0
		r := make([]byte, length)
		for i, v := range data {
			if v == '1' {
				r[i] = 1
			} else {
				r[i] = 0
			}
		}
	} else {
		//每4字节，表示一个十六进制
		r := make([]byte, length/2)
		for i := 0; i < length; i += 4 {
			d, _ := strconv.ParseUint(string(data[i:i+4]), 16, 32)
			//TODO 大小端需要再去确认
			r[i*2] = byte(d << 8)
			r[i*2+1] = byte(d)
		}
	}

	// 错误
	// NAK 帧识别号 站编号 网络编号 PLC编号 上位站编号 ---内容--- 错误码
	// 0x15  F 9  0 0  0 0  F F  0 0  0 0 0 0

	return buf, nil
}

func (t *A3C1) Write(address string, values []byte) error {
	//解析地址
	addr, e := ParseAddress(address)
	if e != nil {
		return e
	}

	length := len(values)

	var value []byte

	//数据预处理，位数据要转成0和1，字数据
	if addr.IsBit {
		value = make([]byte, length)
		for k, v := range values {
			if v != 0 {
				value[k] = '1'
			} else {
				value[k] = '0'
			}
		}
	} else {
		//uint16 数组转字符串
		value := make([]byte, length*2)
		for i := 0; i < length/2; i++ {
			d := values[i]<<8 + values[i+1]
			if addr.Base == 10 {
				copy(value[i*4:], fmt.Sprintf("%d4", d))
			} else {
				copy(value[i*4:], fmt.Sprintf("%X4", d))
			}
		}

		length = len(value) //length * 2
	}

	buf := make([]byte, 20+length)
	copy(buf, []byte("1401")) //命令 读取
	if addr.IsBit {
		copy(buf[4:], []byte("0001")) //子命令 按位
	} else {
		copy(buf[4:], []byte("0000")) //子命令 按字
	}
	copy(buf[8:], addr.Name) // 软元件类型
	if addr.Base == 10 {
		copy(buf[10:], fmt.Sprintf("%d6", addr.Addr)) // 起始地址的地位
	} else {
		copy(buf[10:], fmt.Sprintf("%X6", addr.Addr)) // 起始地址的地位
	}
	if addr.IsBit {
		copy(buf[16:], fmt.Sprintf("%X4", length)) // 软元件点数
	} else {
		copy(buf[16:], fmt.Sprintf("%X4", length/4)) // 软元件点数
	}

	//附加数据
	copy(buf[20:], value)

	//构建命令
	cmd := t.BuildCommand(buf)

	//发送命令
	if err := t.link.Write(cmd); err != nil {
		return err
	}

	//接收响应
	//recv := make([]byte, 15)
	//if _, err := t.link.Read(recv); err != nil {
	//	return err
	//}

	// 正确
	// ACK  帧识别号 站编号 网络编号 PLC编号 上位站编号
	// 0x06  F 9  0 0  0 0  F F  0 0

	// 错误
	// NAK 帧识别号 站编号 网络编号 PLC编号 上位站编号  错误码
	// 0x15  F 9  0 0  0 0  F F  0 0  0 0 0 0

	return nil
}
