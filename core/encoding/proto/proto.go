package proto

import (
	"errors"
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/zmicro-team/zmicro/core/encoding/codec"
)

// Codec is a Marshaller which marshals/unmarshals into/from serialize proto bytes
type Codec struct{}

// ContentType always returns "application/x-protobuf".
func (*Codec) ContentType(_ interface{}) string {
	return "application/x-protobuf"
}
func (*Codec) Marshal(value interface{}) ([]byte, error) {
	message, ok := value.(proto.Message)
	if !ok {
		return nil, errors.New("unable to marshal non proto field")
	}
	return proto.Marshal(message)
}
func (*Codec) Unmarshal(data []byte, value interface{}) error {
	message, ok := value.(proto.Message)
	if !ok {
		return errors.New("unable to unmarshal non proto field")
	}
	return proto.Unmarshal(data, message)
}
func (c *Codec) NewDecoder(r io.Reader) codec.Decoder {
	return codec.DecoderFunc(func(value interface{}) error {
		buffer, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		return c.Unmarshal(buffer, value)
	})
}
func (c *Codec) NewEncoder(w io.Writer) codec.Encoder {
	return codec.EncoderFunc(func(value interface{}) error {
		buffer, err := c.Marshal(value)
		if err != nil {
			return err
		}
		_, err = w.Write(buffer)
		return err
	})
}
