package classification

// multilingualTextUnitCount counts content-bearing units without relying on
// whitespace tokenization. CJK runes count individually, while contiguous runs
// of non-CJK letters/digits count as one unit.
func multilingualTextUnitCount(text string) int {
	return 0
}

func keywordOccurrenceCount(text string, keywords []string, caseSensitive bool) int {
	return 0
}

func countKeywordOccurrences(text string, keyword string) int {
	return 0
}

func keywordBoundaryMatch(text string, start int, end int) bool {
	return false
}

func isBoundaryBlockingRune(r rune) bool {
	return false
}

func containsCJK(text string) bool {
	return false
}

func isCJK(r rune) bool {
	return false
}

func utf8DecodeRuneInString(s string) (rune, int) {
	return rune(0), 0
}

func utf8DecodeLastRuneInString(s string) (rune, int) {
	return rune(0), 0
}
