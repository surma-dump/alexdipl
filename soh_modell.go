package main

import (
	"./stoichio"
	"github.com/surma/goptions"
	"log"
)

const VERSION = "0.1"

func main() {
	var options struct {
		InputFile     string `goptions:"-i, --input, description='File to read', obligatory"`
		goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}

	err := goptions.Parse(&options)
	if err != nil {
		log.Printf("Error: %s", err)
		goptions.PrintHelp()
		return
	}

	matrix, irreversible, err := stoichio.ReadFile(options.InputFile)
	if err != nil {
		log.Fatalf("Could not read file: %s", err)
	}

	_, _ = matrix, irreversible
}
