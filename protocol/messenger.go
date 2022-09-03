package protocol

import (
	"io"
	"sync"
)

type Messenger struct {
	mu   sync.Mutex
	Conn io.ReadWriter
}

func (m *Messenger) Ask(request []byte, response []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, err := m.Conn.Write(request)
	if err != nil {
		return 0, err
	}
	return m.Conn.Read(response)
}

func (m *Messenger) AskAtLeast(request []byte, response []byte, min int) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, err := m.Conn.Write(request)
	if err != nil {
		return 0, err
	}
	return io.ReadAtLeast(m.Conn, response, min)
}
