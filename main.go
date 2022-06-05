package main

import (
	"fmt"

	"github.com/smich42/cgrep/indexset"
)

func main() {
	a := indexset.New[natural](1000)
	b := indexset.New[natural](1000)

	a.PlaceAll(natural{0}, natural{1}, natural{2}, natural{4}, natural{8}, natural{501}, natural{500}, natural{999})
	b.PlaceAll(natural{1}, natural{3}, natural{501}, natural{500}, natural{998})

	a.RemoveAll(natural{500}, natural{998}, natural{1000})
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
