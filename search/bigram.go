package search

import (
	"fmt"
	"regexp"
	"strings"
)

type Bigram struct {
	a, b rune
}

// Computes the index of a bigram.
// Each bigram is assigned a unique index starting from 0 in the following order.
// '  ', ' a', ' b', ..., 'a ', 'aa','ab', ..., 'zz'
func (bg Bigram) Idx() uint {
	return uint(magn(bg.a)*('z'-'a') + magn(bg.b))
}

// Creates a new bigram from the given two-letter string.
func New(ab string) *Bigram {
	bg := Bigram{}
	// Ensure string is a bigram.
	if len(ab) != 2 {
		panic(fmt.Errorf("invalid bigram length [%d], required: 2", len(ab)))
	}

	a, b := rune(ab[0]), rune(ab[1])
	// Ensure bigram contains only lowercase characters and spaces.
	if validate(a) && validate(b) {
		bg.a, bg.b = a, b
	} else {
		panic(fmt.Errorf("invalid character in bigram '%s'", ab))
	}

	return &bg
}

// Sanitises the given string so that only bigram-compliant characters remain.
func Clean(text string) string {
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
