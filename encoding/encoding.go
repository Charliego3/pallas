package encoding

import (
	"fmt"
	"io"
	"strings"

	"github.com/charliego3/pallas/utility"
)

type Codec interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
	Type() string
}

type Coder interface {
	Encoder(w io.Writer) Encoder
	Decoder(r io.Reader) Decoder
}

type Encoder interface {
	Encode(v any) error
}

type Decoder interface {
	Decode(v any) error
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

func CodecWithType(ctype string) Codec {
	if codec, ok := registeredCodec[strings.ToLower(ctype)]; ok {
		return codec
	}
	panic(fmt.Sprintf("forget register Codec? type: [%s]", ctype))
}
