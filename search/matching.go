package search

import (
	"strings"

	"github.com/smich42/cgrep/indexset"
)

// Produces all matches in the haystack whose similarity to the needle
// is over the given threshold.
func Match(needle, haystack string, threshold float32) *[]string {
	haystack = clean(haystack)
	needle = clean(needle)

	words := strings.Split(haystack, " ")
	wordCount := len(strings.Split(needle, " "))

	matches := []string{}

	for i := 0; i < len(words)-wordCount; i++ {
		candidate := words[i]
		for j := 1; j < wordCount; j++ {
			candidate += " " + words[i+j]
		}

		if similarity(needle, candidate) >= threshold {
			matches = append(matches, candidate)
		}
	}

	return &matches
}

func similarity(a, b string) float32 {
	return indexset.Jaccard(bigramsOf(a), bigramsOf(b))
}

func bigramsOf(s string) *indexset.IndexSet[bigram] {
	possibleBigrams := int(LastBigram().Idx())
	bigrams := indexset.NewIndexSet[bigram](possibleBigrams)

	for i := 0; i < len(s)-1; i++ {
		a := rune(s[i])
		b := rune(s[i+1])

		bigrams.Place(newBigram(a, b))
	}

	return bigrams
}
