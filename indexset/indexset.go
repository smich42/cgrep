package indexset

import (
	"fmt"
	"math"
)

type indexable interface {
	Idx() uint
}

type IndexSet[T indexable] struct {
	words []int64
}

const wordSize = 64

func New[T indexable](capacity int) *IndexSet[T] {

	if capacity <= 0 {
		panic(fmt.Errorf("invalid capacity %d: set must be able to hold at least 1 element", capacity))
	}

	wordsNeeded := float64(capacity) / float64(wordSize)
	wholeWordsNeeded := int(math.Ceil(wordsNeeded))

	is := IndexSet[T]{words: make([]int64, wholeWordsNeeded)}

	return &is
}

func (is IndexSet[T]) Place(element indexable) {

	if !is.validate(element) {
		panic(is.outOfRangeError(element))
	}

	word := is.wordOf(element)
	toggledBit := int64(1) << int64(is.inWordPositionOf(element))

	*word |= toggledBit
}

func (is IndexSet[T]) Remove(element indexable) {

	if !is.validate(element) {
		panic(is.outOfRangeError(element))
	}

	word := is.wordOf(element)
	untoggledBit := ^(int64(1) << int64(is.inWordPositionOf(element)))

	*word &= untoggledBit
}

func (is IndexSet[T]) Has(element indexable) bool {

	if !is.validate(element) {
		panic(is.outOfRangeError(element))
	}

	word := is.wordOf(element)
	toggledBit := int64(1) << int64(is.inWordPositionOf(element))

	return (*word)&toggledBit != int64(0)
}

func (is IndexSet[T]) Capacity() int {
	return wordSize * len(is.words)
}

func (is IndexSet[T]) validate(element indexable) bool {
	return element.Idx() < uint(is.Capacity())
}

func (is IndexSet[T]) outOfRangeError(element indexable) error {
	return fmt.Errorf("index %d out of range: capacity = %d", element.Idx(), is.Capacity())
}

func (is IndexSet[T]) wordOf(element indexable) *int64 {
	return &is.words[element.Idx()/wordSize]
}

func (is IndexSet[T]) inWordPositionOf(element indexable) uint {
	return element.Idx() % wordSize
}
