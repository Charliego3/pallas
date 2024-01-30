package encoding

import (
	"strings"

	"github.com/charliego3/pallas/utility"
)

type Codec interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
	Type() string
}

var registeredCodec = make(map[string]Codec)

func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("codec: can not register a nil Codec")
	}

	ctype := codec.Type()
	if utility.IsBlank(ctype) {
		panic("codec: can not register Codec with empty type")
	}

	registeredCodec[strings.ToLower(ctype)] = codec
}

func CodecWithType(ctype string) (Codec, bool) {
	codec, ok := registeredCodec[strings.ToLower(ctype)]
	return codec, ok
}
