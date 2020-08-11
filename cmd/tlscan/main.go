package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/manuelbua/tlscan/pkg/tlscan"
)

func main() {
	tlscan.ShowBanner()
	log.SetFlags(0)

	o := tlscan.ParseOptions()

	var scanner *bufio.Scanner

	if o.HasStdin {
		scanner = bufio.NewScanner(os.Stdin)
	} else if o.HasTargetString {
		scanner = bufio.NewScanner(strings.NewReader(o.Targets))
	} else {
		input, err := os.Open(o.TargetList)
		if err != nil {
			log.Fatalf("Could not open target file %s", o.TargetList)
		}
		scanner = bufio.NewScanner(input)
		defer input.Close()
	}

	limiter := make(chan struct{}, o.Threads)
	outputMutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	for scanner.Scan() {
		in := scanner.Text()
		if len(in) == 0 {
			continue
		}

		s := strings.Split(in, ",")
		var host, port string

		switch len(s) {
		case 2:
			host, port = s[0], s[1]
		case 3:
			host, port = s[1], s[2]
		default:
			log.Printf("Unsupported input format: %s", in)
			continue
		}

		wg.Add(1)
		limiter <- struct{}{}
		go func() {
			defer wg.Done()
			hasTls, err := tlscan.TlsConnect(host, port, o.Timeout)
			if err == nil {
				if (!o.OnlyPlain && !o.OnlyTls) ||
					(o.OnlyTls && hasTls) ||
					(o.OnlyPlain && !hasTls) {
					proto := "http"
					if hasTls {
						proto = "https"
					}
					outputMutex.Lock()
					fmt.Printf("%s://%s:%s\n", proto, host, port)
					outputMutex.Unlock()
				}
			}
			<-limiter
		}()
	}
	wg.Wait()
}
