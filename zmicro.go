package zmicro

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap/zapcore"

	"github.com/zmicro-team/zmicro/core/config"
	"github.com/zmicro-team/zmicro/core/log"
	"github.com/zmicro-team/zmicro/core/transport/http"
	"github.com/zmicro-team/zmicro/core/transport/rpc/server"
	"github.com/zmicro-team/zmicro/core/util/env"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"gopkg.in/natefinch/lumberjack.v2"
)

var cfgFile string

func init() {
	flag.StringVar(&cfgFile, "config", "config.yaml", "config file")
}

type App struct {
	opts       Options
	zc         *zconfig
	rpcServer  *server.Server
	httpServer *http.Server
}

type zconfig struct {
	App struct {
		Mode string
		Name string
	}
	Logger struct {
		Level      string `json:"level"`
		Filename   string `json:"filename"`
		MaxSize    int    `json:"maxSize"`
		MaxBackups int    `json:"maxBackups"`
		MaxAge     int    `json:"maxAge"`
		Compress   bool   `json:"compress"`
	}
	Http struct {
		Addr string
	}
	Rpc struct {
		Addr string
	}
	Tracer struct {
		Addr string
	}
	Registry struct {
		BasePath       string
		EtcdAddr       []string
		UpdateInterval int
	}
}

func New(opts ...Option) *App {
	options := newOptions(opts...)
	flag.Parse()
	_, err := os.Stat(cfgFile)
	if os.IsNotExist(err) {
		log.Fatal("config file not exists")
	}

	c := config.New(config.Path(cfgFile), config.Callbacks(options.ConfigCallbacks...))
	config.ResetDefault(c)

	zc := &zconfig{}
	if err = config.Default().Unmarshal(zc); err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(zc)
	if zc.App.Name == "" {
		log.Fatal("配置项app.name不能为空")
	}

	env.Set(zc.App.Mode)

	level, err := zapcore.ParseLevel(zc.Logger.Level)
	if err != nil {
		level = log.InfoLevel
	}
	if env.IsDevelop() {
		w := &lumberjack.Logger{
			Filename:   zc.Logger.Filename,
			MaxSize:    zc.Logger.MaxSize,
			MaxBackups: zc.Logger.MaxBackups,
			MaxAge:     zc.Logger.MaxAge,
			Compress:   zc.Logger.Compress,
		}
		l := log.NewTee([]io.Writer{os.Stderr, w}, level, log.WithCaller(true))
		log.ResetDefault(l)
	} else {
		w := &lumberjack.Logger{
			Filename:   zc.Logger.Filename,
			MaxSize:    zc.Logger.MaxSize,
			MaxBackups: zc.Logger.MaxBackups,
			MaxAge:     zc.Logger.MaxAge,
			Compress:   zc.Logger.Compress,
		}
		l := log.New(w, level, log.WithCaller(true))
		log.ResetDefault(l)
	}

	app := &App{
		opts: options,
		zc:   zc,
	}

	tracing := false
	if zc.Tracer.Addr != "" {
		setTracerProvider(zc.Tracer.Addr, zc.App.Name)
		tracing = true
	}

	if app.opts.InitRpcServer != nil {
		app.rpcServer = server.NewServer(
			server.Name(zc.App.Name),
			server.Addr(zc.Rpc.Addr),
			server.BasePath(zc.Registry.BasePath),
			server.UpdateInterval(zc.Registry.UpdateInterval),
			server.EtcdAddr(zc.Registry.EtcdAddr),
			server.Tracing(tracing),
		)
		app.rpcServer.Init(server.InitRpcServer(app.opts.InitRpcServer))
	}
	mode := "debug"
	if env.IsProduct() || env.IsStaging() {
		mode = "release"
	}
	if app.opts.InitHttpServer != nil {
		app.httpServer = http.NewServer(
			http.Name(zc.App.Name),
			http.Addr(zc.Http.Addr),
			http.Mode(mode),
			http.Tracing(tracing),
		)
		app.httpServer.Init(http.InitHttpServer(app.opts.InitHttpServer))
	}

	return app
}

func setTracerProvider(endpoint string, name string) *trace.TracerProvider {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		log.Fatal(err.Error())
	}
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp
}

func (a *App) Run() error {
	if a.opts.Before != nil {
		if err := a.opts.Before(); err != nil {
			return err
		}
	}
	if a.rpcServer != nil {
		if err := a.rpcServer.Start(); err != nil {
			return err
		}
	}

	if a.httpServer != nil {
		if err := a.httpServer.Start(); err != nil {
			return err
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	log.Infof("received signal %s", <-ch)

	if a.rpcServer != nil {
		_ = a.rpcServer.Stop()
	}

	if a.httpServer != nil {
		_ = a.httpServer.Stop()
	}

	return nil
}
