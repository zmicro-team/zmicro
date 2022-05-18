# zmicro

## 文档

https://iobrother.com/

## 简介

zmicro 是一套微服务开发解决方案，旨在帮助中小企业与广大 go 爱好者打造一套可落地的微服务方案

zmicro 集成了流行的 web 框架 [gin](https://github.com/gin-gonic/gin) 与 极简的 rpc 框架 [rpcx](https://github.com/smallnest/rpcx)

## 目标

- 极简：简单易学，易于开发与维护
- 效率：通过工具生成 gin 代码，rpcx 代码，错误代码，以及 API 文档，提高开发效率
- 性能：WEB框架 gin 与 RPC 框架 rpcx 在性能上处于业界领先

## 快速开始

### proto文件

```protobuf
syntax = "proto3";

option go_package = "github.com/zmicro-team/zmicro/examples/proto";

package proto;

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

### 安装代码生成插件

```bash
go install github.com/gogo/protobuf/protoc-gen-gofast@latest
go install github.com/rpcxio/protoc-gen-rpcx@latest
```

protoc-gen-rpcx以上方法安装不成功，是因为作者还没发布新版本，请直接下载代码，编译安装

```bash
protoc -I. -I${GOPATH}/src \
  --gofast_out=. --gofast_opt=paths=source_relative \
  --rpcx_out=. --rpcx_opt=paths=source_relative *.proto
```

上述命令生成了 hello.pb.go 与 hello.rpcx.pb.go 两个文件。 hello.pb.go 文件是由protoc-gen-gofast插件生成的， 当然你也可以选择官方的protoc-gen-go插件来生成。 hello.rpcx.pb.go 是由protoc-gen-rpcx插件生成的，它包含服务端的一个骨架， 以及客户端的代码。

### 服务端配置文件

```yaml
app:
  name: "example"
rpc:
  addr: ":5188"

```

### 服务端代码

```go
package main

import (
	"context"
	"fmt"

	"github.com/zmicro-team/zmicro"
	"github.com/zmicro-team/zmicro/core/log"
	"github.com/zmicro-team/zmicro/examples/proto"
	"github.com/smallnest/rpcx/server"
)

func main() {
	app := zmicro.New(zmicro.InitRpcServer(InitRpcServer))

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func InitRpcServer(s *server.Server) error {
	if err := s.RegisterName("Greeter", &GreeterImpl{}, ""); err != nil {
		return err
	}
	return nil
}

type GreeterImpl struct{}

func (s *GreeterImpl) SayHello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloReply) error {
	*rsp = proto.HelloReply{
		Message: fmt.Sprintf("hello %s!", req.Name),
	}

	return nil
}

```

### 客户端代码

```go
package main

import (
	"context"

	"github.com/zmicro-team/zmicro/core/log"
	"github.com/zmicro-team/zmicro/core/transport/rpc/client"
	"github.com/zmicro-team/zmicro/examples/proto"
)

func main() {
	c := client.NewClient(client.WithServiceName("Greeter"), client.WithServiceAddr("127.0.0.1:5188"))
	cli := proto.NewGreeterClient(c.GetXClient())

	req := &proto.HelloRequest{
		Name: "zmicro",
	}

	rsp, err := cli.SayHello(context.Background(), req)
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Infof("reply: %s", rsp.Message)
}
```

## 启动服务器

```bash
go run server.go
```

## 启动客户端

```bash
go run client.go
```

输出

```
{"level":"info","ts":"2022-05-02T16:34:17.754+0800","caller":"log/log.go:59","msg":"reply: hello zmicro!"}
```

## 源码地址

https://github.com/zmicro-team/zmicro/tree/master/examples/greeter
