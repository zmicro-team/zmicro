package msgpack

import (
	"bytes"
	"io"

	msgpack "github.com/ugorji/go/codec"

	"github.com/zmicro-team/zmicro/core/encoding/codec"
)

// Codec is a Codec implementation with xml.
type Codec struct{}

// ContentType always Returns "application/x-msgpack; charset=utf-8".
func (*Codec) ContentType(_ interface{}) string {
	return "application/x-msgpack; charset=utf-8"
}
func (c *Codec) Marshal(v interface{}) ([]byte, error) {
	b := &bytes.Buffer{}
	err := c.NewEncoder(b).Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
func (c *Codec) Unmarshal(data []byte, v interface{}) error {
	return c.NewDecoder(bytes.NewReader(data)).Decode(v)
}
func (*Codec) NewDecoder(r io.Reader) codec.Decoder {
	return msgpack.NewDecoder(r, new(msgpack.MsgpackHandle))
}
func (*Codec) NewEncoder(w io.Writer) codec.Encoder {
	return msgpack.NewEncoder(w, new(msgpack.MsgpackHandle))
}
