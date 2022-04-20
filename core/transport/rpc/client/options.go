package client

import (
	"log"

	"github.com/iobrother/zmicro/core/config"
	etcd_client "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

type Client struct {
	opts    Options
	conf    *clientConfig
	xClient client.XClient
}

type clientConfig struct {
	BasePath string
	EtcdAddr []string
}

func NewClient(opts ...Option) *Client {
	options := newOptions(opts...)

	conf := clientConfig{}
	config.DefaultConfig.Scan("rpc", &conf)
	c := &Client{opts: options, conf: &conf}

	if len(c.conf.EtcdAddr) > 0 {
		d, err := etcd_client.NewEtcdV3Discovery(
			c.conf.BasePath,
			c.opts.ServiceName,
			c.conf.EtcdAddr,
			false,
			nil,
		)
		if err != nil {
			log.Fatal(err)
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
			log.Fatal(err)
		}

		opt := client.DefaultOption
		opt.SerializeType = protocol.ProtoBuffer
		c.xClient = client.NewXClient(c.opts.ServiceName, client.Failtry, client.RoundRobin, d, opt)
	}

	return c
}

func (c *Client) GetXClient() client.XClient {
	return c.xClient
}
