package helper

func ParseUint32(buf []byte) uint32 {
	return uint32(buf[0])<<24 +
		uint32(buf[1])<<16 +
		uint32(buf[2])<<8 +
		uint32(buf[3])
}

func ParseUint32LittleEndian(buf []byte) uint32 {
	return uint32(buf[3])<<24 +
		uint32(buf[2])<<16 +
		uint32(buf[1])<<8 +
		uint32(buf[0])
}

func ParseUint16(buf []byte) uint16 {
	return uint16(buf[0])<<8 + uint16(buf[1])
}

func ParseUint16LittleEndian(buf []byte) uint16 {
	return uint16(buf[1])<<8 + uint16(buf[0])
}

func Uint32ToBytes(value uint32) []byte {
	buf := make([]byte, 4)
	buf[0] = byte(value >> 24)
	buf[1] = byte(value >> 16)
	buf[2] = byte(value >> 8)
	buf[3] = byte(value)
	return buf
}

func Uint32ToBytesLittleEndian(value uint32) []byte {
	buf := make([]byte, 4)
	buf[3] = byte(value >> 24)
	buf[2] = byte(value >> 16)
	buf[1] = byte(value >> 8)
	buf[0] = byte(value)
	return buf
}

func Uint16ToBytes(value uint16) []byte {
	buf := make([]byte, 2)
	buf[0] = byte(value >> 8)
	buf[1] = byte(value)
	return buf
}

func Uint16ToBytesLittleEndian(value uint16) []byte {
	buf := make([]byte, 2)
	buf[1] = byte(value >> 8)
	buf[0] = byte(value)
	return buf
}

func WriteUint32(buf []byte, value uint32) {
	buf[0] = byte(value >> 24)
	buf[1] = byte(value >> 16)
	buf[2] = byte(value >> 8)
	buf[3] = byte(value)
}

func WriteUint32LittleEndian(buf []byte, value uint32) {
	buf[3] = byte(value >> 24)
	buf[2] = byte(value >> 16)
	buf[1] = byte(value >> 8)
	buf[0] = byte(value)
}

func WriteUint24(buf []byte, value uint32) {
	buf[0] = byte(value >> 16)
	buf[1] = byte(value >> 8)
	buf[2] = byte(value)
}

func WriteUint24LittleEndian(buf []byte, value uint32) {
	buf[2] = byte(value >> 16)
	buf[1] = byte(value >> 8)
	buf[0] = byte(value)
}

func WriteUint16(buf []byte, value uint16) {
	buf[0] = byte(value >> 8)
	buf[1] = byte(value)
}

func WriteUint16LittleEndian(buf []byte, value uint16) {
	buf[1] = byte(value >> 8)
	buf[0] = byte(value)
}

func BoolToAscii(buf []byte) []byte {
	length := len(buf)
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		if buf[i] == 0 {
			ret[i] = '0'
		} else {
			ret[i] = '1'
		}
	}
	return ret
}

func AsciiToBool(buf []byte) []byte {
	length := len(buf)
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		if buf[i] == '0' {
			ret[i] = 0
		} else {
			ret[i] = 1
		}
	}
	return ret
}

func Dup(buf []byte) []byte {
	b := make([]byte, len(buf))
	copy(b, buf)
	return b
}

func BoolToByte(buf []bool) []byte {
	r := make([]byte, len(buf))
	for i, v := range buf {
		if v {
			r[i] = 1
		}
	}
	return r
}

func ByteToBool(buf []byte) []bool  {
	r := make([]bool, len(buf))
	for i, v := range buf {
		if v > 0 {
			r[i] = true
		}
	}
	return r
}

//ShrinkBool
func BooleansToBytes(buf []bool) []byte {
	length := len(buf)
	//length = length % 8 == 0 ? length / 8 : length / 8 + 1;
	ln := length >> 3    // length/8
	if length&0x07 > 0 { // length%8
		ln++
	}

	b := make([]byte, ln)

	for i := 0; i < length; i++ {
		if buf[i] {
			//b[i/8] += 1 << (i % 8)
			b[i>>3] += 1 << (i & 0x07)
		}
	}

	return b
}

//ExpandBool
func BytesToBooleans(buf []byte) []bool {
	length := len(buf)
	ln := length << 3 // length * 8
	b := make([]bool, ln)

	for i := 0; i < length; i++ {
		//b[i] = buf[i/8] & (1 << (i % 8))
		if buf[i>>3] & (1 << (i & 0x07)) > 0 {
			b[i] = true
		}
	}

	return b
}
