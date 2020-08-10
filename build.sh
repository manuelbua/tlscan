#!/usr/bin/env bash

version="$(git describe --tags --long)"
go build -o bin/tlscan -ldflags="-X 'github.com/manuelbua/tlscan/pkg/tlscan.version=${version}'" cmd/tlscan/main.go
if [ "$?" = 0 ]; then
  echo "Built version ${version}"
fi