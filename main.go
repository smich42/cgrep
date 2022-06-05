package main

import (
	"fmt"

	"github.com/smich42/cgrep/indexset"
)

func main() {
	a := indexset.New[natural](1000)
	b := indexset.New[natural](1000)

	a.Place(natural{0})
	a.Place(natural{1})
	a.Place(natural{2})
	a.Place(natural{4})
	a.Place(natural{8})
	a.Place(natural{501})
	a.Place(natural{500})
	a.Place(natural{999})

	b.Place(natural{0})
	b.Place(natural{1})
	b.Place(natural{501})
	b.Place(natural{500})
	b.Place(natural{998})

	a.Remove(natural{500})
	a.Remove(natural{1000})

	a.Remove(natural{998})
	a.Place(natural{998})

	union := indexset.Union(a, b)
	intersection := indexset.Intersection(a, b)

	fmt.Println(a.String())
	fmt.Println(b.String())
	fmt.Println(union.String())
	fmt.Println(intersection.String())
}

type natural struct {
	Val int
}

func (i natural) Idx() uint {
	return uint(i.Val)
}
