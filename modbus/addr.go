package modbus

func parseCode(area string) uint8 {
	var code uint8 = 1
	switch area {
	case "C":
		code = 1
	case "D":
		code = 2
	case "H":
		code = 3
	case "I":
		code = 4
	}
	return code
}
