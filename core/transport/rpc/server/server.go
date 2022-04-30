package server

import (
	"context"
	"go.opentelemetry.io/otel"
	"net"
	"strings"
	"time"

	"github.com/iobrother/zmicro/core/log"
	"github.com/iobrother/zmicro/core/util/addr"
	znet "github.com/iobrother/zmicro/core/util/net"
	etcdServerPlugin "github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
)

type Server struct {
	server *server.Server
	opts   Options
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)
	srv := &Server{
		opts: options,
	}
	srv.server = server.NewServer()
	return srv
}

func (s *Server) Init(opts ...Option) {
	for _, opt := range opts {
		opt(&s.opts)
	}
}

func (s *Server) Start(l net.Listener) error {
	a := l.Addr().String()
	if s.opts.Tracing {
		tracer := otel.Tracer("rpcx")
		p := serverplugin.NewOpenTelemetryPlugin(tracer, nil)
		s.server.Plugins.Add(p)
	}
	s.register(a)

	if s.opts.InitRpcServer != nil {
		if err := s.opts.InitRpcServer(s.server); err != nil {
			return err
		}
	}

	log.Infof("Server [RPCX] listening on %s", a)
	go func() {
		if err := s.server.ServeListener("tcp", l); err != nil {
			log.Fatal(err.Error())
		}
	}()
	return nil
}

func (s *Server) Stop() error {
	log.Info("Server [RPCX] stopping")
	return s.server.Shutdown(context.Background())
}

func (s *Server) register(a string) {
	if len(s.opts.EtcdAddr) == 0 {
		return
	}

	var err error
	var host, port string
	if cnt := strings.Count(a, ":"); cnt >= 1 {
		host, port, err = net.SplitHostPort(a)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	} else {
		host = a
	}

	address, err := addr.Extract(host)
	if err != nil {
		log.Fatal(err.Error())
	}

	if port != "" {
		address = znet.HostPort(address, port)
	}

	r := &etcdServerPlugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + address,
		EtcdServers:    s.opts.EtcdAddr,
		BasePath:       s.opts.BasePath,
		UpdateInterval: time.Duration(s.opts.UpdateInterval) * time.Second,
	}
	err = r.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	s.server.Plugins.Add(r)

	log.Infof("Registering server: %s", address)
}
