package protocol

type AddressResolver func(area string, addr string) (Addr, error)

type Addr interface {
	String() string
	Diff(from Addr) (int, bool)
}
