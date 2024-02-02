package xml

import (
	"encoding/xml"
	"io"

	"github.com/charliego3/pallas/encoding"
)

// Type xml Codec name
const Type = "xml"

// codec is a XML Codec implemention
type codec struct{}

// Marshal object to bytes with xml
func (codec) Marshal(v any) ([]byte, error) {
	return xml.Marshal(v)
}

// Unmarshal bytes to object in xml
func (codec) Unmarshal(data []byte, v any) error {
	return xml.Unmarshal(data, v)
}

// Type xml Codec type name
func (codec) Type() string {
	return Type
}

func (codec) Encoder(w io.Writer) encoding.Encoder {
	return xml.NewEncoder(w)
}

func (codec) Decoder(r io.Reader) encoding.Decoder {
	return xml.NewDecoder(r)
}

func init() {
	encoding.RegisterCodec(new(codec))
}
