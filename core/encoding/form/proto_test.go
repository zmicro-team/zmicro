package form

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/zmicro-team/zmicro/core/encoding/internal/examplepb"
)

func TestProto_Codec(t *testing.T) {
	codec := New("json")

	want := &examplepb.Complex{
		Id:      2233,
		NoOne:   "2233",
		Simple:  &examplepb.Simple{Component: "5566"},
		Simples: []string{"3344", "5566"},
		B:       true,
		Sex:     examplepb.Sex_woman,
		Age:     18,
		A:       19,
		Count:   3,
		Price:   11.23,
		D:       22.22,
		Byte:    []byte("123"),

		Timestamp: &timestamppb.Timestamp{Seconds: 20, Nanos: 2},
		Duration:  &durationpb.Duration{Seconds: 120, Nanos: 22},
		Field:     &fieldmaskpb.FieldMask{Paths: []string{"1", "2"}},

		Double:  &wrapperspb.DoubleValue{Value: 12.33},
		Float:   &wrapperspb.FloatValue{Value: 12.34},
		Int64:   &wrapperspb.Int64Value{Value: 64},
		Int32:   &wrapperspb.Int32Value{Value: 32},
		Uint64:  &wrapperspb.UInt64Value{Value: 64},
		Uint32:  &wrapperspb.UInt32Value{Value: 32},
		Bool:    &wrapperspb.BoolValue{Value: false},
		String_: &wrapperspb.StringValue{Value: "golang"},
		Bytes:   &wrapperspb.BytesValue{Value: []byte("123")},

		Map: map[string]string{"key": "https://go.dev"},
	}
	content, err := codec.Marshal(want)
	require.NoError(t, err)
	require.Equal(t, "a=19&age=18&b=true&bool=false&byte=MTIz&bytes=MTIz&count=3&d=22.22&double=12.33&duration="+
		"2m0.000000022s&field=1%2C2&float=12.34&id=2233&int32=32&int64=64&map%5Bkey%5D=https%3A%2F%2Fgo.dev&"+
		"numberOne=2233&price=11.23&sex=woman&simples=3344&simples=5566&string=golang"+
		"&timestamp=1970-01-01T00%3A00%3A20.000000002Z&uint32=32&uint64=64&very_simple.component=5566", string(content))

	got := &examplepb.Complex{}
	err = codec.Unmarshal(content, got)
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(got, got, protocmp.Transform()))
}
