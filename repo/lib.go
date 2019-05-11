package repo

func findInRunes(rs []rune, c rune, start int) int {
	for i, r := range rs {
		if i < start {
			continue
		}
		if r == c {
			return i
		}
	}
	return -1
}
