package tlscan

import (
	"fmt"
	"os"
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

`

const unversioned = "(unversioned)"

var version = unversioned

func ShowBanner() {
	fmt.Fprintf(os.Stderr, banner, version)
}
