#!/usr/bin/env bash

set -eo pipefail

ROOT=$(git rev-parse --show-toplevel 2>/dev/null)

cd $ROOT/proto
echo "generating proto and gRPC gateway files..."
buf generate --template buf.gen.gogo.yaml
cd ..

# move proto files to the right places
xpath=$(head -n 1 go.mod | sed 's/^module //')
cp -r $xpath/* ./

echo "cleaning up..."
go mod tidy
rm -rf github.com

echo "done"
