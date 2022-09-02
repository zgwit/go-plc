package protocol

import (
	"net"
	"time"
)

type Socket struct {
	Conn    net.Conn
	Timeout time.Duration
}

func (s *Socket) Read(b []byte) (int, error) {
	err := s.Conn.SetReadDeadline(time.Now().Add(s.Timeout))
	if err != nil {
		return 0, err
	}
	return s.Conn.Read(b)
}

func (s *Socket) Write(b []byte) (int, error) {
	err := s.Conn.SetWriteDeadline(time.Now().Add(s.Timeout))
	if err != nil {
		return 0, err
	}
	return s.Conn.Read(b)
}

func (s *Socket) Close() error {
	return s.Conn.Close()
}
