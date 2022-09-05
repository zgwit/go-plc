package protocol

type Parser func(code string, addr string) (Addr, error)

type Area struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	//regex ??
}

type Manifest struct {
	Name     string          `json:"name"`
	Label    string          `json:"label"`
	Version  string          `json:"version"`
	Areas    []Area          `json:"areas"`
	Station  bool            `json:"station"`
	Resolver AddressResolver `json:"-"`
	Factory  CreateFactory   `json:"-"`
}
