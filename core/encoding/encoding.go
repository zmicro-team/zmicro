package encoding

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/zmicro-team/zmicro/core/encoding/codec"
	"github.com/zmicro-team/zmicro/core/encoding/form"
	"github.com/zmicro-team/zmicro/core/encoding/jsonpb"
	"github.com/zmicro-team/zmicro/core/encoding/msgpack"
	"github.com/zmicro-team/zmicro/core/encoding/proto"
	"github.com/zmicro-team/zmicro/core/encoding/toml"
	"github.com/zmicro-team/zmicro/core/encoding/xml"
	"github.com/zmicro-team/zmicro/core/encoding/yaml"
)

const defaultMemory = 32 << 20

// Content-Type MIME of the most common data formats.
const (
	// MIMEWildcard is the fallback MIME type used for requests which do not match
	// a registered MIME type.
	MIMEWildcard = "*"
	// MIMEURI is sepcial form query.
	MIMEQuery = "__MIME__/QUERY"
	// MIMEURI is sepcial form uri.
	MIMEURI = "__MIME__/URI"

	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
	MIMETOML              = "application/toml"
)

var (
	acceptHeader      = http.CanonicalHeaderKey("Accept")
	contentTypeHeader = http.CanonicalHeaderKey("Content-Type")
)

// Encoding is a mapping from MIME types to Marshalers.
type Encoding struct {
	mimeMap map[string]codec.Marshaler
}

// New encoding with default Marshalers
func New() *Encoding {
	return &Encoding{
		mimeMap: map[string]codec.Marshaler{
			MIMEWildcard: &HTTPBodyCodec{
				Marshaler: &jsonpb.Codec{
					MarshalOptions: protojson.MarshalOptions{
						UseProtoNames:  true,
						UseEnumNumbers: true,
					},
					UnmarshalOptions: protojson.UnmarshalOptions{
						DiscardUnknown: true,
					},
				},
			},
			MIMEQuery: &form.QueryCodec{Codec: form.New("json")},
			MIMEURI:   &form.UriCodec{Codec: form.New("json")},

			MIMEPOSTForm:          form.New("json"),
			MIMEMultipartPOSTForm: &form.MultipartCodec{Codec: form.New("json")},
			MIMEJSON: &jsonpb.Codec{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:  true,
					UseEnumNumbers: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
			MIMEXML:      &xml.Codec{},
			MIMEXML2:     &xml.Codec{},
			MIMEPROTOBUF: &proto.Codec{},
			MIMEMSGPACK:  &msgpack.Codec{},
			MIMEMSGPACK2: &msgpack.Codec{},
			MIMEYAML:     &yaml.Codec{},
			MIMETOML:     &toml.Codec{},
		},
	}
}

// Register a marshaler for a case-sensitive MIME type string
// ("*" to match any MIME type).
// you can override default marshaler with same MIME type
func (r *Encoding) Register(mime string, marshaler codec.Marshaler) error {
	if len(mime) == 0 {
		return errors.New("empty MIME type")
	}
	if marshaler == nil {
		return errors.New("<nil> marshaller not allow")
	} else {
		r.mimeMap[mime] = marshaler
	}
	return nil
}

// Get returns the marshalers with a case-sensitive MIME type string
// It checks the MIME type on the Encoding.
// Otherwise, it follows the above logic for "*" Marshaler.
func (r *Encoding) Get(mime string) codec.Marshaler {
	m := r.mimeMap[mime]
	if m == nil {
		m = r.mimeMap[MIMEWildcard]
	}
	return m
}

// Delete remove the MIME type marshaler.
// MIMEWildcard, MIMEQuery, MIMEURI should be always exist.
func (r *Encoding) Delete(mime string) error {
	if mime == MIMEWildcard ||
		mime == MIMEQuery ||
		mime == MIMEURI {
		return fmt.Errorf("MIME(%s) can't delete, you can override", mime)
	}
	delete(r.mimeMap, mime)
	return nil
}

// InboundForRequest returns the inbound Content-Type and marshalers for this request.
// It checks the registry on the Encoding for the MIME type set by the Content-Type header.
// If it isn't set (or the request Content-Type is empty), checks for "*".
// If there are multiple Content-Type headers set, choose the first one that it can
// exactly match in the registry.
// Otherwise, it follows the above logic for "*" Marshaler.
func (r *Encoding) InboundForRequest(req *http.Request) (string, codec.Marshaler) {
	var err error
	var marshaler codec.Marshaler
	var contentType string

	for _, contentTypeVal := range req.Header[contentTypeHeader] {
		contentType, _, err = mime.ParseMediaType(contentTypeVal)
		if err != nil {
			continue
		}
		if m, ok := r.mimeMap[contentType]; ok {
			marshaler = m
			break
		}
	}
	if marshaler == nil {
		contentType = MIMEWildcard
		marshaler = r.mimeMap[MIMEWildcard]
	}
	return contentType, marshaler
}

// OutboundForRequest returns the outbound marshalers for this request.
// It checks the registry on the Encoding for the MIME type set by the Accept header.
// If it isn't set (or the request Accept is empty), checks for "*".
// If there are multiple Accept headers set, choose the first one that it can
// exactly match in the registry.
// Otherwise, it follows the above logic for "*" Marshaler.
func (r *Encoding) OutboundForRequest(req *http.Request) codec.Marshaler {
	var marshaler codec.Marshaler
	for _, acceptVal := range req.Header[acceptHeader] {
		headerValues := parseAcceptHeader(acceptVal)
		for _, value := range headerValues {
			if m, ok := r.mimeMap[value]; ok {
				marshaler = m
				break
			}
		}
	}
	if marshaler == nil {
		marshaler = r.mimeMap[MIMEWildcard]
	}

	return marshaler
}

// Bind checks the Method and Content-Type to select codec.Marshaler automatically,
// Depending on the "Content-Type" header different bind are used, for example:
//
//	"application/json" --> JSON codec.Marshaler
//	"application/xml"  --> XML codec.Marshaler
//
// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
// It decodes the json payload into the struct specified as a pointer.
func (r *Encoding) Bind(req *http.Request, v any) error {
	if req.Method == http.MethodGet {
		return r.BindQuery(req, v)
	}
	contentType, marshaller := r.InboundForRequest(req)
	if contentType == MIMEMultipartPOSTForm {
		m, ok := marshaller.(codec.FormMarshaler)
		if !ok {
			return fmt.Errorf("not supported marshaller(%v)", contentType)
		}
		if err := req.ParseMultipartForm(defaultMemory); err != nil {
			return err
		}
		return m.Decode(req.MultipartForm.Value, v)
	}
	return marshaller.NewDecoder(req.Body).
		Decode(v)
}

// BindQuery binds the passed struct pointer using the query codec.Marshaler.
func (r *Encoding) BindQuery(req *http.Request, v any) error {
	marshaller := r.Get(MIMEQuery)
	m, ok := marshaller.(codec.FormMarshaler)
	if !ok {
		return fmt.Errorf("not supported marshaller(%v)", MIMEQuery)
	}
	return m.Decode(req.URL.Query(), v)
}

// BindUri binds the passed struct pointer using the uri codec.Marshaler.
// NOTE: before use this, you should set uri params in the request context with RequestWithUri.
func (r *Encoding) BindUri(req *http.Request, v any) error {
	raws := FromRequestUri(req)
	if raws == nil {
		return errors.New("must be request with uri in context")
	}
	marshaller := r.Get(MIMEURI)
	m, ok := marshaller.(codec.FormMarshaler)
	if !ok {
		return fmt.Errorf("not supported marshaller(%v)", MIMEURI)
	}
	return m.Decode(raws, v)
}

// Render writes the response headers and calls the outbound marshalers for this request.
// It checks the registry on the Encoding for the MIME type set by the Accept header.
// If it isn't set (or the request Accept is empty), checks for "*". for example:
//
//	"application/json" --> JSON codec.Marshaler
//	"application/xml"  --> XML codec.Marshaler
//
// If there are multiple Accept headers set, choose the first one that it can
// exactly match in the registry.
// Otherwise, it follows the above logic for "*" Marshaler.
func (r *Encoding) Render(w http.ResponseWriter, req *http.Request, v any) error {
	if v == nil {
		return nil
	}
	marshaller := r.OutboundForRequest(req)
	data, err := marshaller.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", marshaller.ContentType(v))
	_, err = w.Write(data)
	return err
}

func parseAcceptHeader(header string) []string {
	values := strings.Split(header, ",")
	for i := 0; i < len(values); i++ {
		values[i] = strings.TrimSpace(values[i])
	}
	return values
}
