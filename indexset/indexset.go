package indexset

import (
	"fmt"
	"math"
	"math/bits"
	"strconv"
	"strings"
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
func NewIndexSet[T indexable](capacity int) *IndexSet[T] {
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

// Adds all given elements to the set.
func (is IndexSet[T]) PlaceAll(elements ...indexable) {
	for _, element := range elements {
		is.Place(element)
	}
}

// Removes a new element from the set.
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

// Removes all given elements from the set.
func (is IndexSet[T]) RemoveAll(elements ...indexable) {
	for _, element := range elements {
		is.Remove(element)
	}
}

// Checks whether the set contains a given element.
func (is IndexSet[T]) Has(element indexable) bool {
	if !is.validate(element) {
		panic(is.outOfRangeError(element))
	}

	word := is.wordOf(element)
	return toggled(is.wordPositionOf(element), *word)
}

// Checks whether the set contains all given elements.
func (is IndexSet[T]) HasAll(elements ...indexable) bool {
	for _, element := range elements {
		if !is.Has(element) {
			return false
		}
	}
	return true
}

// Checks whether the set contains any of the given elements.
func (is IndexSet[T]) HasAny(elements ...indexable) bool {
	for _, element := range elements {
		if is.Has(element) {
			return true
		}
	}
	return false
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

// Builts a pretty string of for the indices in the set.
func (is IndexSet[T]) String() string {
	indices := []string{}

	for i := 0; i < len(is.words); i++ {
		for j := 0; j < wordSize; j++ {
			if toggled(uint(j), is.words[i]) {
				indices = append(indices, strconv.Itoa(i*wordSize+j))
			}
		}
	}

	return "{" + strings.Join(indices, ", ") + "}"
}

// Returns a new IndexSet which is the union of two given sets.
func Union[T indexable](a, b *IndexSet[T]) *IndexSet[T] {
	// The capacity of the new IndexSet will be that of the largest of the two sets.
	jointCapacity := int(math.Max(
		float64(a.Capacity()),
		float64(b.Capacity())))

	c := NewIndexSet[T](jointCapacity)

	for i := 0; i < len(c.words); i++ {
		var wordA, wordB uint64 = 0, 0

		if i < len(a.words) {
			wordA = a.words[i]
		}
		if i < len(b.words) {
			wordB = b.words[i]
		}

		// Any bits that are 1 in either constituent set will be 1 in the result.
		c.words[i] = wordA | wordB
	}

	return c
}

// Returns a new IndexSet which is the intersection of two given sets.
func Intersection[T indexable](a, b *IndexSet[T]) *IndexSet[T] {
	// The capacity of the new IndexSet will be that of the smallest of the two sets.
	jointCapacity := int(math.Min(
		float64(a.Capacity()),
		float64(b.Capacity())))

	c := NewIndexSet[T](jointCapacity)

	for i := 0; i < len(c.words); i++ {
		// Only bits that are 1 in both constituent sets will be 1 in the result.
		c.words[i] = a.words[i] & b.words[i]
	}

	return c
}

// Computes the Jaccard similarity of two sets.
func Jaccard[T indexable](a, b *IndexSet[T]) float32 {
	intersection := Intersection(a, b).Count()
	union := Union(a, b).Count()

	if union == 0 {
		return 1.0 // Empty sets are identical
	}

	return float32(intersection) / float32(union)
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

// Checks whether a given bit in a word is 1.
func toggled(bitpos uint, word uint64) bool {
	toggledBit := uint64(1) << bitpos
	// AND between the word and a number whose i^th bit is the only one toggled
	// does not produce 0 IFF the i^th bit of the word is activated.
	return (word & toggledBit) != uint64(0)
}
