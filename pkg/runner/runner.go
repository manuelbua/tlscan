package runner

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/manuelbua/tlscan/pkg/progress"
	"github.com/manuelbua/tlscan/pkg/scanner/http"
	"log"
	"strings"
	"sync"
)

type Runner struct {
	options *Options
	input   *Input
	// progress tracking
	progress progress.IProgress

	// output coloring
	colorizer aurora.Aurora
	//decolorizer *regexp.Regexp
}

func New() (*Runner, error) {
	opts := ParseOptions()
	runner := &Runner{
		options: opts,
		input:   NewInput(opts),
	}

	// output coloring
	useColor := !opts.NoColor
	runner.colorizer = aurora.NewAurora(useColor)
	//if useColor {
	//	// compile a decolorization regex to cleanup file output messages
	//	compiled, err := regexp.Compile("\\x1B\\[[0-9;]*[a-zA-Z]")
	//	if err != nil {
	//		return nil, err
	//	}
	//	runner.decolorizer = compiled
	//}

	// progress tracking
	runner.progress = progress.NewProgress(opts.NoColor, !opts.NoProgressBar)

	return runner, nil
}

func (r *Runner) Run() {
	input := r.input
	opts := r.options

	log.Printf("Processing %s hosts.", r.colorizer.Bold(input.Count).String())

	scanner := bufio.NewScanner(strings.NewReader(input.Data))

	limiter := make(chan struct{}, opts.Threads)
	outputMutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	httpScanner := http.New(opts.Timeout)
	p := r.progress
	p.InitProgressbar(input.Count)

	for scanner.Scan() {
		in := scanner.Text()
		var host, port, _ = ParseLine(in)

		wg.Add(1)
		limiter <- struct{}{}
		go func() {
			defer wg.Done()
			hasTls, err := httpScanner.Scan(host, port)
			if err == nil {
				if (!opts.OnlyPlain && !opts.OnlyTls) ||
					(opts.OnlyTls && hasTls) ||
					(opts.OnlyPlain && !hasTls) {
					proto := "http"
					if hasTls {
						proto = "https"
					}
					outputMutex.Lock()
					fmt.Printf("%s://%s:%s\n", proto, host, port)
					outputMutex.Unlock()
				}
			}
			p.Update()
			<-limiter
		}()
	}
	wg.Wait()
	p.Wait()
}

func ParseLine(line string) (host, port, hostport string) {
	s := strings.Split(line, ",")

	switch len(s) {
	case 2:
		host, port = s[0], s[1]
	case 3:
		host, port = s[1], s[2]
	default:
		log.Printf("Unsupported input format: %s", line)
	}
	return host, port, fmt.Sprintf("%s:%s", host, port)
}
