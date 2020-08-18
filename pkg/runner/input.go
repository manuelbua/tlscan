package runner

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type Input struct {
	Data  [][3]string
	Count int64
}

func NewInput(o *Options) *Input {
	i := &Input{}

	var scanner *bufio.Scanner
	var input *os.File
	var err error

	if o.HasStdin {
		scanner = bufio.NewScanner(os.Stdin)
	} else if o.HasTargetString {
		scanner = bufio.NewScanner(strings.NewReader(o.Targets))
	} else {
		input, err = os.Open(o.TargetList)
		if err != nil {
			log.Fatalf("Could not open target file %s", o.TargetList)
		}
		scanner = bufio.NewScanner(input)
	}

	// Sanitize input, deduplicate and precompute total number of targets
	var usedInput = make(map[string]bool)
	dupeCount := 0
	i.Count = 0
	for scanner.Scan() {
		line := scanner.Text()

		// skip empty lines
		if len(line) == 0 {
			continue
		}

		var ip, host, port, ipHostPort, err = parseLine(line)
		if err != nil {
			log.Println(err)
		} else {
			// deduplication
			if _, ok := usedInput[ipHostPort]; !ok {
				usedInput[ipHostPort] = true
				i.Count++
				i.Data = append(i.Data, [3]string{ip, host, port})
			} else {
				dupeCount++
			}
		}
	}

	if input != nil && input.Close() != nil {
		log.Fatalf("Couldn't close input file %s", o.TargetList)
	}

	if dupeCount > 0 {
		log.Printf("Supplied input was automatically deduplicated (%d removed).", dupeCount)
	}

	return i
}

func parseLine(line string) (ip, host, port, hostPort string, err error) {
	s := strings.Split(line, ",")
	ip = ""

	switch len(s) {
	case 2:
		host, port = s[0], s[1]
	case 3:
		ip, host, port = s[0], s[1], s[2]
	default:
		return "", "", "", "", errors.New(fmt.Sprintf("Unsupported input format: %s", line))
	}

	hostPort = fmt.Sprintf("%s:%s:%s", ip, host, port)
	return ip, host, port, hostPort, nil
}
