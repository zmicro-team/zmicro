package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const version = "0.3.0"

var args = struct {
	ShowVersion         bool
	Omitempty           bool
	AllowDeleteBody     bool
	AllowEmptyPatchBody bool
	UseInvoke2          bool
}{
	ShowVersion:         false,
	Omitempty:           true,
	AllowDeleteBody:     false,
	AllowEmptyPatchBody: false,
	UseInvoke2:          false,
}

func init() {
	flag.BoolVar(&args.ShowVersion, "version", false, "print the version and exit")
	flag.BoolVar(&args.Omitempty, "omitempty", true, "omit if google.api is empty")
	flag.BoolVar(&args.AllowDeleteBody, "allow_delete_body", false, "allow delete body")
	flag.BoolVar(&args.AllowEmptyPatchBody, "allow_empty_patch_body", false, "allow empty patch body")
	flag.BoolVar(&args.UseInvoke2, "use_invoke2", false, "use invoke2")
}

func main() {
	flag.Parse()
	if args.ShowVersion {
		fmt.Printf("protoc-gen-zmicro-resty %v\n", version)
		return
	}

	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(runProtoGen)
}
