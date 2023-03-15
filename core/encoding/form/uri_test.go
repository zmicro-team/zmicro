package form

import (
	"testing"

	"github.com/zmicro-team/zmicro/core/encoding/internal/examplepb"
)

type NoProtoSub struct {
	Name string `json:"name"`
}

type NoProtoHello struct {
	Name string      `json:"name"`
	Sub  *NoProtoSub `json:"sub"`
}

func TestEncodeURL(t *testing.T) {
	type args struct {
		pathTemplate string
		msg          any
		needQuery    bool
	}
	codec := New("json").EnableProto()
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"proto: no any param",
			args{
				"http://hello.dev/sub",
				&examplepb.HelloRequest{
					Name: "test",
					Sub:  &examplepb.Sub{Name: "2233!!!"},
				},
				false,
			},
			`http://hello.dev/sub`,
		},
		{
			"proto: param",
			args{
				"http://hello.dev/{name}/sub/{sub.name}",
				&examplepb.HelloRequest{
					Name: "test",
					Sub:  &examplepb.Sub{Name: "2233!!!"},
				},
				false,
			},
			`http://hello.dev/test/sub/2233!!!`,
		},
		{
			"proto: param with proto [json_name=naming]",
			args{
				"http://hello.dev/{name}/sub/{sub.naming}",
				&examplepb.HelloRequest{
					Name: "test",
					Sub:  &examplepb.Sub{Name: "5566!!!"},
				},
				false,
			},
			`http://hello.dev/test/sub/5566!!!`,
		},
		{
			"proto: param with empty",
			args{
				"http://hello.dev/{name}/sub/{sub.name}",
				&examplepb.HelloRequest{
					Name: "test",
				},
				false,
			},
			`http://hello.dev/test/sub/`,
		},
		{
			"proto: param not match",
			args{
				"http://hello.dev/{name}/sub/{sub.name33}",
				&examplepb.HelloRequest{
					Name: "test",
				},
				false,
			},
			`http://hello.dev/test/sub/{sub.name33}`,
		},
		{
			"proto: param with query",
			args{
				"http://hello.dev/{name}/sub",
				&examplepb.HelloRequest{
					Name: "go",
					Sub:  &examplepb.Sub{Name: "golang"},
				},
				true,
			},
			`http://hello.dev/go/sub?sub.naming=golang`,
		},

		{
			"no proto: no any param",
			args{
				"http://hello.dev/sub",
				&NoProtoHello{
					Name: "test",
					Sub:  &NoProtoSub{Name: "2233!!!"},
				},
				false,
			},
			`http://hello.dev/sub`,
		},
		{
			"no proto: param",
			args{
				"http://hello.dev/{name}/sub/{sub.name}",
				&NoProtoHello{
					Name: "test",
					Sub:  &NoProtoSub{Name: "2233!!!"},
				},
				false,
			},
			`http://hello.dev/test/sub/2233!!!`,
		},
		{
			"no proto: param with empty",
			args{
				"http://hello.dev/{name}/sub/{sub.name}",
				&NoProtoHello{
					Name: "test",
				},
				false,
			},
			`http://hello.dev/test/sub/`,
		},
		{
			"no proto: param not match",
			args{
				"http://hello.dev/{name}/sub/{sub.name33}",
				&NoProtoHello{
					Name: "test",
				},
				false,
			},
			`http://hello.dev/test/sub/{sub.name33}`,
		},
		{
			"no proto: param with query",
			args{
				"http://hello.dev/{name}/sub",
				&NoProtoHello{
					Name: "go",
					Sub:  &NoProtoSub{Name: "golang"},
				},
				true,
			},
			`http://hello.dev/go/sub?sub.name=golang`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := codec.EncodeURL(tt.args.pathTemplate, tt.args.msg, tt.args.needQuery); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
