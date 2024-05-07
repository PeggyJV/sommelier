#!/usr/bin/env bash

set -eo pipefail

echo "generating proto and gRPC gateway files..."
cd proto
buf mod update
cd ..
buf generate

# move proto files to the right places
xpath=$(head -n 1 go.mod | sed 's/^module //')
cp -r $xpath/* ./

echo "cleaning up..."
rm -rf github.com
go mod tidy

echo "done"
