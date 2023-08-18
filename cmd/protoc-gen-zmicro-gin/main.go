package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const version = "0.2.0"

var args = struct {
	ShowVersion            bool
	Omitempty              bool
	AllowDeleteBody        bool
	AllowEmptyPatchBody    bool
	RpcMode                string
	UseEncoding            bool
	DisableErrorBadRequest bool
	DisableClient          bool
}{
	ShowVersion:            false,
	Omitempty:              true,
	AllowDeleteBody:        false,
	AllowEmptyPatchBody:    false,
	RpcMode:                "rpcx",
	UseEncoding:            false,
	DisableErrorBadRequest: false,
	DisableClient:          true,
}

func init() {
	flag.BoolVar(&args.ShowVersion, "version", false, "print the version and exit")
	flag.BoolVar(&args.Omitempty, "omitempty", true, "omit if google.api is empty")
	flag.BoolVar(&args.AllowDeleteBody, "allow_delete_body", false, "allow delete body")
	flag.BoolVar(&args.AllowEmptyPatchBody, "allow_empty_patch_body", false, "allow empty patch body")
	flag.StringVar(&args.RpcMode, "rpc_mode", "rpcx", "rpc mode, default use rpcx rpc, options: rpcx,official")
	flag.BoolVar(&args.UseEncoding, "use_encoding", false, "use the framework encoding")
	flag.BoolVar(&args.DisableErrorBadRequest, "disable_error_bad_request", false, "disable error bad request")
	flag.BoolVar(&args.DisableClient, "disable_client", true, "disable use client")
}

func main() {
	flag.Parse()
	if args.ShowVersion {
		fmt.Printf("protoc-gen-zmicro-gin %v\n", version)
		return
	}

	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(runProtoGen)
}
