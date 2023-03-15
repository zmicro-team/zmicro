package http

import (
	"net/http"
)

// CallSettings allow fine-grained control over how calls are made.
type CallSettings struct {
	// Content-Type
	contentType string
	// custom header
	header http.Header
	// Path overwrite api call
	Path string
	// no auth
	noAuth bool
}

func DefaultCallOption(path string) CallSettings {
	return CallSettings{
		contentType: "application/json",
		Path:        path,
		header:      make(http.Header),
		noAuth:      false,
	}
}

// CallOption is an option used by Invoke to control behaviors of RPC calls.
// CallOption works by modifying relevant fields of CallSettings.
type CallOption func(cs *CallSettings)

// WithCoContentType use encoding.MIMExxx
func WithCoContentType(contentType string) CallOption {
	return func(cs *CallSettings) {
		cs.contentType = contentType
	}
}

// WithCoPath
func WithCoPath(path string) CallOption {
	return func(cs *CallSettings) {
		cs.Path = path
	}
}

// WithCoHeader
func WithCoHeader(k, v string) CallOption {
	return func(cs *CallSettings) {
		cs.header.Add(k, v)
	}
}

// WithCoNoAuth
func WithCoNoAuth() CallOption {
	return func(cs *CallSettings) {
		cs.noAuth = true
	}
}
