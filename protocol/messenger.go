package protocol

type Messenger struct {
	Conn io.ReadWriter
}

func (m *Messenger) Ask(request []byte, response []byte) (int, error) {
	_, err := m.Conn.Write(request)
	if err != nil {
		return 0, err
	}
	return m.Conn.Read(response)
}

func (m *Messenger) AskAtLeast(request []byte, response []byte, min int) (int, error) {
	_, err := m.Conn.Write(request)
	if err != nil {
		return 0, err
	}
	return io.ReadAtLeast(m.Conn, response, min)
}
