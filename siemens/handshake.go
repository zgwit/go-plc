package siemens

import (
	"bytes"
	"github.com/zgwit/go-plc/helper"
)

const handshake1_200_smart = "03 00 00 16 11 E0 00 00 00 01 00 C1 02 10 00 C2 02 03 00 C0 01 0A"
const handshake1_200 = "03 00 00 16 11 E0 00 00 00 01 00 C1 02 4D 57 C2 02 4D 57 C0 01 09"
const handshake1_300 = "03 00 00 16 11 E0 00 00 00 01 00 C0 01 0A C1 02 01 02 C2 02 01 02"
const handshake1_400 = "03 00 00 16 11 E0 00 00 00 01 00 C0 01 0A C1 02 01 00 C2 02 01 03"
const handshake1_1200 = "03 00 00 16 11 E0 00 00 00 01 00 C0 01 0A C1 02 01 02 C2 02 01 00"
const handshake1_1500 = "03 00 00 16 11 E0 00 00 00 01 00 C0 01 0A C1 02 01 02 C2 02 01 00"

const handshake2_200_smart = "03 00 00 19 02 F0 80 32 01 00 00 CC C1 00 08 00 00 F0 00 00 01 00 01 03 C0"
const handshake2_200 = "03 00 00 19 02 F0 80 32 01 00 00 00 00 00 08 00 00 F0 00 00 01 00 01 03 C0"
const handshake2_300 = "03 00 00 19 02 F0 80 32 01 00 00 04 00 00 08 00 00 F0 00 00 01 00 01 01 E0"
const handshake2_400 = "03 00 00 19 02 F0 80 32 01 00 00 04 00 00 08 00 00 F0 00 00 01 00 01 01 E0"
const handshake2_1200 = "03 00 00 19 02 F0 80 32 01 00 00 04 00 00 08 00 00 F0 00 00 01 00 01 01 E0"
const handshake2_1500 = "03 00 00 19 02 F0 80 32 01 00 00 04 00 00 08 00 00 F0 00 00 01 00 01 01 E0"

func parseHex(str string) []byte {
	//删除空格
	buf := bytes.ReplaceAll([]byte(str), []byte{' '}, []byte{})
	return helper.FromHex(buf)
	//str = strings.ReplaceAll(str, " ", "")
	//buf := helper.FromHex([]byte(str))
	//return buf
}

func setRackSlot(handshake1 []byte, rack, slot uint8) {
	handshake1[21] = rack*0x20 + slot
}
