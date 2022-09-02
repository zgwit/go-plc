package protocol

import (
	"net"
	"time"
)

type Socket struct {
	conn    net.Conn
	timeout time.Duration
}

func (s *Socket) Read(b []byte) (int, error) {
	err := s.conn.SetReadDeadline(time.Now().Add(s.timeout))
	if err != nil {
		return 0, err
	}
	return s.conn.Read(b)
}

func (s *Socket) Write(b []byte) (int, error) {
	err := s.conn.SetWriteDeadline(time.Now().Add(s.timeout))
	if err != nil {
		return 0, err
	}
	return s.conn.Read(b)
}

func (s *Socket) Close() error {
	return s.conn.Close()
}
