package runner

import (
	"flag"
	"fmt"
	"github.com/manuelbua/go-version"
	"log"
	"os"
)

type Options struct {
	Timeout         float64 // Timeout is the number of seconds to wait for the handshake to complete
	Threads         int     // Threads is the number of concurrent connections to make
	Target          string  // Target represents a single target
	TargetList      string  // TargetList is the file with a list of targets
	UserAgent       string  // UserAgent indicates the user agent to use
	OnlyTls         bool    // OnlyTls indicates to only produce output for TLS-enabled servers
	OnlyPlain       bool    // OnlyPlain indicates to only produce output for non-TLS-enabled servers
	HasStdin        bool    // HasStdin indicates if input is present at stdin
	HasTargetString bool    // HasTargetString indicates if Target is valid
	HasTargetList   bool    // HasTargetList indicates if TargetList is valid
	NoColor         bool    // NoColor indicates to not colorize output
	NoProgressBar   bool    // NoProgressBar indicates to not use a progressbar
	ShowVersion     bool    // Shows version and exit
}

func ParseOptions() *Options {
	options := &Options{}

	flag.StringVar(&options.Target, "t", "", "Specify a single target (stdin has precedence)")
	flag.StringVar(&options.TargetList, "tl", "", "Specify a file with a list of targets, one per line (stdin has precedence)")
	flag.StringVar(&options.UserAgent, "ua", "tlscan (https://github.com/manuelbua/tlscan)", "Specify a custom User-Agent")
	flag.Float64Var(&options.Timeout, "timeout", 10, "Seconds to wait for the handshake to complete")
	flag.IntVar(&options.Threads, "c", 20, "Number of concurrent connections to make")
	flag.BoolVar(&options.OnlyTls, "https", false, "Output only TLS-enabled servers")
	flag.BoolVar(&options.OnlyPlain, "http", false, "Output only non-TLS-enabled servers")
	flag.BoolVar(&options.NoColor, "nc", false, "Do not colorize output")
	flag.BoolVar(&options.NoProgressBar, "nobar", false, "Do not use a progressbar")
	flag.BoolVar(&options.ShowVersion, "v", false, "Shows version and exit")

	flag.Parse()

	if options.ShowVersion {
		fmt.Println(version.GetVersionLong())
		os.Exit(0)
	}

	ShowBanner()

	if hasStdin() {
		options.HasStdin = true
	} else if len(options.Target) > 0 {
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
