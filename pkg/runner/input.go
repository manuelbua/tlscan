package runner

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Input struct {
	Data  string
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
	sb := strings.Builder{}
	i.Count = 0
	for scanner.Scan() {
		line := scanner.Text()

		// skip empty lines
		if len(line) == 0 {
			continue
		}

		var host, _ = ParseLine(line)

		// deduplication
		if _, ok := usedInput[host]; !ok {
			usedInput[host] = true
			i.Count++
			sb.WriteString(line)
			sb.WriteString("\n")
		} else {
			dupeCount++
		}
	}

	if input != nil && input.Close() != nil {
		log.Fatalf("Couldn't close input file %s", o.TargetList)
	}

	i.Data = sb.String()
	if dupeCount > 0 {
		log.Printf("Supplied input was automatically deduplicated (%d removed).", dupeCount)
	}

	return i
}
