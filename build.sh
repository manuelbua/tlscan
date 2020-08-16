#!/usr/bin/env bash

version="$(git describe --tags --long)"

go build -o bin/tlscan -ldflags="-X 'github.com/manuelbua/tlscan/pkg/runner.version=${version}'" cmd/tlscan/main.go
GOOS=linux GOARCH=amd64 go build -o bin/tlscan_linux64 -ldflags="-X 'github.com/manuelbua/tlscan/pkg/runner.version=${version}'" cmd/tlscan/main.go

if [ "$?" = 0 ]; then
  echo "Built version ${version}"
fi
