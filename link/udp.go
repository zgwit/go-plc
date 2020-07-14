package link

import (
	"net"
)

type UDPLink struct {
	Address string
	Conn    net.Conn
}

func NewUDP() *UDPLink {
	link := &UDPLink{}
	return link
}

func (l *UDPLink) Open() (err error) {
	l.Conn, err = net.Dial("udp", l.Address)
	return
}

func (l *UDPLink) Read(buf []byte) (int, error) {
	return l.Conn.Read(buf)
}

func (l *UDPLink) Write(buf []byte) (int, error) {
	return l.Conn.Write(buf)
}

func (l *UDPLink) Close() error {
	return l.Conn.Close()
}
