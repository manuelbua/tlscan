package main

import (
	"github.com/manuelbua/tlscan/pkg/runner"
	"log"
)

func main() {
	log.SetFlags(0)

	runner, err := runner.New()
	if err != nil {
		log.Fatal("Couldn't create runner")
	}

	runner.Run()
}
