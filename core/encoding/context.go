package encoding

import (
	"context"
	"net/http"
	"net/url"
)

type ctxUriKey struct{}

// RequestWithUri sets the URL variables for the given request,
// Arguments are not modified, a shallow copy is returned.
// URL variables can be set by making a route that captures
// the required variables, starting a server and sending the request
// to that server.
func RequestWithUri(req *http.Request, uri url.Values) *http.Request {
	if uri == nil {
		uri = url.Values{}
	}
	ctx := context.WithValue(req.Context(), ctxUriKey{}, uri)
	return req.WithContext(ctx)
}

// FromRequestUri returns the route variables for the current request, if any.
func FromRequestUri(req *http.Request) url.Values {
	if rv := req.Context().Value(ctxUriKey{}); rv != nil {
		return rv.(url.Values)
	}
	return nil
}
