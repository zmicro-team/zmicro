package encoding

import (
	"google.golang.org/genproto/googleapis/api/httpbody"

	"github.com/zmicro-team/zmicro/core/encoding/codec"
)

// HTTPBodyCodec is a Marshaler which supports marshaling of a
// google.api.HttpBody message as the full response body if it is
// the actual message used as the response. If not, then this will
// simply fallback to the Marshaler specified as its default Marshaler.
type HTTPBodyCodec struct {
	codec.Marshaler
}

// ContentType returns its specified content type in case v is a
// google.api.HttpBody message, otherwise it will fall back to the default Marshalers
// content type.
func (h *HTTPBodyCodec) ContentType(v interface{}) string {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetContentType()
	}
	return h.Marshaler.ContentType(v)
}

// Marshal marshals "v" by returning the body bytes if v is a
// google.api.HttpBody message, otherwise it falls back to the default Marshaler.
func (h *HTTPBodyCodec) Marshal(v interface{}) ([]byte, error) {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.Data, nil
	}
	return h.Marshaler.Marshal(v)
}
