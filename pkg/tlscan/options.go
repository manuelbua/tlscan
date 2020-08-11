package tlscan

import (
	"flag"
	"log"
	"os"
)

type Options struct {
	Timeout         float64 // Timeout is the number of seconds to wait for the handshake to complete
	Threads         int     // Threads is the number of concurrent connections to make
	Targets         string  // Target is a single target
	TargetList      string  // TargetList is the file with a list of targets
	OnlyTls 		bool    // OnlyTls indicates to only produce output for TLS-enabled servers
	OnlyPlain  		bool  	// OnlyPlain indicates to only produce output for non-TLS-enabled servers
	HasStdin        bool    // HasStdin indicates if input is present at stdin
	HasTargetString bool    // HasSingleTarget indicates if Target is valid
	HasTargetList   bool    // HasTargetList indicates if TargetList is valid
}

func ParseOptions() *Options {
	options := &Options{}

	flag.StringVar(&options.Targets, "t", "", "Specify a single containing one or more targets, newline separated")
	flag.StringVar(&options.TargetList, "tL", "", "Specify a file with a list of targets, one per line")
	flag.Float64Var(&options.Timeout, "timeout", 5, "Seconds to wait for the handshake to complete")
	flag.IntVar(&options.Threads, "c", 20, "Number of concurrent connection to make")
	flag.BoolVar(&options.OnlyTls, "https", false, "Output only TLS-enabled servers")
	flag.BoolVar(&options.OnlyPlain, "plain", false, "Output only non-TLS-enabled servers")

	flag.Parse()

	if hasStdin() {
		options.HasStdin = true
	} else if len(options.Targets) > 0 {
		options.HasTargetString = true
	} else {
		if len(options.TargetList) > 0 {
			isFile, err := isFilePath(options.TargetList)
			if err != nil {
				log.Fatalf("Error opening target list: %s", err)
			} else {
				options.HasTargetList = isFile
			}
		}
	}

	if !options.HasStdin && !options.HasTargetString && !options.HasTargetList {
		log.Fatal("Please supply some input.")
	}

	return options
}

func hasStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return false
	}
	return true
}

func isFilePath(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Mode().IsRegular(), nil
}
