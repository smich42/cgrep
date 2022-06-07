package main

import (
	"fmt"

	"github.com/smich42/cgrep/search"
)

func main() {
	for _, match := range *search.Match("hello, aworld!", "hello world this is em speaking", 0.75) {
		fmt.Printf("'%s'\n", match)
	}
}
