package encoding

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/zmicro-team/zmicro/core/encoding/codec"
	"github.com/zmicro-team/zmicro/core/encoding/internal/examplepb"
	"github.com/zmicro-team/zmicro/core/encoding/json"
	"github.com/zmicro-team/zmicro/core/encoding/msgpack"
)

var marshalers = []dummyMarshaler{0, 1}

func TestEncoding_Register(t *testing.T) {
	t.Run("empty MIME type", func(t *testing.T) {
		registry := New()

		err := registry.Register("", &json.Codec{})
		require.Error(t, err)
	})
	t.Run("<nil> marshaller not allow", func(t *testing.T) {
		registry := New()

		err := registry.Register(MIMEURI, nil)
		require.Error(t, err)
	})
	t.Run("remove MIME type", func(t *testing.T) {
		registry := New()

		got := registry.Get(MIMEMSGPACK2)
		_, ok := got.(*msgpack.Codec)
		require.True(t, ok, "should be got MIME wildcard marshaler")

		err := registry.Delete(MIMEMSGPACK2)
		require.NoError(t, err)

		got = registry.Get(MIMEMSGPACK2)
		_, ok = got.(*HTTPBodyCodec)
		require.True(t, ok, "should be got MIME wildcard marshaler")
	})
	t.Run("remove not allow MIME type", func(t *testing.T) {
		registry := New()

		err := registry.Delete(MIMEURI)
		require.Error(t, err)
		err = registry.Delete(MIMEQuery)
		require.Error(t, err)
		err = registry.Delete(MIMEURI)
		require.Error(t, err)
	})
}

func TestEncoding_MarshalerForRequest_Wildcard(t *testing.T) {
	var registry = New()

	r, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf(`http.NewRequest("GET", "http://example.com", nil) failed with %v; want success`, err)
	}

	r.Header.Set("Accept", "application/unknown")
	r.Header.Set("Content-Type", "application/unknown")
	_, in := registry.InboundForRequest(r)
	if _, ok := in.(*HTTPBodyCodec); !ok {
		t.Errorf("in = %#v; want a HTTPBodyCodec", in)
	}
	out := registry.OutboundForRequest(r)
	if _, ok := out.(*HTTPBodyCodec); !ok {
		t.Errorf("out = %#v; want a HTTPBodyCodec", in)
	}
}

func TestEncoding_MarshalerForRequest_NotWildcard(t *testing.T) {
	var registry = New()

	err := registry.Register("application/x-0", &marshalers[0])
	require.NoError(t, err)
	err = registry.Register("application/x-1", &marshalers[1])
	require.NoError(t, err)

	tests := []struct {
		name        string
		contentType string
		accept      string
		wantIn      codec.Marshaler
		wantOut     codec.Marshaler
	}{
		// You can specify a marshaler for a specific MIME type.
		// The output marshaler follows the input one unless specified.
		{
			name:        "",
			contentType: "application/x-0",
			accept:      "application/x-0",
			wantIn:      &marshalers[0],
			wantOut:     &marshalers[0],
		},
		// You can also separately specify an output marshaler
		{
			name:        "",
			contentType: "application/x-0",
			accept:      "application/x-1",
			wantIn:      &marshalers[0],
			wantOut:     &marshalers[1],
		},
		{
			name:        "",
			contentType: "application/x-1; charset=UTF-8",
			accept:      "application/x-1",
			wantIn:      &marshalers[1],
			wantOut:     &marshalers[1],
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := http.NewRequest("GET", "http://example.com", nil)
			if err != nil {
				t.Fatalf(`http.NewRequest("GET", "http://example.com", nil) failed with %v; want success`, err)
			}
			r.Header.Set("Accept", test.accept)
			r.Header.Set("Content-Type", test.contentType)
			_, in := registry.InboundForRequest(r)
			if got, want := in, test.wantIn; got != want {
				t.Errorf("in = %#v; want %#v", got, want)
			}
			out := registry.OutboundForRequest(r)
			if got, want := out, test.wantOut; got != want {
				t.Errorf("out = %#v; want %#v", got, want)
			}
		})
	}
}

type dummyMarshaler int

func (dummyMarshaler) ContentType(_ interface{}) string { return "" }
func (dummyMarshaler) Marshal(interface{}) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (dummyMarshaler) Unmarshal([]byte, interface{}) error {
	return errors.New("not implemented")
}

func (dummyMarshaler) NewDecoder(r io.Reader) codec.Decoder {
	return dummyDecoder{}
}
func (dummyMarshaler) NewEncoder(w io.Writer) codec.Encoder {
	return dummyEncoder{}
}

func (m dummyMarshaler) GoString() string {
	return fmt.Sprintf("dummyMarshaler(%d)", m)
}

type dummyDecoder struct{}

func (dummyDecoder) Decode(interface{}) error {
	return errors.New("not implemented")
}

type dummyEncoder struct{}

func (dummyEncoder) Encode(interface{}) error {
	return errors.New("not implemented")
}

type TestMode struct {
	Id   string `json:"id" yaml:"id" xml:"id" toml:"id" msgpack:"id"`
	Name string `json:"name" yaml:"name" xml:"name" toml:"name" msgpack:"name"`
}

var protoMessage = &examplepb.ABitOfEverything{
	SingleNested:        &examplepb.ABitOfEverything_Nested{},
	RepeatedStringValue: nil,
	MappedStringValue:   nil,
	MappedNestedValue:   nil,
	RepeatedEnumValue:   nil,
	TimestampValue:      &timestamppb.Timestamp{},
	Uuid:                "6EC2446F-7E89-4127-B3E6-5C05E6BECBA7",
	Nested: []*examplepb.ABitOfEverything_Nested{
		{
			Name:   "foo",
			Amount: 12345,
		},
	},
	Uint64Value: 0xFFFFFFFFFFFFFFFF,
	EnumValue:   examplepb.NumericEnum_ONE,
	OneofValue: &examplepb.ABitOfEverything_OneofString{
		OneofString: "bar",
	},
	MapValue: map[string]examplepb.NumericEnum{
		"a": examplepb.NumericEnum_ONE,
		"b": examplepb.NumericEnum_ZERO,
	},
}

func TestEncoding_Bind(t *testing.T) {
	registry := New()
	tests := []struct {
		name    string
		genReq  func() (*http.Request, error)
		want    any
		wantErr bool
	}{
		{
			"default: marshaler",
			func() (*http.Request, error) {
				marshaler := registry.Get(MIMEWildcard)

				b, err := marshaler.Marshal(&examplepb.Complex{
					Id:     11,
					Uint32: wrapperspb.UInt32(1234),
					Bool:   wrapperspb.Bool(true),
				})
				if err != nil {
					return nil, err
				}
				r, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(b))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/unknown")
				return r, nil
			},
			&examplepb.Complex{
				Id:     11,
				Uint32: wrapperspb.UInt32(1234),
				Bool:   wrapperspb.Bool(true),
			},
			false,
		},
		{
			"form - application/x-www-form-urlencoded",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader([]byte(`id=foo&name=bar`)))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"form - method get so it query",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com?id=foo&name=bar", nil)
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"form - MultipartForm",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "http://example.com", nil)
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "multipart/form-data")
				r.MultipartForm = &multipart.Form{
					Value: map[string][]string{
						"id":   {"foo"},
						"name": {"bar"},
					},
					File: nil,
				}
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"json",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader([]byte(`{"id":"foo"}`)))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/json")
				return r, nil
			},
			&examplepb.SimpleMessage{
				Id: "foo",
			},
			false,
		},
		{
			"proto",
			func() (*http.Request, error) {
				buf := &bytes.Buffer{}

				m := registry.Get("application/x-protobuf")
				err := m.NewEncoder(buf).Encode(protoMessage)
				if err != nil {
					return nil, err
				}
				r, err := http.NewRequest(http.MethodPost, "http://example.com", buf)
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-protobuf")
				return r, nil
			},
			protoMessage,
			false,
		},
		{
			"yaml",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader([]byte("id: foo\nname: bar")))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-yaml")
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"xml",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader([]byte("<TestMode><id>foo</id><name>bar</name></TestMode>")))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/xml")
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"toml",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader([]byte("id=\"foo\"\nname=\"bar\"")))
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/toml")
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"msgpack",
			func() (*http.Request, error) {
				buf := &bytes.Buffer{}

				m := registry.Get("application/x-msgpack")
				err := m.NewEncoder(buf).Encode(&TestMode{
					Id:   "foo",
					Name: "bar",
				})
				if err != nil {
					return nil, err
				}

				r, err := http.NewRequest(http.MethodPost, "http://example.com", buf)
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-msgpack")
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.genReq()
			if err != nil {
				t.Errorf("genReq() error = %v", err)
			}
			got := alloc(reflect.TypeOf(tt.want))
			if err = registry.Bind(req, got.Interface()); (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, ok := tt.want.(proto.Message); ok {
				if diff := proto.Equal(got.Interface().(proto.Message), tt.want.(proto.Message)); !diff {
					t.Errorf("got = %v, want %v", got, tt.want)
				}
			} else {
				require.Equal(t, got.Interface(), tt.want)
			}
		})
	}
}

func TestEncoding_BindQuery(t *testing.T) {
	registry := New()

	tests := []struct {
		name    string
		genReq  func() (*http.Request, error)
		want    any
		wantErr bool
	}{
		{
			"form - no proto",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com?id=foo&name=bar", nil)
				if err != nil {
					return nil, err
				}
				return r, nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"form - proto",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com?id=11&uint32=1234&bool=true", nil)
				if err != nil {
					return nil, err
				}
				return r, nil
			},
			&examplepb.Complex{
				Id:     11,
				Uint32: wrapperspb.UInt32(1234),
				Bool:   wrapperspb.Bool(true),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.genReq()
			if err != nil {
				t.Errorf("genReq() error = %v", err)
			}
			got := alloc(reflect.TypeOf(tt.want))
			if err = registry.BindQuery(req, got.Interface()); (err != nil) != tt.wantErr {
				t.Errorf("BindQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, ok := tt.want.(proto.Message); ok {
				if diff := proto.Equal(got.Interface().(proto.Message), tt.want.(proto.Message)); !diff {
					t.Errorf("got = %v, want %v", got, tt.want)
				}
			} else {
				require.Equal(t, got.Interface(), tt.want)
			}
		})
	}
}

func TestEncoding_BindUri(t *testing.T) {
	registry := New()

	tests := []struct {
		name    string
		genReq  func() (*http.Request, error)
		want    any
		wantErr bool
	}{
		{
			"uri - no proto",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com/foo/bar", nil)
				if err != nil {
					return nil, err
				}

				param := url.Values{}
				param.Add("id", "foo")
				param.Add("name", "bar")
				return RequestWithUri(r, param), nil
			},
			&TestMode{
				Id:   "foo",
				Name: "bar",
			},
			false,
		},
		{
			"uri - proto",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com?id=11&uint32=1234&bool=true", nil)
				if err != nil {
					return nil, err
				}
				param := url.Values{}
				param.Add("id", "11")
				param.Add("uint32", "1234")
				param.Add("bool", "true")
				return RequestWithUri(r, param), nil
			},
			&examplepb.Complex{
				Id:     11,
				Uint32: wrapperspb.UInt32(1234),
				Bool:   wrapperspb.Bool(true),
			},
			false,
		},
		{
			"uri - always in context",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com?id=11&uint32=1234&bool=true", nil)
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return RequestWithUri(r, nil), nil
			},
			&examplepb.Complex{},
			false,
		},
		{
			"uri - no existing in context",
			func() (*http.Request, error) {
				r, err := http.NewRequest(http.MethodGet, "http://example.com?id=11&uint32=1234&bool=true", nil)
				if err != nil {
					return nil, err
				}
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return r, nil
			},
			&examplepb.Complex{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.genReq()
			if err != nil {
				t.Errorf("genReq() error = %v", err)
			}
			got := alloc(reflect.TypeOf(tt.want))
			if err = registry.BindUri(req, got.Interface()); (err != nil) != tt.wantErr {
				t.Errorf("BindQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, ok := tt.want.(proto.Message); ok {
				if diff := proto.Equal(got.Interface().(proto.Message), tt.want.(proto.Message)); !diff {
					t.Errorf("got = %v, want %v", got, tt.want)
				}
			} else {
				require.Equal(t, got.Interface(), tt.want)
			}
		})
	}
}

// helper
func alloc(t reflect.Type) reflect.Value {
	if t == nil {
		return reflect.ValueOf(new(interface{}))
	}
	return reflect.New(t.Elem())
}

func TestEncoding_Render(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		genReq func() (*http.Request, error)
		v      any
	}
	tests := []struct {
		name     string
		encoding *Encoding
		args     args
		want     string
		wantErr  bool
	}{
		{
			"<nil> payload",
			New(),
			args{
				w: httptest.NewRecorder(),
				genReq: func() (*http.Request, error) {
					return http.NewRequest(http.MethodGet, "http://example.com", nil)
				},
				v: nil,
			},
			"",
			false,
		},
		{
			"<nil> payload",
			New(),
			args{
				w: httptest.NewRecorder(),
				genReq: func() (*http.Request, error) {
					req, err := http.NewRequest(http.MethodPost, "http://example.com", nil)
					if err != nil {
						return nil, err
					}
					req.Header.Set("Accept", "application/json; charset=utf-8")
					return req, nil
				},
				v: TestMode{
					Id:   "foo",
					Name: "bar",
				},
			},
			`{"id":"foo","name":"bar"}`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.args.genReq()
			if err != nil {
				t.Errorf("genReq() error = %v", err)
			}
			if err = tt.encoding.Render(tt.args.w, req, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
			}
			w := tt.args.w.(*httptest.ResponseRecorder)
			if got := w.Body.String(); got != tt.want {
				t.Errorf("Render() result got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAcceptHeader(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   []string
	}{
		{
			"",
			"application/json, text/plain, */*",
			[]string{"application/json", "text/plain", "*/*"},
		},
		{
			"",
			"application/json,text/plain,   */*",
			[]string{"application/json", "text/plain", "*/*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseAcceptHeader(tt.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAcceptHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}
