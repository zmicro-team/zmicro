package form

import (
	"testing"
)

func TestTagParsing(t *testing.T) {
	name, opts := parseTag("")
	if name != "" {
		t.Fatalf("name = %q, want ''", name)
	}
	if opts.Contains("foobar") == true {
		t.Errorf("Contains(%q) = %v", "foobar", false)
	}

	name, opts = parseTag("field,foobar,foo")
	if name != "field" {
		t.Fatalf("name = %q, want field", name)
	}
	for _, tt := range []struct {
		opt  string
		want bool
	}{
		{"foobar", true},
		{"foo", true},
		{"bar", false},
	} {
		if opts.Contains(tt.opt) != tt.want {
			t.Errorf("Contains(%q) = %v", tt.opt, !tt.want)
		}
	}
}

func Test_isValidTag(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			"empty",
			"",
			false,
		},
		{
			"valid",
			"aa",
			true,
		},
		{
			"valid",
			".+aa",
			true,
		},
		{
			"invalid",
			"`",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidTag(tt.args); got != tt.want {
				t.Errorf("isValidTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
