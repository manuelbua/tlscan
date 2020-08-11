package runner

import (
	"fmt"
	"os"
	"strings"
)

const banner = `
  __  .__
_/  |_|  |   ______ ____ _____    ____
\   __\  |  /  ___// ___\\__  \  /    \
 |  | |  |__\___ \\  \___ / __ \|   |  \
 |__| |____/____  >\___  >____  /___|  /
                \/     \/     \/     \/
%40s

Probes HTTP servers for TLS support.
Input is <host,port> or <ip,host,port>.

`

const unversioned = "(unversioned)"

var version = unversioned

func ShowBanner() {
	fmt.Fprintf(os.Stderr, banner, GetVersion())
}

func GetVersion() string {
	if len(strings.TrimSpace(version)) == 0 {
		return unversioned
	}
	return version
}