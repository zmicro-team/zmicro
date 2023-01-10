package encoding

import (
	"bytes"
	"testing"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/zmicro-team/zmicro/core/encoding/jsonpb"
)

func TestCodec_ContentType(t *testing.T) {
	m := HTTPBodyCodec{
		&jsonpb.Codec{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
		},
	}
	expected := "CustomContentType"
	message := &httpbody.HttpBody{
		ContentType: expected,
	}
	res := m.ContentType(nil)
	if res != "application/json; charset=utf-8" {
		t.Errorf("content type not equal (%q, %q)", res, expected)
	}
	res = m.ContentType(message)
	if res != expected {
		t.Errorf("content type not equal (%q, %q)", res, expected)
	}
}

func TestCodec_Marshal(t *testing.T) {
	m := HTTPBodyCodec{
		&jsonpb.Codec{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
		},
	}
	expected := []byte("Some test")
	message := &httpbody.HttpBody{
		Data: expected,
	}
	res, err := m.Marshal(message)
	if err != nil {
		t.Errorf("m.Marshal(%#v) failed with %v; want success", message, err)
	}
	if !bytes.Equal(res, expected) {
		t.Errorf("Marshalled data not equal (%q, %q)", res, expected)
	}
}
