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
	Count     int64
	DupeCount int

	src                        *os.File
	scanner                    *bufio.Scanner
	scanned                    string
	ip, host, port, ipHostPort string
	usedInput                  map[string]bool
}

func NewInput(o *Options) *Input {
	scanner, input := newScanner(o)

	i := &Input{
		scanner:   scanner,
		src:       input,
		Count:     0,
		DupeCount: 0,
		usedInput: make(map[string]bool),
	}

	if !o.HasStdin {
		// precompute totals
		for i.Scan() {
			i.Count++
		}

		if i.DupeCount > 0 {
			log.Printf("Supplied input was automatically deduplicated (%d removed).", i.DupeCount)
		}

		if o.HasTargetList {
			_, err := input.Seek(0, 0)
			if err != nil {
				log.Fatalf("Couldn't seek input")
			}
		}
		i.scanner, i.src = newScanner(o)
		i.usedInput = make(map[string]bool)
	} else {
		i.Count, i.DupeCount = -1, -1
	}

	return i
}

func newScanner(o *Options) (scanner *bufio.Scanner, input *os.File) {
	var err error
	input = nil

	if o.HasStdin {
		scanner = bufio.NewScanner(os.Stdin)
	} else if o.HasTargetString {
		scanner = bufio.NewScanner(strings.NewReader(o.Target))
	} else {
		input, err = os.Open(o.TargetList)
		if err != nil {
			log.Fatalf("Could not open target file %s", o.TargetList)
		}
		scanner = bufio.NewScanner(input)
	}
	return scanner, input
}

func (i *Input) Begin() bool {
	return i.Count != 0
}

func (i *Input) Scan() bool {
	var err error

	for {
		scanned := i.scanner.Scan()
		if !scanned {
			return false
		}

		line := i.scanner.Text()
		if len(line) == 0 {
			continue
		}

		i.ip, i.host, i.port, i.ipHostPort, err = parseLine(line)
		if err != nil {
			log.Println(err)
			continue
		}

		if _, dupe := i.usedInput[i.ipHostPort]; dupe {
			i.DupeCount++
			continue
		}

		i.usedInput[i.ipHostPort] = true
		return true
	}
}

func (i *Input) Data() (ip, host, port, ipHostPort string) {
	return i.ip, i.host, i.port, i.ipHostPort
}

func (i *Input) End() {
	if i.src != nil {
		err := i.src.Close()
		if err != nil {
			log.Fatal("Couldn't close input file")
		}
	}
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
