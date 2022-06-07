package search

import (
	"fmt"
	"regexp"
	"strings"
)

type bigram struct {
	a, b rune
}

// Computes the index of a bigram.
// Each bigram is assigned a unique index starting from 0 in the following order.
// '  ', ' a', ' b', ..., 'a ', 'aa','ab', ..., 'zz'
func (bg bigram) Idx() uint {
	return uint(magn(bg.a)*('z'-'a') + magn(bg.b))
}

func FirstBigram() bigram {
	return *newBigram(' ', ' ')
}

func LastBigram() bigram {
	return *newBigram('z', 'z')
}

// Creates a new bigram from the given two-letter string.
func newBigram(a, b rune) *bigram {
	bg := bigram{}

	// Ensure bigram contains only lowercase characters and spaces.
	if validate(a) && validate(b) {
		bg.a, bg.b = a, b
	} else {
		panic(fmt.Errorf("invalid character in bigram '%c%c'", a, b))
	}

	return &bg
}

// Sanitises the given string so that only bigram-compliant characters remain.
func clean(text string) string {
	sanitised := strings.Builder{}

	for _, c := range strings.ToLower(text) {
		if validate(c) {
			sanitised.WriteRune(c)
		} else {
			// Treat character as punctuation splitting words
			sanitised.WriteRune(' ')
		}
	}

	text = sanitised.String()
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	return text
}

// Checks whether character can be part of a bigram,
// i.e. whether it's lowercase english or a space.
func validate(c rune) bool {
	return magn(c) >= 0 && magn(c) <= magn('z')
}

// Returns the magnitude of a rune, i.e. its position in the sequence
// ' ', 'a', 'b', ..., 'z'
func magn(c rune) int {
	if c == ' ' {
		return 0
	}
	return int(c - 'a' + 1)
}
