package encoding

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	t.Run("in context", func(t *testing.T) {
		want := url.Values{
			"id":   {"foo"},
			"name": {"bar"},
		}

		req, err := http.NewRequest(http.MethodGet, "http://example.com/foo/bar", nil)
		require.NoError(t, err)
		req = RequestWithUri(req, want)
		got := FromRequestUri(req)
		require.NotNil(t, got)
		require.Equal(t, got, want)
	})

	t.Run("empty in context", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://example.com/foo/bar", nil)
		require.NoError(t, err)
		got := FromRequestUri(req)
		require.Nil(t, got)
	})
}
