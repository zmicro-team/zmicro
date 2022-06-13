package transport

import (
	"context"
)

type Transporter interface {
	// Kind transporter
	// grpc
	// http
	Kind() Kind
	// FullPath Service full method or path
	FullPath() string
	// ClientIp client ip
	ClientIp() string
	// RequestHeader return transport request header
	// http: http.Header
	// grpc: metadata.MD
	RequestHeader() Header
	// ResponseHeader return transport response header
	// http: http.Header
	// grpc: metadata.MD
	ResponseHeader() Header
}

// Header is the storage medium used by a Header.
type Header interface {
	Get(key string) string
	Add(key, value string)
	Set(key string, value string)
	Keys() []string
	Clone() Header
}

// Kind defines the type of Transport
type Kind string

func (k Kind) String() string { return string(k) }

// Defines a set of transport kind
const (
	GRPC Kind = "grpc"
	HTTP Kind = "http"
)

type ctxTransportKey struct{}

// WithValueTransporter returns a new Context that carries value.
func WithValueTransporter(ctx context.Context, p Transporter) context.Context {
	return context.WithValue(ctx, ctxTransportKey{}, p)
}

// FromTransporter returns the Transporter value stored in ctx, if any.
func FromTransporter(ctx context.Context) (p Transporter, ok bool) {
	p, ok = ctx.Value(ctxTransportKey{}).(Transporter)
	return
}

// MustFromTransporter returns the Transporter value stored in ctx.
func MustFromTransporter(ctx context.Context) Transporter {
	p, ok := ctx.Value(ctxTransportKey{}).(Transporter)
	if !ok {
		panic("transport: must be set Transporter into context but it is not!!!")
	}
	return p
}
