package json

import (
	"encoding/json"
	"io"

	"github.com/charliego3/pallas/encoding"
)

// Type json Codec name
const Type = "json"

// codec is a Codec implemention with json
type codec struct{}

// Marshal data v to bytes, if err not nil
func (codec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal bytes to any pointer
func (codec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// Type json Codec type name
func (codec) Type() string {
	return Type
}

func (codec) Encoder(w io.Writer) encoding.Encoder {
	return json.NewEncoder(w)
}

func (codec) Decoder(r io.Reader) encoding.Decoder {
	return json.NewDecoder(r)
}

func init() {
	encoding.RegisterCodec(new(codec))
}
