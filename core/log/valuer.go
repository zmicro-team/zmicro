package log

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

func FromBinary(key string, vf func(context.Context) []byte) Valuer {
	return func(ctx context.Context) Field {
		return zap.Binary(key, vf(ctx))
	}
}
func FromBool(key string, vf func(context.Context) bool) Valuer {
	return func(ctx context.Context) Field {
		return zap.Bool(key, vf(ctx))
	}
}
func FromBoolp(key string, vf func(context.Context) *bool) Valuer {
	return func(ctx context.Context) Field {
		return zap.Boolp(key, vf(ctx))
	}
}
func FromByteString(key string, vf func(context.Context) []byte) Valuer {
	return func(ctx context.Context) Field {
		return zap.ByteString(key, vf(ctx))
	}
}
func FromComplex128(key string, vf func(context.Context) complex128) Valuer {
	return func(ctx context.Context) Field {
		return zap.Complex128(key, vf(ctx))
	}
}
func FromComplex128p(key string, vf func(context.Context) *complex128) Valuer {
	return func(ctx context.Context) Field {
		return zap.Complex128p(key, vf(ctx))
	}
}
func FromComplex64(key string, vf func(context.Context) complex64) Valuer {
	return func(ctx context.Context) Field {
		return zap.Complex64(key, vf(ctx))
	}
}
func FromComplex64p(key string, vf func(context.Context) *complex64) Valuer {
	return func(ctx context.Context) Field {
		return zap.Complex64p(key, vf(ctx))
	}
}
func FromFloat64(key string, vf func(context.Context) float64) Valuer {
	return func(ctx context.Context) Field {
		return zap.Float64(key, vf(ctx))
	}
}
func FromFloat64p(key string, vf func(context.Context) *float64) Valuer {
	return func(ctx context.Context) Field {
		return zap.Float64p(key, vf(ctx))
	}
}
func FromFloat32(key string, vf func(context.Context) float32) Valuer {
	return func(ctx context.Context) Field {
		return zap.Float32(key, vf(ctx))
	}
}
func FromFloat32p(key string, vf func(context.Context) *float32) Valuer {
	return func(ctx context.Context) Field {
		return zap.Float32p(key, vf(ctx))
	}
}
func FromInt(key string, vf func(context.Context) int) Valuer {
	return func(ctx context.Context) Field {
		return zap.Int(key, vf(ctx))
	}
}
func FromIntp(key string, vf func(context.Context) *int) Valuer {
	return func(ctx context.Context) Field {
		return zap.Intp(key, vf(ctx))
	}
}
func FromInt64(key string, vf func(context.Context) int64) Valuer {
	return func(ctx context.Context) Field {
		return zap.Int64(key, vf(ctx))
	}
}
func FromInt64p(key string, vf func(context.Context) *int64) Valuer {
	return func(ctx context.Context) Field {
		return zap.Int64p(key, vf(ctx))
	}
}
func FromInt32(key string, vf func(context.Context) int32) Valuer {
	return func(ctx context.Context) Field {
		return zap.Int32(key, vf(ctx))
	}
}
func FromInt32p(key string, vf func(context.Context) *int32) Valuer {
	return func(ctx context.Context) Field {
		return zap.Int32p(key, vf(ctx))
	}
}
func FromInt16(key string, vf func(context.Context) int16) Valuer {
	return func(ctx context.Context) Field {
		return zap.Int16(key, vf(ctx))
	}
}
func FromInt16p(key string, vf func(context.Context) *int16) Valuer {
	return func(ctx context.Context) Field {
		return zap.Int16p(key, vf(ctx))
	}
}
func FromInt8(key string, vf func(context.Context) int8) Valuer {
	return func(ctx context.Context) Field {
		return zap.Int8(key, vf(ctx))
	}
}
func FromInt8p(key string, vf func(context.Context) *int8) Valuer {
	return func(ctx context.Context) Field {
		return zap.Int8p(key, vf(ctx))
	}
}
func FromUint(key string, vf func(context.Context) uint) Valuer {
	return func(ctx context.Context) Field {
		return zap.Uint(key, vf(ctx))
	}
}
func FromUintp(key string, vf func(context.Context) *uint) Valuer {
	return func(ctx context.Context) Field {
		return zap.Uintp(key, vf(ctx))
	}
}
func FromUint64(key string, vf func(context.Context) uint64) Valuer {
	return func(ctx context.Context) Field {
		return zap.Uint64(key, vf(ctx))
	}
}
func FromUint64p(key string, vf func(context.Context) *uint64) Valuer {
	return func(ctx context.Context) Field {
		return zap.Uint64p(key, vf(ctx))
	}
}
func FromUint32(key string, vf func(context.Context) uint32) Valuer {
	return func(ctx context.Context) Field {
		return zap.Uint32(key, vf(ctx))
	}
}
func FromUint32p(key string, vf func(context.Context) *uint32) Valuer {
	return func(ctx context.Context) Field {
		return zap.Uint32p(key, vf(ctx))
	}
}
func FromUint16(key string, vf func(context.Context) uint16) Valuer {
	return func(ctx context.Context) Field {
		return zap.Uint16(key, vf(ctx))
	}
}
func FromUint16p(key string, vf func(context.Context) *uint16) Valuer {
	return func(ctx context.Context) Field {
		return zap.Uint16p(key, vf(ctx))
	}
}
func FromUint8(key string, vf func(context.Context) uint8) Valuer {
	return func(ctx context.Context) Field {
		return zap.Uint8(key, vf(ctx))
	}
}
func FromUint8p(key string, vf func(context.Context) *uint8) Valuer {
	return func(ctx context.Context) Field {
		return zap.Uint8p(key, vf(ctx))
	}
}
func FromString(key string, vf func(context.Context) string) Valuer {
	return func(ctx context.Context) Field {
		return zap.String(key, vf(ctx))
	}
}
func FromStringp(key string, vf func(context.Context) *string) Valuer {
	return func(ctx context.Context) Field {
		return zap.Stringp(key, vf(ctx))
	}
}
func FromUintptr(key string, vf func(context.Context) uintptr) Valuer {
	return func(ctx context.Context) Field {
		return zap.Uintptr(key, vf(ctx))
	}
}
func FromUintptrp(key string, vf func(context.Context) *uintptr) Valuer {
	return func(ctx context.Context) Field {
		return zap.Uintptrp(key, vf(ctx))
	}
}
func FromReflect(key string, vf func(context.Context) any) Valuer {
	return func(ctx context.Context) Field {
		return zap.Reflect(key, vf(ctx))
	}
}
func FromStringer(key string, vf func(context.Context) fmt.Stringer) Valuer {
	return func(ctx context.Context) Field {
		return zap.Stringer(key, vf(ctx))
	}
}
func FromTime(key string, vf func(context.Context) time.Time) Valuer {
	return func(ctx context.Context) Field {
		return zap.Time(key, vf(ctx))
	}
}
func FromTimep(key string, vf func(context.Context) *time.Time) Valuer {
	return func(ctx context.Context) Field {
		return zap.Timep(key, vf(ctx))
	}
}
func FromDuration(key string, vf func(context.Context) time.Duration) Valuer {
	return func(ctx context.Context) Field {
		return zap.Duration(key, vf(ctx))
	}
}
func FromDurationp(key string, vf func(context.Context) *time.Duration) Valuer {
	return func(ctx context.Context) Field {
		return zap.Durationp(key, vf(ctx))
	}
}
func FromAny(key string, vf func(context.Context) any) Valuer {
	return func(ctx context.Context) Field {
		return zap.Any(key, vf(ctx))
	}
}
func FromNamespace(key string) Valuer {
	field := zap.Namespace(key)
	return func(ctx context.Context) Field {
		return field
	}
}
func FromStack(key string) Valuer {
	return func(ctx context.Context) Field {
		return zap.Stack(key)
	}
}
func FromStackSkip(key string, skip int) Valuer {
	return func(ctx context.Context) Field {
		return zap.StackSkip(key, skip)
	}
}

func ImmutBinary(key string, v []byte) Valuer {
	field := zap.Binary(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutBool(key string, v bool) Valuer {
	field := zap.Bool(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutBoolp(key string, v *bool) Valuer {
	field := zap.Boolp(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutByteString(key string, v []byte) Valuer {
	field := zap.ByteString(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutComplex128(key string, v complex128) Valuer {
	field := zap.Complex128(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutComplex128p(key string, v *complex128) Valuer {
	field := zap.Complex128p(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutComplex64(key string, v complex64) Valuer {
	field := zap.Complex64(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutComplex64p(key string, v *complex64) Valuer {
	field := zap.Complex64p(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutFloat64(key string, v float64) Valuer {
	field := zap.Float64(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutFloat64p(key string, v *float64) Valuer {
	field := zap.Float64p(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutFloat32(key string, v float32) Valuer {
	field := zap.Float32(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutFloat32p(key string, v *float32) Valuer {
	field := zap.Float32p(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutInt(key string, v int) Valuer {
	field := zap.Int(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutIntp(key string, v *int) Valuer {
	field := zap.Intp(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutInt64(key string, v int64) Valuer {
	field := zap.Int64(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutInt64p(key string, v *int64) Valuer {
	field := zap.Int64p(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutInt32(key string, v int32) Valuer {
	field := zap.Int32(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutInt32p(key string, v *int32) Valuer {
	field := zap.Int32p(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutInt16(key string, v int16) Valuer {
	field := zap.Int16(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutInt16p(key string, v *int16) Valuer {
	field := zap.Int16p(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutInt8(key string, v int8) Valuer {
	field := zap.Int8(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutInt8p(key string, v *int8) Valuer {
	field := zap.Int8p(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}

func ImmutUint(key string, v uint) Valuer {
	field := zap.Uint(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutUintp(key string, v *uint) Valuer {
	field := zap.Uintp(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutUint64(key string, v uint64) Valuer {
	field := zap.Uint64(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutUint64p(key string, v *uint64) Valuer {
	field := zap.Uint64p(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutUint32(key string, v uint32) Valuer {
	field := zap.Uint32(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutUint32p(key string, v *uint32) Valuer {
	field := zap.Uint32p(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}

func ImmutUint16(key string, v uint16) Valuer {
	field := zap.Uint16(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutUint16p(key string, v *uint16) Valuer {
	field := zap.Uint16p(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutUint8(key string, v uint8) Valuer {
	field := zap.Uint8(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutUint8p(key string, v *uint8) Valuer {
	field := zap.Uint8p(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}

func ImmutString(key string, v string) Valuer {
	field := zap.String(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutStringp(key string, v *string) Valuer {
	field := zap.Stringp(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutUintptr(key string, v uintptr) Valuer {
	field := zap.Uintptr(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutUintptrp(key string, v *uintptr) Valuer {
	field := zap.Uintptrp(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutReflect(key string, v any) Valuer {
	field := zap.Reflect(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutStringer(key string, v fmt.Stringer) Valuer {
	field := zap.Stringer(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutTime(key string, v time.Time) Valuer {
	field := zap.Time(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutTimep(key string, v *time.Time) Valuer {
	field := zap.Timep(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutDuration(key string, v time.Duration) Valuer {
	field := zap.Duration(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutDurationp(key string, v *time.Duration) Valuer {
	field := zap.Durationp(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}
func ImmutAny(key string, v any) Valuer {
	field := zap.Any(key, v)
	return func(ctx context.Context) Field {
		return field
	}
}

/**************************** helper ******************************************/
func caller(depth int) (file string, line int) {
	d := depth
	_, file, line, _ = runtime.Caller(d)
	if strings.LastIndex(file, "/log/log.go") > 0 {
		d++
		_, file, line, _ = runtime.Caller(d)
	}
	if strings.LastIndex(file, "/log/default.go") > 0 {
		d++
		_, file, line, _ = runtime.Caller(d)
	}
	return file, line
}

// Caller returns a Valuer that returns a pkg/file:line description of the caller.
func Caller(depth int) Valuer {
	return func(context.Context) Field {
		file, line := caller(depth)
		idx := strings.LastIndexByte(file, '/')
		return zap.String("caller", file[idx+1:]+":"+strconv.Itoa(line))
	}
}

// File returns a Valuer that returns a pkg/file:line description of the caller.
func File(depth int) Valuer {
	return func(context.Context) Field {
		file, line := caller(depth)
		return zap.String("file", file+":"+strconv.Itoa(line))
	}
}

// Package returns a Valuer that returns a immutable Valuer which key is pkg
func Package(v string) Valuer {
	return ImmutString("pkg", v)
}

func App(v string) Valuer {
	return ImmutString("app", v)
}
func Component(v string) Valuer {
	return ImmutString("component", v)
}
func Module(v string) Valuer {
	return ImmutString("module", v)
}
func Unit(v string) Valuer {
	return ImmutString("unit", v)
}
func Kind(v string) Valuer {
	return ImmutString("kind", v)
}
func Type(v string) Valuer {
	return ImmutString("type", v)
}
func TraceId(f func(c context.Context) string) Valuer {
	return FromString("traceId", f)
}
func RequestId(f func(c context.Context) string) Valuer {
	return FromString("requestId", f)
}
