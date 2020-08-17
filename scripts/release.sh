#!/usr/bin/env bash

basepath=$(dirname $0)
version_pkg=""
main_go="${basepath}/../cmd/tlscan/main.go"
ldflags_sh="go-version-ldflags.sh"
outbin="${basepath}/../bin/tlscan"

go build -o "${outbin}" -ldflags="$(${ldflags_sh} "${version_pkg}")" "${main_go}"
GOOS=linux GOARCH=amd64 go build -o "${outbin}_linux64" -ldflags="$(${ldflags_sh} "${version_pkg}")" "${main_go}"

#upx --brute "${outbin}"
#upx --brute "${outbin}_linux64"

"${outbin}" -v

