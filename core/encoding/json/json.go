package json

import (
	"encoding/json"
	"io"

	"github.com/zmicro-team/zmicro/core/encoding/codec"
)

// Codec is a Marshaler which marshals/unmarshals into/from JSON
// with the standard "encoding/json" package of Golang.
// Although it is generally faster for simple proto messages than JSONPb,
// it does not support advanced features of protobuf, e.g. map, oneof, ....
//
// The NewEncoder and NewDecoder types return *json.Encoder and
// *json.Decoder respectively.
type Codec struct {
	// UseNumber causes the Decoder to unmarshal a number into an interface{} as a
	// Number instead of as a float64.
	UseNumber bool
	// DisallowUnknownFields causes the Decoder to return an error when the destination
	// is a struct and the input contains object keys which do not match any
	// non-ignored, exported fields in the destination.
	DisallowUnknownFields bool
}

// ContentType always Returns "application/json".
func (*Codec) ContentType(_ interface{}) string {
	return "application/json; charset=utf-8"
}
func (*Codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (*Codec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
func (c *Codec) NewDecoder(r io.Reader) codec.Decoder {
	decoder := json.NewDecoder(r)
	if c.UseNumber {
		decoder.UseNumber()
	}
	if c.DisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	return decoder
}
func (c *Codec) NewEncoder(w io.Writer) codec.Encoder {
	return json.NewEncoder(w)
}
func (c *Codec) Delimiter() []byte {
	return []byte("\n")
}
