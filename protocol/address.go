package protocol

import "iot-master/model"

//Addr 地址
//type Addr struct {
//	Station  uint8  `json:"slave,omitempty"`
//	Code   uint8  `json:"code"`
//	DB  uint16 `json:"block,omitempty"`
//	Offset uint16 `json:"offset"`
//}

type Addr interface {
	String() string
	Lookup(data []byte, from Addr, tp model.DataType, le bool, precision int) (interface{}, bool)
	//Flatten(value interface{}, tp DataType, le bool, precision int) []byte
}

//
//type Url struct {
//	Code   uint16 //uint8
//	Offset uint32
//	Extra  uint32
//}
