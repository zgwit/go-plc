package protocol

import "iot-master/connect"

type Options map[string]interface{}

func (opts Options) GetInt(name string, value int) int {
	v, ok := opts[name]
	if ok {
		return v.(int)
	}
	return value
}

type Factory func(tunnel connect.Tunnel, opts Options) Protocol

type Parser func(code string, addr string) (Addr, error)

type Code struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	//regex ??
}

type Desc struct {
	Name    string  `json:"name"`
	Label   string  `json:"label"`
	Version string  `json:"version"`
	Codes   []Code  `json:"codes"`
	Parser  Parser  `json:"-"`
	Factory Factory `json:"-"`
}
