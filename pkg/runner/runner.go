package runner

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/manuelbua/tlscan/pkg/progress"
	"github.com/manuelbua/tlscan/pkg/scanner/http"
	"log"
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
	runner.progress = progress.NewProgress(opts.NoColor, !opts.NoProgressBar && !opts.HasStdin)

	return runner, nil
}

func (r *Runner) Run() {
	input := r.input
	opts := r.options

	if input.Begin() {
		if !opts.HasStdin {
			log.Printf("Processing %s hosts.", r.colorizer.Bold(input.Count).String())
		}

		uniqueOutput := make(map[string]bool)
		limiter := make(chan struct{}, opts.Threads)
		outputMutex := sync.Mutex{}
		wg := sync.WaitGroup{}

		httpScanner := http.NewScanner(opts.Timeout, opts.UserAgent)

		p := r.progress
		p.InitProgressbar(input.Count)

		for input.Scan() {
			var ip, host, port, _ = input.Data()

			wg.Add(1)
			limiter <- struct{}{}
			go func() {
				defer wg.Done()
				hasTls, err := httpScanner.Scan(ip, host, port)
				if err == nil {
					if (!opts.OnlyPlain && !opts.OnlyTls) ||
						(opts.OnlyTls && hasTls) ||
						(opts.OnlyPlain && !hasTls) {
						proto := "http"
						if hasTls {
							proto = "https"
						}

						outputMutex.Lock()
						out := fmt.Sprintf("%s://%s:%s\n", proto, host, port)
						if _, ok := uniqueOutput[out]; !ok {
							uniqueOutput[out] = true
							fmt.Print(out)
						}
						outputMutex.Unlock()
					}
				}
				p.Update()
				<-limiter
			}()
		}
		wg.Wait()
		p.Wait()

		input.End()
	} else {
		log.Println("No input found.")
	}
}
