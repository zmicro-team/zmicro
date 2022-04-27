package main

import (
	"context"
	"flag"
	"os"

	"github.com/iobrother/zmicro/core/config"
	"github.com/iobrother/zmicro/core/log"
	"github.com/iobrother/zmicro/core/transport/rpc/client"
	"github.com/iobrother/zmicro/examples/greeter/proto"
)

var cfgFile string

func init() {
	flag.StringVar(&cfgFile, "config", "config.yaml", "config file")
}

func main() {

	flag.Parse()
	_, err := os.Stat(cfgFile)
	if os.IsNotExist(err) {
		log.Fatal("config file not exists")
	}

	config.ResetDefault(config.New(config.Path(cfgFile)))

	c := client.NewClient(client.WithServiceName("Greeter"), client.WithServiceAddr("127.0.0.1:5188"))
	cli := proto.NewGreeterClient(c.GetXClient())

	args := &proto.HelloRequest{
		Name: "zmicro",
	}

	reply, err := cli.SayHello(context.Background(), args)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("reply: %s", reply.Message)
}
