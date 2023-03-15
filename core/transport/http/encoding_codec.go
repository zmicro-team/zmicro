package http

import (
	"encoding/json"
	"io"

	"github.com/zmicro-team/zmicro/core/encoding/codec"
	"github.com/zmicro-team/zmicro/core/encoding/jsonpb"
)

// Codec is a Marshaler which marshals/unmarshals into/from JSON/
// marshals use encoding/json
// unmarshals use google.golang.org/protobuf/encoding/protojson
type Codec struct {
	*jsonpb.Codec
}

// ContentType always Returns "application/json; charset=utf-8".
func (*Codec) ContentType(_ interface{}) string {
	return "application/json; charset=utf-8"
}
func (*Codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (c *Codec) Unmarshal(data []byte, v interface{}) error {
	return c.Codec.Unmarshal(data, v)
}
func (c *Codec) NewEncoder(w io.Writer) codec.Encoder {
	return json.NewEncoder(w)
}
func (c *Codec) NewDecoder(r io.Reader) codec.Decoder {
	return c.Codec.NewDecoder(r)
}
func (c *Codec) Delimiter() []byte {
	return []byte("\n")
}
