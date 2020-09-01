package runner

import (
	"fmt"
	"github.com/manuelbua/go-version"
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

Probe targets at the specified port for
HTTPS/HTTP support and print an URL if
a connection can be established.

Input format is "host,port" or "ip,host,port".

`

func ShowBanner() {
	fmt.Fprintf(os.Stderr, banner, version.GetVersion()+"/@dudez")
}
