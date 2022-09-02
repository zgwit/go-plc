package mitsubishi

import (
	"github.com/zgwit/go-plc/protocol"
)

var Fx_Program = protocol.Desc{
	Name:    "Fx-Program",
	Version: "1.0",
	Label:   "Fx-Program",
	Factory: func(tunnel io.ReadWriter, opts string) protocol.Protocol {
		return &FxProgram{
			link: tunnel,
		}
	},
	Parser: ParseFxProgramAddress,
}
