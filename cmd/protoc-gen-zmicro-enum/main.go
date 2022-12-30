package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

const version = "v0.0.1"

var showVersion = flag.Bool("version", false, "print the version and exit")
var customTemplate = flag.String("template", "", "use custom template")
var suffix = flag.String("suffix", ".mapping.pb.go", "use custom file suffix")

var merge = flag.Bool("merge", false, "merge in a file")
var filename = flag.String("filename", "", "filename when merge enabled")
var _package = flag.String("package", "", "package name when merge enabled")
var goPackage = flag.String("go_package", "", "go package when merge enabled")

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-zmicro-enum %v\n", version)
		return
	}

	protogen.Options{ParamFunc: flag.CommandLine.Set}.Run(runProtoGen)
}
