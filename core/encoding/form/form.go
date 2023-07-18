package form

import (
	"io"
	"net/url"
	"reflect"

	"github.com/go-playground/form/v4"
	"google.golang.org/protobuf/proto"

	"github.com/zmicro-team/zmicro/core/encoding/codec"
)

type Codec struct {
	Encoder *form.Encoder
	Decoder *form.Decoder
	TagName string
	// UseProtoNames uses proto field name instead of
	// lowerCamelCase name in JSON field names.
	UseProtoNames bool
	// UseEnumNumbers emits enum values as numbers.
	UseEnumNumbers bool
}

// New returns a new Codec,
// default tag name is "json", and proto message use proto field.
func New(tagName string) *Codec {
	encoder := form.NewEncoder()
	encoder.SetTagName(tagName)
	decoder := form.NewDecoder()
	decoder.SetTagName(tagName)
	return &Codec{
		encoder,
		decoder,
		tagName,
		true,
		true,
	}
}

// DisableUseProtoNames disable proto field name, use lowerCamelCase name
// in JSON field names.
func (c *Codec) DisableUseProtoNames() *Codec {
	c.UseProtoNames = false
	return c
}

// DisableUseProtoNames disable emits enum values as numbers.
func (c *Codec) DisableUseEnumNumbers() *Codec {
	c.UseEnumNumbers = false
	return c
}

// ContentType always Returns "application/x-www-form-urlencoded; charset=utf-8"
func (*Codec) ContentType(_ interface{}) string {
	return "application/x-www-form-urlencoded; charset=utf-8"
}
func (c *Codec) Marshal(v any) ([]byte, error) {
	vs, err := c.Encode(v)
	if err != nil {
		return nil, err
	}
	return []byte(vs.Encode()), nil
}
func (c *Codec) Unmarshal(data []byte, v any) error {
	vs, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}
	return c.Decode(vs, v)
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

func (c *Codec) Encode(v any) (url.Values, error) {
	var vs url.Values
	var err error

	if m, ok := v.(proto.Message); ok {
		vs, err = EncodeValues(m, c.UseProtoNames, c.UseEnumNumbers)
	} else {
		vs, err = c.Encoder.Encode(v)
	}
	if err != nil {
		return nil, err
	}
	for k, vv := range vs {
		if len(vv) == 0 {
			delete(vs, k)
		}
	}
	return vs, nil
}

func (c *Codec) Decode(vs url.Values, v any) error {
	if m, ok := v.(proto.Message); ok {
		return DecodeValues(m, vs)
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	if m, ok := rv.Interface().(proto.Message); ok {
		return DecodeValues(m, vs)
	}
	return c.Decoder.Decode(v, vs)
}

type MultipartCodec struct {
	*Codec
}

func (*MultipartCodec) ContentType(_ interface{}) string {
	return "multipart/form-data"
}

type QueryCodec struct {
	*Codec
}

func (*QueryCodec) ContentType(_ interface{}) string {
	return "__MIME__/Query"
}

type UriCodec struct {
	*Codec
}

func (*UriCodec) ContentType(_ interface{}) string {
	return "__MIME__/URI"
}
