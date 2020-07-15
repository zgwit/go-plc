package link

import (
	"github.com/tarm/serial"
	"time"
)

type SerialLink struct {
	Name     string
	Baud     int
	Size     byte
	Parity   byte
	StopBits byte

	ReadTimeout time.Duration

	port *serial.Port
}

func NewSerial() *SerialLink {
	link := &SerialLink{
		Baud:     115200,
		Size:     7,
		Parity:   'N',
		StopBits: 1,
	}
	return link
}

func (l *SerialLink) Open() error {
	config := &serial.Config{
		Name: l.Name,
		Baud: l.Baud,
		ReadTimeout: time.Millisecond * 500,
		Size:     l.Size,
		Parity:   serial.Parity(l.Parity),
		StopBits: serial.StopBits(l.StopBits),
	}
	port, err := serial.OpenPort(config)
	if err != nil {
		return err
	}
	l.port = port
	return nil
}

func (l *SerialLink) Request(req []byte) ([]byte, error) {
	//fmt.Println("serial send", req)

	//发送
	n, e := l.port.Write(req)
	if e != nil {
		return nil, e
	}
	if n < len(req) {
		//TODO 此处应该继续发送，直到发送完
		_, _ = l.port.Write(req[n:])
	}

	buf := make([]byte, 1024)
	//读取
	var sum int = 0
	for {
		n, e := l.port.Read(buf[sum:])
		if e != nil {
			return nil, e
		}
		//没有了
		if n == 0 {
			//return buf[:sum], nil
			break
		}

		sum += n
		//收满了
		if sum == len(buf) {
			//return  buf, nil
			break //TODO 此处应该扩张
		}
	}

	//fmt.Println("serial recv", buf[:sum])

	return buf[:sum], nil
}

func (l *SerialLink) Close() error {
	return l.port.Close()
}
