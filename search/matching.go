package search

import (
	"strings"

	"github.com/smich42/cgrep/indexset"
)

// Produces all matches in the haystack whose similarity to the needle
// is over the given threshold.
func match(needle, haystack string, threshold float32) *[]string {
	haystack = clean(haystack)
	needle = clean(needle)

	words := strings.Split(haystack, " ")
	// All matches must have the same word count as the search pattern.
	wordCount := len(strings.Split(needle, " "))

	matches := []string{}

	for i := 0; i <= len(words)-wordCount; i++ {
		// Produce the next sentence with the same word count as the search pattern.
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

// Ranks the similarity of two strings from 0 to 1.
func similarity(a, b string) float32 {
	return indexset.Jaccard(bigramsOf(a), bigramsOf(b))
}

// Produces an IndexSet containing all bigrams in the given string.
func bigramsOf(s string) *indexset.IndexSet[bigram] {
	// Bigrams range from '  ' (index 0) to 'zz' (index 26*26).
	// The IndexSet must accommodate them all.
	possibleBigrams := int(lastBigram().Idx())
	bigrams := indexset.NewIndexSet[bigram](possibleBigrams)

	for i := 0; i < len(s)-1; i++ {
		a := rune(s[i])
		b := rune(s[i+1])

		bigrams.Place(newBigram(a, b))
	}

	return bigrams
}
