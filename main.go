package main

import (
	"fmt"

	"github.com/smich42/cgrep/indexset"
)

func main() {
	fmt.Println("Hello, World!")

	is := indexset.New[natural](1000)

	fmt.Println(is.Count())

	is.Place(natural{0})
	is.Place(natural{1})
	is.Place(natural{2})
	is.Place(natural{4})
	is.Place(natural{8})
	is.Place(natural{501})
	is.Place(natural{500})
	is.Place(natural{999})

	fmt.Println(is.Count())

	is.Remove(natural{500})
	is.Remove(natural{1000})

	fmt.Println(is.Count())

	is.Remove(natural{998})
	is.Place(natural{998})

	for _, n := range []int{0, 1, 2, 3, 500, 501, 998, 999, 1000} {
		fmt.Println(n, is.Has(natural{n}))
	}
}

type natural struct {
	Val int
}

func (i natural) Idx() uint {
	return uint(i.Val)
}
