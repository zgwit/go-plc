package mitsubishi

import "iot-master/connect"

// A1eAdapter A1E协议
type A1eAdapter struct {
	PlcNumber byte
	link      connect.Tunnel
}

func NewA1eAdapter() *A1eAdapter {
	a := A1eAdapter{}
	a.PlcNumber = 0xFF
	return &a
}

//Read 读取数据
func (t *A1eAdapter) Read(address string, length int) ([]byte, error) {

	//解析地址
	addr, e := ParseA1EAddress(address)
	if e != nil {
		return nil, e
	}

	//副标题
	var subTitle byte = 0x01
	if addr.IsBit {
		subTitle = 0x00
	}

	//构建命令
	buf := make([]byte, 12)
	buf[0] = subTitle        //副标题，0：bit 1：byte
	buf[1] = t.PlcNumber     //PLC号
	buf[2] = 0x0A            // CPU监视定时器（L）这里设置为0x00,0x0A，等待CPU返回的时间为10*250ms=2.5秒
	buf[3] = 0x00            // CPU监视定时器（H）
	buf[4] = byte(addr.Addr) // 起始软元件（开始读取的地址）
	buf[5] = byte(addr.Addr >> 8)
	buf[6] = byte(addr.Addr >> 16)
	buf[7] = byte(addr.Addr >> 24)
	buf[8] = 0x20          // 软元件代码（L）
	buf[9] = addr.Code     // 软元件代码（H）
	buf[10] = byte(length) // 软元件点数
	buf[11] = 0x00

	//发送命令
	if err := t.link.Write(buf); err != nil {
		return nil, err
	}

	//接收响应
	recv := make([]byte, 2+length)
	//if _, err := t.link.Read(recv); err != nil {
	//	return nil, err
	//}

	// 80/81 00 ....
	return recv[2:], nil
}

//Write 写入数据
func (t *A1eAdapter) Write(address string, values []byte) error {

	length := len(values)

	//解析地址
	addr, e := ParseAddress(address)
	if e != nil {
		return e
	}

	//副标题
	var subTitle byte = 0x03
	if addr.IsBit {
		subTitle = 0x02
	}

	//构建命令
	buf := make([]byte, 12+length)
	buf[0] = subTitle        //副标题，2：bit 3：byte
	buf[1] = t.PlcNumber     //PLC号
	buf[2] = 0x0A            // CPU监视定时器（L）这里设置为0x00,0x0A，等待CPU返回的时间为10*250ms=2.5秒
	buf[3] = 0x00            // CPU监视定时器（H）
	buf[4] = byte(addr.Addr) // 起始软元件（开始读取的地址）
	buf[5] = byte(addr.Addr >> 8)
	buf[6] = byte(addr.Addr >> 16)
	buf[7] = byte(addr.Addr >> 24)
	buf[8] = 0x20          // 软元件代码（L）
	buf[9] = addr.Code     // 软元件代码（H）
	buf[10] = byte(length) // 软元件点数
	buf[11] = 0x00

	//附加数据
	copy(buf[12:], values)

	//发送命令
	if err := t.link.Write(buf); err != nil {
		return err
	}

	//接收响应
	//recv := make([]byte, 2)
	//if _, err := t.link.Read(recv); err != nil {
	//	return err
	//}

	// 82/83 00
	return nil
}
