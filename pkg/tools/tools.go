//go:build tools

package tools

// These imports only exist to keep go.mod entries for packages that are referenced in BUILD files,
// but not in Go code.

import (
	_ "github.com/gogo/protobuf/proto"
	_ "google.golang.org/grpc"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
