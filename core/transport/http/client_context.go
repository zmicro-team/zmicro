package http

import "context"

type ctxCallOptionKey struct{}

// WithValueCallOption returns a new Context that carries value.
func WithValueCallOption(ctx context.Context, cs CallSettings) context.Context {
	return context.WithValue(ctx, ctxCallOptionKey{}, cs)
}

// FromValueCallOption returns the CallSettings value stored in ctx, if any.
func FromValueCallOption(ctx context.Context) (cs CallSettings, ok bool) {
	cs, ok = ctx.Value(ctxCallOptionKey{}).(CallSettings)
	return
}

// MustFromValueCallOption returns the CallSettings value stored in ctx.
func MustFromValueCallOption(ctx context.Context) CallSettings {
	cs, ok := ctx.Value(ctxCallOptionKey{}).(CallSettings)
	if !ok {
		panic("transport: must be set CallSettings into context but it is not!!!")
	}
	return cs
}
