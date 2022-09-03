package protocol

type AddrResolver func(area string, addr string) (Addr, error)

type Addr interface {
	String() string
	Diff(from Addr) (int, bool)
}
