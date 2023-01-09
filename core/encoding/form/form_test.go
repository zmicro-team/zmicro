package form

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type LoginRequest struct {
	Username string `json:"username,omitempty" form:"uname"`
	Password string `json:"password,omitempty" form:"passwd"`
}

type TestModel struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Empty string `json:"empty,omitempty"`
}

type TestNotProto struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Empty string `json:"empty,omitempty"`
}

func TestNew(t *testing.T) {
	codec := New("form")
	req := &LoginRequest{
		Username: "username",
		Password: "password",
	}
	content, err := codec.Marshal(req)
	require.NoError(t, err)
	require.Equal(t, []byte("passwd=password&uname=username"), content)
}

func TestFormCodec(t *testing.T) {
	codec := New("json")

	t.Run("Content Type", func(t *testing.T) {
		require.Equal(t, "application/x-www-form-urlencoded; charset=utf-8", codec.ContentType(struct{}{}))
	})

	t.Run("Marshal", func(t *testing.T) {
		req := &LoginRequest{
			Username: "username",
			Password: "password",
		}
		content, err := codec.Marshal(req)
		require.NoError(t, err)
		require.Equal(t, []byte("password=password&username=username"), content)

		req = &LoginRequest{
			Username: "username",
			Password: "",
		}
		content, err = codec.Marshal(req)
		require.NoError(t, err)
		require.Equal(t, []byte("username=username"), content)

		m := &TestModel{
			ID:   1,
			Name: "username",
		}
		content, err = codec.Marshal(m)
		require.NoError(t, err)
		require.Equal(t, []byte("id=1&name=username"), content)
	})
	t.Run("Unmarshal", func(t *testing.T) {
		want := &LoginRequest{
			Username: "username",
			Password: "password",
		}
		got := new(LoginRequest)
		err := codec.Unmarshal([]byte(`password=password&username=username`), got)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})
	t.Run("Marshal/Unmarshal", func(t *testing.T) {
		want := &LoginRequest{
			Username: "username",
			Password: "password",
		}
		content, err := codec.Marshal(want)
		require.NoError(t, err)

		got := new(LoginRequest)
		err = codec.Unmarshal(content, got)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})
}
func TestEncode(t *testing.T) {
	codec := New("json")
	// TODO: encode proto
	tests := []struct {
		name string
		args any
		want url.Values
	}{
		{
			"full",
			TestNotProto{
				Name:  "test",
				URL:   "https://go.dev",
				Empty: "empty",
			},
			url.Values{
				"name":  []string{"test"},
				"url":   []string{"https://go.dev"},
				"empty": []string{"empty"},
			},
		},
		{
			"omitempty empty values",
			TestNotProto{
				Name: "test",
				URL:  "https://go.dev",
			},
			url.Values{
				"name": []string{"test"},
				"url":  []string{"https://go.dev"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := codec.Encode(tt.args)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "Encode(%v)", tt.args)
		})
	}
}
func TestDecode(t *testing.T) {
	type TestDecode struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	codec := New("json")
	p1 := TestDecode{}
	type args struct {
		vars   url.Values
		target any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    any
	}{
		{
			name: "test",
			args: args{
				vars:   map[string][]string{"name": {"golang"}, "url": {"https://go.dev"}},
				target: &p1,
			},
			wantErr: false,
			want:    &TestDecode{"golang", "https://go.dev"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := codec.Decode(tt.args.vars, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("Decode() target = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}
