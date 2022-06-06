package search

type bigram struct {
	a, b rune
}

func (bg bigram) Idx() uint {
	return magn(bg.a)*('z'-'a') + magn(bg.b)
}

func magn(c rune) uint {
	return uint(c - 'a')
}
