package logging

import (
	"bytes"
	"context"
	"io"
	"mime"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/zmicro-team/zmicro/core/log"
	"github.com/zmicro-team/zmicro/core/util/env"
)

var (
	slowThreshold = time.Millisecond * 500
)

type rspWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w rspWriter) Write(b []byte) (int, error) {
	n, err := w.body.Write(b)
	if err != nil {
		return n, err
	}

	return w.ResponseWriter.Write(b)
}

func (w rspWriter) WriteString(s string) (int, error) {
	n, err := w.body.WriteString(s)
	if err != nil {
		return n, err
	}

	return w.ResponseWriter.WriteString(s)
}

func copyHeader(header http.Header) http.Header {
	h := http.Header{}
	for k, v := range header {
		h[k] = v
	}
	return h
}

func skipRequestBody(c *gin.Context) bool {
	// skip request body rule
	// skip MultipartForm file
	v := c.Request.Header.Get("Content-Type")
	d, params, err := mime.ParseMediaType(v)
	if err != nil || !(d == "multipart/form-data" || d == "multipart/mixed") {
		return false
	}
	_, ok := params["boundary"]
	return ok
}

func skipResponseBody(c *gin.Context) bool {
	// TODO: add skip response body rule
	return false
}

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		var start = time.Now()

		respBodyBuf := &bytes.Buffer{}
		reqBody := "skip request body"

		if !skipRequestBody(c) {
			buf := &bytes.Buffer{}
			c.Request.Body = io.NopCloser(io.TeeReader(c.Request.Body, buf))
			if _, err := io.ReadAll(c.Request.Body); err != nil {
				c.Abort()
				return
			}

			reqBody = buf.String()
			c.Request.Body = io.NopCloser(buf)
		}
		c.Writer = &rspWriter{c.Writer, respBodyBuf}

		defer func() {
			if checkPrefix(c.Request.RequestURI, "/swagger") {
				return
			}
			fields := make([]zap.Field, 0, 14)
			traceId := getTraceId(c.Request.Context())
			if traceId != "" {
				fields = append(fields, zap.String("trace_id", traceId))
			}
			fields = append(fields, zap.String("type", "http"))
			fields = append(fields, zap.Int("status", c.Writer.Status()))
			fields = append(fields, zap.String("method", c.Request.Method))
			fields = append(fields, zap.String("route", c.FullPath()))
			fields = append(fields, zap.String("target", c.Request.RequestURI))
			fields = append(fields, zap.String("client_ip", c.ClientIP()))
			fields = append(fields, zap.String("peer", c.Request.RemoteAddr))
			fields = append(fields, zap.Int("size", c.Writer.Size()))
			duration := time.Since(start)
			fields = append(fields, zap.Duration("duration", duration))
			fields = append(fields, zap.Any("req", map[string]interface{}{
				"header": copyHeader(c.Request.Header),
				"body":   reqBody,
			}))

			if env.Get() == env.Develop {
				respBody := "skip response body"
				if !skipResponseBody(c) {
					respBody = respBodyBuf.String()
				}
				fields = append(fields, zap.Any("rsp", map[string]interface{}{
					"header": copyHeader(c.Writer.Header()),
					"body":   respBody,
				}))
			}
			// slow log
			if duration > slowThreshold {
				log.Warn("slow", fields...)
			}
			log.Info("access", fields...)
		}()

		c.Next()
	}
}

func getTraceId(ctx context.Context) string {
	if sc := trace.SpanContextFromContext(ctx); sc.HasTraceID() {
		return sc.SpanID().String()
	}
	return ""
}

func checkPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

// Recovery returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs but the stack info is too large.
func Recovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					log.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.ByteString("request", httpRequest),
					)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error))
					c.Abort()
					return
				}

				fields := make([]zap.Field, 0, 3)
				fields = append(fields,
					zap.Any("error", err),
					zap.ByteString("request", httpRequest),
				)
				if stack {
					fields = append(fields, zap.ByteString("stack", debug.Stack()))
				}
				log.Error("recovery from panic", fields...)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
