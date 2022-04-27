package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iobrother/zmicro/core/config"
	"github.com/iobrother/zmicro/core/log"
)

type Server struct {
	opts Options
	conf *httpConfig
	*gin.Engine
	server *http.Server
}

type httpConfig struct {
	Name string
	Mode string
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)
	conf := httpConfig{}
	if err := config.Scan("rpc", &conf); err != nil {
		log.Fatal(err.Error())
	}

	srv := &Server{
		opts: options,
		conf: &conf,
	}

	gin.SetMode(conf.Mode)
	r := gin.New()
	srv.Engine = r
	srv.server = &http.Server{Handler: srv.Engine}
	return srv
}

func (s *Server) Init(opts ...Option) error {
	for _, opt := range opts {
		opt(&s.opts)
	}
	return nil
}

func (s *Server) Start(l net.Listener) error {
	addr := l.Addr().String()
	log.Infof("Server [GIN] listening on %s", addr)
	if s.opts.InitHttpServer != nil {
		s.opts.InitHttpServer(s.Engine)
	}

	go func() {
		if err := s.server.Serve(l); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Fatal(err.Error())
			}
		}
	}()
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	_ = s.server.Shutdown(ctx)
	return nil
}
