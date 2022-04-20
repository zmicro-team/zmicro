package zmicro

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/iobrother/zmicro/core/config"
	"github.com/iobrother/zmicro/core/transport/http"
	"github.com/iobrother/zmicro/core/transport/rpc/server"
)

var cfgFile string

func init() {
	flag.StringVar(&cfgFile, "config", "config.yaml", "config file")
}

type App struct {
	opts       Options
	conf       *appConfig
	rpcServer  *server.Server
	httpServer *http.Server
}

type appConfig struct {
	Name string
	Addr string
}

func New(opts ...Option) *App {
	options := newOptions(opts...)
	flag.Parse()
	_, err := os.Stat(cfgFile)
	if os.IsNotExist(err) {
		log.Fatal("config file not exists")
	}

	if config.DefaultConfig, err = config.NewConfig(config.Path(cfgFile)); err != nil {
		log.Fatal(err)
	}

	conf := appConfig{}
	if err = config.Scan("app", &conf); err != nil {
		log.Fatal(err)
	}
	app := &App{opts: options, conf: &conf, rpcServer: server.NewServer()}
	app.rpcServer.Init(server.WithInitRpcServer(app.opts.InitRpcServer))
	return app
}

func (a *App) Run() error {

	l, err := net.Listen("tcp", a.conf.Addr)
	if err != nil {
		return err
	}

	if a.rpcServer != nil {
		if err = a.rpcServer.Start(l); err != nil {
			return err
		}
	}

	if a.httpServer != nil {
		if err = a.httpServer.Start(l); err != nil {
			return err
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	log.Printf("received signal %s", <-ch)

	if a.rpcServer != nil {
		a.rpcServer.Stop()
	}

	if a.httpServer != nil {
		a.httpServer.Stop()
	}

	return nil
}
