package tlscan

import (
	"flag"
	"os"
)

type Options struct {
	Timeout float64 // Timeout is the number of seconds to wait for the handshake to complete
	Threads int     // Threads is the number of concurrent connections to make
}

func ParseOptions() *Options {
	options := &Options{}

	flag.Float64Var(&options.Timeout, "timeout", 5, "Seconds to wait for the handshake to complete")
	flag.IntVar(&options.Threads, "c", 20, "Number of concurrent connection to make")

	flag.Parse()
	return options
}

func HasStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return false
	}
	return true
}
