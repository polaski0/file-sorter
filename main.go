package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	src, dest := parseArguments()
	fmt.Printf("Sorting %v to %v...\n", strings.Join(src, ", "), dest)

	c := NewCounter()

	sr := NewSorter(c, src, dest)
	err := sr.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func parseArguments() (source []string, destination string) {
	args := os.Args[1:]
	source = []string{}
	destination = ""

	if len(args) < 1 {
		log.Fatal("Source directory to be sorted required")
	}

	for i, s := range args {
		if len(args) == 1 {
			source = append(source, path.Clean(s))
			break
		}

		if i != len(args)-1 {
			source = append(source, path.Clean(s))
		} else {
			destination = s
		}
	}

	if len(source) == 0 {
		log.Fatal("Source directory to be filtered required")
	}

	if destination == "" {
		destination = "./tmp" // Create a temporary directory on the current file
	}

	return source, path.Clean(destination)
}
