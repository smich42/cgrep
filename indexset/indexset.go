package indexset

import (
	"fmt"
	"math"
	"math/bits"
)

// Types storable in an IndexSet must implement an indexing function.
// The indexing function must determine a unique integer index
// for each possible value of the stored type, starting at 0 and
// producing a contiguous interval over the integers.
type indexable interface {
	Idx() uint
}

type IndexSet[T indexable] struct {
	words []uint64 // A sequence of bits 0..n representing whether each element 0..n is in stored.
}

const wordSize = 64

// Constructs a new IndexSet with a given minimum capacity.
// In reality the capacity will be an integer multiple of wordSize.
func New[T indexable](capacity int) *IndexSet[T] {
	if capacity <= 0 {
		panic(fmt.Errorf("invalid capacity %d: set must be able to hold at least 1 element", capacity))
	}

	wordsNeeded := float64(capacity) / float64(wordSize)
	wholeWordsNeeded := int(math.Ceil(wordsNeeded))

	is := IndexSet[T]{words: make([]uint64, wholeWordsNeeded)}

	return &is
}

// Adds a new element to the set.
func (is IndexSet[T]) Place(element indexable) {
	// Ensure element is in range.
	if !is.validate(element) {
		panic(is.outOfRangeError(element))
	}

	word := is.wordOf(element)
	toggledBit := uint64(1) << is.wordPositionOf(element)

	// To turn on the i^th bit of the given word,
	// use OR with a number whose i^th bit is the only one toggled.
	// E.g. word = word | 00000100 will clearly toggle the word's 3rd bit.
	*word |= toggledBit
}

// Removes a new element to the set.
func (is IndexSet[T]) Remove(element indexable) {
	if !is.validate(element) {
		panic(is.outOfRangeError(element))
	}

	word := is.wordOf(element)
	untoggledBit := ^(uint64(1) << is.wordPositionOf(element))

	// To turn off the i^th bit of the given word,
	// use AND with a number whose i^th bit is the only one untoggled.
	// E.g. word = word & 111011 will clearly turn off the word's 3rd bit.
	*word &= untoggledBit
}

// Checks whether the set contains a given element.
func (is IndexSet[T]) Has(element indexable) bool {
	if !is.validate(element) {
		panic(is.outOfRangeError(element))
	}

	word := is.wordOf(element)
	toggledBit := uint64(1) << is.wordPositionOf(element)

	// AND between the word and a number whose i^th bit is the only one toggled
	// does not produce 0 IFF the i^th bit of the word is activated.
	return ((*word) & toggledBit) != uint64(0)
}

// The number of elements currently stored in the set.
func (is IndexSet[T]) Count() int {
	count := 0

	for _, word := range is.words {
		count += bits.OnesCount64(word)
	}

	return count
}

// The number of elements the set can store.
func (is IndexSet[T]) Capacity() int {
	return wordSize * len(is.words)
}

// Ensures the given element can be stored in the set.
func (is IndexSet[T]) validate(element indexable) bool {
	return element.Idx() < uint(is.Capacity())
}

// Produces an error indicating the given element is out of range.
func (is IndexSet[T]) outOfRangeError(element indexable) error {
	return fmt.Errorf("index %d out of range: capacity = %d", element.Idx(), is.Capacity())
}

// Retrieves the specific word in the bit sequence which contains
// the index of a given element.
func (is IndexSet[T]) wordOf(element indexable) *uint64 {
	return &is.words[element.Idx()/wordSize]
}

// Returns the position of the index of a the given element relative to the
// start of the word it belongs to.
func (is IndexSet[T]) wordPositionOf(element indexable) uint {
	return element.Idx() % wordSize
}
