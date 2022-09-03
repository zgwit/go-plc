package mitsubishi

import (
	"github.com/zgwit/go-plc/protocol"
	"io"
)

var ManifestFxProgram = protocol.Manifest{
	Name:     "Fx-Program",
	Version:  "1.0",
	Label:    "Fx-Program",
	Resolver: ParseFxProgramAddress,
	Factory: func(link io.ReadWriter, opts string) protocol.Protocol {
		return &FxProgram{
			link: protocol.Messenger{Conn: link},
		}
	},
}
