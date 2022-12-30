package client

import (
	etcd_client "github.com/rpcxio/rpcx-etcd/client"
	otelClient "github.com/rpcxio/rpcx-plugins/client/otel"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"go.opentelemetry.io/otel"

	"github.com/zmicro-team/zmicro/core/log"
)

type Client struct {
	opts    Options
	xClient client.XClient
}

func NewClient(opts ...Option) (*Client, error) {
	options := newOptions(opts...)

	c := &Client{opts: options}

	if len(c.opts.EtcdAddr) > 0 {
		d, err := etcd_client.NewEtcdV3Discovery(
			c.opts.BasePath,
			c.opts.ServiceName,
			c.opts.EtcdAddr,
			false,
			nil,
		)
		if err != nil {
			log.Errorf("NewClient error=%v", err)
			return nil, err
		}
		opt := client.DefaultOption
		opt.SerializeType = protocol.ProtoBuffer
		c.xClient = client.NewXClient(
			c.opts.ServiceName,
			client.Failtry,
			client.RoundRobin,
			d,
			opt,
		)
	} else {
		d, err := client.NewPeer2PeerDiscovery("tcp@"+c.opts.ServiceAddr, "")
		if err != nil {
			log.Errorf("NewClient error=%v", err)
			return nil, err
		}

		opt := client.DefaultOption
		opt.SerializeType = protocol.ProtoBuffer
		c.xClient = client.NewXClient(c.opts.ServiceName, client.Failtry, client.RoundRobin, d, opt)
	}

	if c.opts.Tracing {
		tracer := otel.Tracer("rpcx")
		p := otelClient.NewOpenTelemetryPlugin(tracer, nil)
		pc := client.NewPluginContainer()
		pc.Add(p)
		c.xClient.SetPlugins(pc)
	}

	return c, nil
}

func (c *Client) GetXClient() client.XClient {
	return c.xClient
}
