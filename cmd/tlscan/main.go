package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"tlscan/pkg/tlscan"
)

func main() {
	o := tlscan.ParseOptions()

	if !tlscan.HasStdin() {
		println("No data at stdin (host,port)")
		return
	}

	limiter := make(chan struct{}, o.Threads)
	outputMutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		in := scanner.Text()
		if len(in) == 0 {
			continue
		}

		s := strings.Split(in, ",")
		host, port := s[0], s[1]

		wg.Add(1)
		limiter <- struct{}{}
		go func() {
			defer wg.Done()
			hasTls, err := tlscan.TlsConnect(host, port, o.Timeout)
			if err == nil {
				proto := "http"
				if hasTls {
					proto = "https"
				}
				outputMutex.Lock()
				fmt.Printf("%s://%s:%s\n", proto, host, port)
				outputMutex.Unlock()
			}
			<-limiter
		}()
	}
	wg.Wait()
}
