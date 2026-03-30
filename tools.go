//go:build tools

// Package tools exists solely to pin tool dependencies in go.mod.
// It is not included in normal builds.
package tools

import (
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

