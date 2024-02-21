// Code generated by protoc-gen-zmicro-resty. DO NOT EDIT.
// versions:
// - protoc-gen-zmicro-resty v0.3.0
// - protoc                v4.24.0
// source: api/hello.proto

package api

import (
	context "context"
	errors "errors"
	http "github.com/zmicro-team/zmicro/core/transport/http"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = errors.New
var _ = context.TODO
var _ = http.NewClient

// GreeterHTTPClient
type GreeterHTTPClient interface {
	// SayHello
	SayHello(context.Context, *HelloRequest, ...http.CallOption) (*HelloReply, error)
}

type GreeterHTTPClientImpl struct {
	cc *http.Client
}

// NewGreeterHTTPClient
func NewGreeterHTTPClient(c *http.Client) GreeterHTTPClient {
	return &GreeterHTTPClientImpl{
		cc: c,
	}
}

// SayHello
func (c *GreeterHTTPClientImpl) SayHello(ctx context.Context, req *HelloRequest, opts ...http.CallOption) (*HelloReply, error) {
	var err error
	var resp HelloReply

	settings := c.cc.CallSetting("/hello/{name}", opts...)
	path := c.cc.EncodeURL(settings.Path, req, true)
	err = c.cc.Invoke2(ctx, "GET", path, nil, &resp, settings)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
