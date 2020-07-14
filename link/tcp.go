package link

import (
	"net"
)

type TCPLink struct {
	Address string
	Conn    net.Conn
}

func NewTCP() *TCPLink {
	link := &TCPLink{}
	return link
}

func (l *TCPLink) Open() (err error) {
	l.Conn, err = net.Dial("tcp", l.Address)
	return
}

func (l *TCPLink) Read(buf []byte) (int, error) {
	//延迟100毫秒
	//l.Conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))

	return l.Conn.Read(buf)
}

func (l *TCPLink) Write(buf []byte) (int, error) {
	return l.Conn.Write(buf)
}

func (l *TCPLink) Close() error {
	return l.Conn.Close()
}
