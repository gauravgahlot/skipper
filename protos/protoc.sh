#!/usr/bin/env bash

set -e

for proto in event; do
	echo "Generating ${proto}.pb.go..."
	protoc -I ./ -I ./common/ -I "$GOPATH"/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis "${proto}/${proto}.proto" --go_out=plugins=grpc:./
done
goimports -w .
