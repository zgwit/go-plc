package helper

func Sum(buf []byte) byte {
	var sum byte = 0
	l := len(buf)
	for i := 0; i < l; i++ {
		sum += buf[i]
	}
	return sum
}

func Xor(buf []byte) byte {
	var xor byte = buf[0]
	l := len(buf)
	for i := 1; i < l; i++ {
		xor ^= buf[i]
	}
	return xor
}