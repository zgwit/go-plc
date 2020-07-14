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

func (l *SerialLink) Read(buf []byte) (int, error) {
	var sum int = 0
	for {
		n, e := l.port.Read(buf[sum:])
		if e != nil {
			return 0, e
		}
		sum += n
		//没有了
		if n == 0 {
			return sum, nil
		}
		//收满了
		if sum == len(buf) {
			return  sum, nil
		}
	}
}

func (l *SerialLink) Write(buf []byte) (int, error) {
	return l.port.Write(buf)
}

func (l *SerialLink) Close() error {
	return l.port.Close()
}
