#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"


go generate "$DIR/main.go"
GOOS=js GOARCH=wasm go build -o "$DIR/main.wasm" "$DIR/main.go"
