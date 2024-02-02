package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "1.0.0"

var (
	showVersion = flag.Bool("version", false, "show http generator version")
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-http version: %v\n", version)
		return
	}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}

			generate(gen, f)
		}
		return nil
	})
}

func addHeader(gen *protogen.Plugin, f *protogen.File, filename string) *protogen.GeneratedFile {
	g := gen.NewGeneratedFile(filename, f.GoImportPath)
	g.P("// Code generated by protoc-gen-pallas-http. DO NOT EDIT.")
	g.P("//")
	g.P("// proto-gen-pallas-http version: ", version)
	g.P("// protoc version: ", getProtocVersion(gen))
	if f.Proto.GetOptions().GetDeprecated() {
		g.P("// ", f.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source file: ", f.Desc.Path())
	}
	g.P()
	g.P("package ", f.GoPackageName)
	g.P()
	return g
}