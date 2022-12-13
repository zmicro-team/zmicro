package logging

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
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

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		var start = time.Now()
		var writer *rspWriter
		var buf bytes.Buffer
		var body string
		c.Request.Body = ioutil.NopCloser(io.TeeReader(c.Request.Body, &buf))
		if _, err := ioutil.ReadAll(c.Request.Body); err != nil {
			c.Abort()
			return
		}

		body = buf.String()
		c.Request.Body = ioutil.NopCloser(&buf)
		writer = &rspWriter{c.Writer, &bytes.Buffer{}}
		c.Writer = writer

		defer func() {
			// if c.Request.RequestURI
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
				"body": func() interface{} {
					if c.Request.MultipartForm != nil {
						return "MultipartForm File"
					} else {
						return body
					}
				}(),
			}))

			if env.Get() == env.Develop {
				fields = append(fields, zap.Any("rsp", map[string]interface{}{
					"header": copyHeader(c.Writer.Header()),
					"body":   writer.body.String(),
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
