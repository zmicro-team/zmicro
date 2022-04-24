package errors

import (
	"fmt"
)

// IsBadRequest determines if err is an error which indicates a BadRequest error.
// It supports wrapped errors.
func IsBadRequest(err error) bool {
	return Code(err) == 400
}

// ErrBadRequest new BadRequest error that is mapped to a 400 response.
func ErrBadRequest(detail string) *Error {
	return New(400, "请求参数错误", detail)
}

// ErrBadRequestf new BadRequest error that is mapped to a 400 response.
func ErrBadRequestf(format string, args ...any) *Error {
	return New(400, "请求参数错误", fmt.Sprintf(format, args...))
}

// IsUnauthorized determines if err is an error which indicates a Unauthorized error.
// It supports wrapped errors.
func IsUnauthorized(err error) bool {
	return Code(err) == 401
}

// ErrUnauthorized new Unauthorized error that is mapped to a 401 response.
func ErrUnauthorized(detail string) *Error {
	return Newf(401, "未授权", detail)
}

// ErrUnauthorizedf new Unauthorized error that is mapped to a 401 response.
func ErrUnauthorizedf(format string, args ...any) *Error {
	return Newf(401, "未授权", fmt.Sprintf(format, args...))
}

// IsForbidden determines if err is an error which indicates a Forbidden error.
// It supports wrapped errors.
func IsForbidden(err error) bool {
	return Code(err) == 403
}

// ErrForbidden new Forbidden error that is mapped to a 403 response.
func ErrForbidden(detail string) *Error {
	return Newf(403, "禁止访问", detail)
}

// ErrForbiddenf new Forbidden error that is mapped to a 403 response.
func ErrForbiddenf(format string, args ...any) *Error {
	return Newf(403, "禁止访问", fmt.Sprintf(format, args...))
}

// IsNotFound determines if err is an error which indicates an NotFound error.
// It supports wrapped errors.
func IsNotFound(err error) bool {
	return Code(err) == 404
}

// ErrNotFound new NotFound error that is mapped to a 404 response.
func ErrNotFound(detail string) *Error {
	return Newf(404, "没有找到,已删除或不存在", detail)
}

// ErrNotFoundf new NotFound error that is mapped to a 404 response.
func ErrNotFoundf(format string, args ...any) *Error {
	return Newf(404, "没有找到,已删除或不存在", fmt.Sprintf(format, args...))
}

// IsConflict determines if err is an error which indicates a Conflict error.
// It supports wrapped errors.
func IsConflict(err error) bool {
	return Code(err) == 409
}

// ErrConflict new Conflict error that is mapped to a 409 response.
func ErrConflict(detail string) *Error {
	return Newf(409, "资源冲突", detail)
}

// ErrConflictf new Conflict error that is mapped to a 409 response.
func ErrConflictf(format string, args ...any) *Error {
	return Newf(409, "资源冲突", fmt.Sprintf(format, args...))
}

func IsInternalServer(err error) bool {
	return Code(err) == 500
}
func ErrInternalServer(detail string) *Error {
	return New(500, "服务器错误", detail)
}
func ErrInternalServerf(format string, args ...any) *Error {
	return New(500, "服务器错误", fmt.Sprintf(format, args...))
}

// IsServiceUnavailable determines if err is an error which indicates a Unavailable error.
// It supports wrapped errors.
func IsServiceUnavailable(err error) bool {
	return Code(err) == 503
}

// ErrServiceUnavailable new ServiceUnavailable error that is mapped to a HTTP 503 response.
func ErrServiceUnavailable(detail string) *Error {
	return Newf(503, "服务器不可用", detail)
}

// ErrServiceUnavailablef new ServiceUnavailable error that is mapped to a HTTP 503 response.
func ErrServiceUnavailablef(format string, args ...any) *Error {
	return Newf(503, "服务器不可用", fmt.Sprintf(format, args...))
}

// IsGatewayTimeout determines if err is an error which indicates a GatewayTimeout error.
// It supports wrapped errors.
func IsGatewayTimeout(err error) bool {
	return Code(err) == 504
}

// ErrGatewayTimeout new GatewayTimeout error that is mapped to a HTTP 504 response.
func ErrGatewayTimeout(detail string) *Error {
	return Newf(504, "网关超时", detail)
}

// ErrGatewayTimeoutf new GatewayTimeout error that is mapped to a HTTP 504 response.
func ErrGatewayTimeoutf(format string, args ...any) *Error {
	return Newf(504, "网关超时", fmt.Sprintf(format, args...))
}

// IsClientClosed determines if err is an error which indicates a IsClientClosed error.
// It supports wrapped errors.
func IsClientClosed(err error) bool {
	return Code(err) == 499
}

// ErrClientClosed new ClientClosed error that is mapped to a HTTP 499 response.
func ErrClientClosed(message string) *Error {
	return Newf(499, "客户端关闭", message)
}

// ErrClientClosedf new ClientClosed error that is mapped to a HTTP 499 response.
func ErrClientClosedf(format string, args ...any) *Error {
	return Newf(499, "客户端关闭", fmt.Sprintf(format, args...))
}
