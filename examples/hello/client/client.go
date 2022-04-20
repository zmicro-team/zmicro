package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/iobrother/zmicro/core/config"
	"github.com/iobrother/zmicro/core/transport/rpc/client"
	"github.com/iobrother/zmicro/examples/hello/proto"
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

	if config.DefaultConfig, err = config.NewConfig(config.Path(cfgFile)); err != nil {
		log.Fatal(err)
	}

	c := client.NewClient(client.WithServiceName("Greeter"), client.WithServiceAddr("127.0.0.1:5188"))
	cli := proto.NewGreeterClient(c.GetXClient())

	args := &proto.HelloRequest{
		Name: "zmicro",
	}

	reply, _ := cli.SayHello(context.Background(), args)
	log.Println("reply: ", reply.Message)
}
