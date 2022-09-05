package siemens

import (
	"fmt"
	"github.com/zgwit/go-plc/protocol"
)

type S7 struct {
	handshake1 []byte
	handshake2 []byte

	link protocol.Messenger
	buf  []byte
}

func (s *S7) HandShake() error {
	_, err := s.link.Ask(s.handshake1, s.buf)
	if err != nil {
		return err
	}
	//TODO 检查结果
	_, err = s.link.Ask(s.handshake2, s.buf)
	if err != nil {
		return err
	}
	//TODO 检查结果
	return nil
}

func (s *S7) Read(station int, address protocol.Addr, size int) ([]byte, error) {
	addr := address.(*Address)

	var vt uint8 = VariableTypeWord
	offset := addr.Offset
	if addr.HasBit {
		vt = VariableTypeBit
		offset = offset*8 + uint32(addr.Bits)
	}

	pack := S7Package{
		Type:      MessageTypeJobRequest,
		Reference: 0,
		param: S7Parameter{
			Code:  ParameterTypeRead,
			Count: 1,
			Type:  vt,
			Areas: []S7ParameterArea{
				{
					Code:   addr.Code,
					DB:     addr.DB,
					Size:   uint16(size),
					Offset: offset,
				},
			},
		},
	}

	cmd := pack.encode()

	n, err := s.link.Ask(cmd, s.buf)
	if err != nil {
		return nil, err
	}

	//解析数据
	var resp S7Package
	err = resp.decode(s.buf[:n])
	if err != nil {
		return nil, err
	}

	return resp.data[0].Data, nil
}

func (s *S7) Write(station int, address protocol.Addr, data []byte) error {
	addr := address.(*Address)
	length := len(data)

	var vt uint8 = VariableTypeWord
	offset := addr.Offset
	if addr.HasBit {
		vt = VariableTypeBit
		offset = offset*8 + uint32(addr.Bits)
	}

	pack := S7Package{
		Type:      MessageTypeJobRequest,
		Reference: 1,
		param: S7Parameter{
			Code:  ParameterTypeWrite,
			Count: 1,
			Type:  vt,
			Areas: []S7ParameterArea{
				{
					Code:   addr.Code,
					DB:     addr.DB,
					Size:   uint16(length),
					Offset: offset,
				},
			},
		},
		data: []S7Data{{
			Type:  vt + 2, //transport size 3 bit 4 byte/word/qword
			Count: uint16(length),
			Data:  data,
		}},
	}

	cmd := pack.encode()

	n, err := s.link.Ask(cmd, s.buf)
	if err != nil {
		return err
	}

	//解析结果
	var resp S7Package
	err = resp.decode(s.buf[:n])
	if err != nil {
		return err
	}
	code := resp.data[0].Code
	if code != 0xff {
		return fmt.Errorf("错误码 %d", code)
	}

	/*
		0x00	Reserved	未定义，预留
		0x01	Hardware error	硬件错误
		0x03	Accessing the object not allowed	对象不允许访问
		0x05	Invalid addr	无效地址，所需的地址超出此PLC的极限
		0x06	Data type not supported	数据类型不支持
		0x07	Data type inconsistent	日期类型不一致
		0x0a	Object does not exist	对象不存在
		0xff	Success	成功
	*/

	return nil
}
