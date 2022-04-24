package server

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/iobrother/zmicro/core/config"
	"github.com/iobrother/zmicro/core/log"
	"github.com/iobrother/zmicro/core/util/addr"
	znet "github.com/iobrother/zmicro/core/util/net"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
)

type Server struct {
	registry *serverplugin.EtcdV3RegisterPlugin
	server   *server.Server
	conf     *rpcConfig
	opts     Options
}

type rpcConfig struct {
	BasePath       string
	UpdateInterval int
	EtcdAddr       []string
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)

	conf := rpcConfig{}
	if err := config.Scan("rpc", &conf); err != nil {
		log.Fatal(err.Error())
	}
	srv := &Server{
		opts: options,
		conf: &conf,
	}
	srv.server = server.NewServer()
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
	log.Infof("Server [RPCX] listening On %s", addr)
	s.register(addr)

	if s.opts.InitRpcServer != nil {
		s.opts.InitRpcServer(s.server)
	}

	go func() {
		if err := s.server.ServeListener("tcp", l); err != nil {
			log.Fatal(err.Error())
		}
	}()
	return nil
}

func (s *Server) Stop() error {
	log.Info("Server [RPCX] stopping")
	_ = s.server.Shutdown(context.Background())
	return nil
}

func (s *Server) register(address string) {
	if len(s.conf.EtcdAddr) == 0 {
		return
	}

	var err error
	var host, port string
	if cnt := strings.Count(address, ":"); cnt >= 1 {
		host, port, err = net.SplitHostPort(address)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	} else {
		host = address
	}

	a, err := addr.Extract(host)
	if err != nil {
		log.Fatal(err.Error())
	}

	if port != "" {
		a = znet.HostPort(a, port)
	}

	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + a,
		EtcdServers:    s.conf.EtcdAddr,
		BasePath:       s.conf.BasePath,
		UpdateInterval: time.Duration(s.conf.UpdateInterval) * time.Second,
	}
	err = r.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	s.server.Plugins.Add(r)

	s.registry = r
	log.Infof("Registering server: %s", a)
}
