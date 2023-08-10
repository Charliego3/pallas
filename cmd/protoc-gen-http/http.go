package main

import "google.golang.org/protobuf/compiler/protogen"

// generate generates a _grpc.http.go file containing gRPC service definitions.
func generate(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}

	return nil
}
