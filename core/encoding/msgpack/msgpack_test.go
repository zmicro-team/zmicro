package msgpack

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodec_ContentType(t *testing.T) {
	codec := Codec{}

	want := "application/x-msgpack; charset=utf-8"
	if got := codec.ContentType(struct{}{}); got != want {
		t.Errorf("m.ContentType(_) failed, got = %q; want %q; ", got, want)
	}
}

type testMode struct {
	Foo string `msgpack:"foo"`
}

func TestCodec_Marshal_Unmarshal(t *testing.T) {
	codec := Codec{}

	want := &testMode{Foo: "FOO"}
	got := &testMode{}

	b, err := codec.Marshal(want)
	require.NoError(t, err)

	err = codec.Unmarshal(b, got)
	require.NoError(t, err)

	require.Equal(t, want, got)
}

func TestCodec_Encoder_Decoder(t *testing.T) {
	codec := Codec{}

	want := &testMode{Foo: "FOO"}
	got := &testMode{}

	buf := &bytes.Buffer{}
	err := codec.NewEncoder(buf).Encode(want)
	require.NoError(t, err)

	err = codec.NewDecoder(buf).Decode(got)
	require.NoError(t, err)

	require.Equal(t, want, got)
}
