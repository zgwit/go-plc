package siemens

import (
	"github.com/zgwit/go-plc/protocol"
)

// PPI 协议
type PPI struct {
	link protocol.Messenger
}

// Read 读到数据
func (t *PPI) Read(address protocol.Addr, length int) ([]byte, error) {
	return nil, nil
}

// Write 写入数据
func (t *PPI) Write(address protocol.Addr, values []byte) error {
	return nil
}
