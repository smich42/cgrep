package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/smich42/cgrep/search"
)

const usageHelp = "Usage: cgrep '[search string]' [0-1.0 similarity] [directory]"

func main() {
	argcnt := len(os.Args) - 1

	// Require at least two arguments: pattern and similarity threshold.
	if argcnt <= 1 {
		fmt.Println(usageHelp)
		return
	}

	searchPattern := os.Args[1]
	searchThreshold, converr := strconv.ParseFloat(os.Args[2], 32)
	// Ensure threshold could be parsed.
	if converr != nil {
		fmt.Println("Second argument must be a floating-point number.")
		return
	}
	// Default to current working directory if not provided.
	searchPath := "."
	if argcnt >= 3 {
		searchPath = os.Args[3]
	}

	results, _ := search.SearchDir(searchPattern, searchPath, float32(searchThreshold))

	for filepath, matches := range *results {
		for _, match := range *matches {
			fmt.Printf("[%s] '%s'\n", filepath, match)
		}
	}
}
