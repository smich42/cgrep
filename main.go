package main

import (
	"fmt"

	"github.com/smich42/cgrep/search"
)

func main() {
	results, _ := search.SearchDir("quick brown fox", "testdata/", 0.5)
	for filepath, matches := range *results {
		fmt.Println(filepath)
		for _, match := range *matches {
			fmt.Println("\t" + match)
		}
	}
}
