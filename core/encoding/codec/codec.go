package codec

import (
	"io"
	"net/url"
)

// Marshaler defines a conversion between byte sequence and gRPC payloads / fields.
type Marshaler interface {
	// ContentType returns the Content-Type which this marshaler is responsible for.
	// The parameter describes the type which is being marshalled, which can sometimes
	// affect the content type returned.
	ContentType(v interface{}) string
	// Marshal marshals "v" into byte sequence.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal unmarshals "data" into "v".
	// "v" must be a pointer value.
	Unmarshal(data []byte, v interface{}) error
	// NewDecoder returns a Decoder which reads byte sequence from "r".
	NewDecoder(r io.Reader) Decoder
	// NewEncoder returns an Encoder which writes bytes sequence into "w".
	NewEncoder(w io.Writer) Encoder
}

// FormCodec encode or decode a url.values
type FormCodec interface {
	Encode(v any) (url.Values, error)
	Decode(vs url.Values, v any) error
}

// UriEncoder encode to url path
type UriEncoder interface {
	// EncodeURL encode v to url path.
	// pathTemplate is a template of url path like http://helloworld.dev/{name}/sub/{sub.name},
	EncodeURL(pathTemplate string, v any, needQuery bool) string
}

// FormMarshaler defines a conversion between byte sequence and gRPC payloads / fields.
type FormMarshaler interface {
	Marshaler
	FormCodec
}

// UriMarshaler defines a conversion between byte sequence and gRPC payloads / fields.
type UriMarshaler interface {
	Marshaler
	FormCodec
	UriEncoder
}

// Decoder decodes a byte sequence
type Decoder interface {
	Decode(v interface{}) error
}

// Encoder encodes gRPC payloads / fields into byte sequence.
type Encoder interface {
	Encode(v interface{}) error
}

// DecoderFunc adapts an decoder function into Decoder.
type DecoderFunc func(v interface{}) error

// Decode delegates invocations to the underlying function itself.
func (f DecoderFunc) Decode(v interface{}) error { return f(v) }

// EncoderFunc adapts an encoder function into Encoder
type EncoderFunc func(v interface{}) error

// Encode delegates invocations to the underlying function itself.
func (f EncoderFunc) Encode(v interface{}) error { return f(v) }
