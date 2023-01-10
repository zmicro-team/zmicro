package yaml

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

func TestCodec_ContentType(t *testing.T) {
	var m Codec

	want := "application/x-yaml; charset=utf-8"
	if got := m.ContentType(struct{}{}); got != want {
		t.Errorf("m.ContentType(_) failed, got = %q; want %q; ", got, want)
	}
}

func TestCodec_Marshal(t *testing.T) {
	codec := Codec{}

	value := map[string]string{"v": "hi"}
	got, err := codec.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal should not return err(%v)", err)
	}
	if string(got) != "v: hi\n" {
		t.Fatalf("want \"v: hi\n\" return \"%s\"", string(got))
	}
}

func TestCodec_Unmarshal(t *testing.T) {
	codec := Codec{}

	for _, tt := range unmarshalerTests {
		v := reflect.ValueOf(tt.value).Type()
		value := reflect.New(v)
		err := codec.Unmarshal([]byte(tt.data), value.Interface())
		if err != nil {
			t.Fatalf("Unmarshal should not return err(%v)", err)
		}
	}
	spec := struct {
		A string
		B map[string]interface{}
	}{A: "a"}
	err := codec.Unmarshal([]byte("v: hi"), &spec.B)
	if err != nil {
		t.Fatalf("Unmarshal should not return err(%v)", err)
	}
}
func TestCodecNewEncoder(t *testing.T) {
	codec := Codec{}

	buf := &bytes.Buffer{}
	value := map[string]string{"v": "hi"}

	err := codec.NewEncoder(buf).Encode(value)
	if err != nil {
		t.Fatalf("Marshal should not return err(%v)", err)
	}

	if got := buf.String(); got != "v: hi\n" {
		t.Fatalf("want \"v: hi\n\" return \"%s\"", got)
	}
}

func TestCodec_NewDecoder(t *testing.T) {
	codec := Codec{}

	for _, tt := range unmarshalerTests {
		v := reflect.ValueOf(tt.value).Type()
		value := reflect.New(v)
		if tt.data == "" { // todo: empty string case io.EOF
			continue
		}
		err := codec.NewDecoder(bytes.NewBufferString(tt.data)).Decode(value.Interface())
		if err != nil {
			t.Fatalf("NewDecoder should not return err(%v)", err)
		}
	}
	spec := struct {
		A string
		B map[string]interface{}
	}{A: "a"}
	err := codec.Unmarshal([]byte("v: hi"), &spec.B)
	if err != nil {
		t.Fatalf("NewDecoder should not return err(%v)", err)
	}
}

var unmarshalerTests = []struct {
	data  string
	value interface{}
}{
	{
		"",
		(*struct{})(nil),
	},
	{
		"{}", &struct{}{},
	},
	{
		"v: hi",
		map[string]string{"v": "hi"},
	},
	{
		"v: hi", map[string]interface{}{"v": "hi"},
	},
	{
		"v: true",
		map[string]string{"v": "true"},
	},
	{
		"v: true",
		map[string]interface{}{"v": true},
	},
	{
		"v: 10",
		map[string]interface{}{"v": 10},
	},
	{
		"v: 0b10",
		map[string]interface{}{"v": 2},
	},
	{
		"v: 0xA",
		map[string]interface{}{"v": 10},
	},
	{
		"v: 4294967296",
		map[string]int64{"v": 4294967296},
	},
	{
		"v: 0.1",
		map[string]interface{}{"v": 0.1},
	},
	{
		"v: .1",
		map[string]interface{}{"v": 0.1},
	},
	{
		"v: .Inf",
		map[string]interface{}{"v": math.Inf(+1)},
	},
	{
		"v: -.Inf",
		map[string]interface{}{"v": math.Inf(-1)},
	},
	{
		"v: -10",
		map[string]interface{}{"v": -10},
	},
	{
		"v: -.1",
		map[string]interface{}{"v": -0.1},
	},
}
