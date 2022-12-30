package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const version = "0.1.0"

var showVersion = flag.Bool("version", false, "print the version and exit")
var errorsPackage = flag.String("epk", "github.com/zmicro-team/zmicro/core/errors", "errors core package in your project")

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-zmicro-errno %v\n", version)
		return
	}

	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(runPortoGen)
}
