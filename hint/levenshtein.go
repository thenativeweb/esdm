package hint

// Levenshtein returns the minimum edit distance between a
// and b - the number of single-character insertions,
// deletions, or substitutions needed to turn one into the
// other. The implementation operates on Unicode code
// points so multi-byte characters in identifiers count
// correctly.
func Levenshtein(a, b string) int {
	aRunes := []rune(a)
	bRunes := []rune(b)

	if len(aRunes) == 0 {
		return len(bRunes)
	}
	if len(bRunes) == 0 {
		return len(aRunes)
	}

	previousRow := make([]int, len(bRunes)+1)
	currentRow := make([]int, len(bRunes)+1)

	for j := range previousRow {
		previousRow[j] = j
	}

	for i := 1; i <= len(aRunes); i++ {
		currentRow[0] = i
		for j := 1; j <= len(bRunes); j++ {
			cost := 1
			if aRunes[i-1] == bRunes[j-1] {
				cost = 0
			}
			currentRow[j] = min(
				previousRow[j]+1,      // deletion
				currentRow[j-1]+1,     // insertion
				previousRow[j-1]+cost, // substitution
			)
		}
		previousRow, currentRow = currentRow, previousRow
	}

	return previousRow[len(bRunes)]
}

// ClosestMatch returns the candidate whose Levenshtein
// distance to name is smallest, plus that distance. If
// candidates is empty, ok is false.
//
// Callers decide how strict the match has to be by
// comparing the returned distance against a threshold.
func ClosestMatch(name string, candidates []string) (best string, distance int, ok bool) {
	if len(candidates) == 0 {
		return "", 0, false
	}

	best = candidates[0]
	distance = Levenshtein(name, best)
	for _, candidate := range candidates[1:] {
		candidateDistance := Levenshtein(name, candidate)
		if candidateDistance < distance {
			best = candidate
			distance = candidateDistance
		}
	}

	return best, distance, true
}

// SuggestionThreshold reports the maximum edit distance
// at which a candidate is still considered a plausible
// "did you mean?" suggestion for the given name. Short
// names tolerate up to two edits, longer names allow up
// to a third of their length.
func SuggestionThreshold(name string) int {
	threshold := len([]rune(name)) / 3
	if threshold < 2 {
		return 2
	}
	return threshold
}

// Best returns the closest candidate to name if it is
// within SuggestionThreshold(name), is not name itself,
// and candidates is non-empty. Otherwise it returns "",
// false. This is the one-shot convenience wrapper used
// when callers only need the candidate name and do not
// care about the distance.
func Best(name string, candidates []string) (string, bool) {
	best, distance, ok := ClosestMatch(name, candidates)
	if !ok {
		return "", false
	}
	if distance > SuggestionThreshold(name) {
		return "", false
	}
	if best == name {
		return "", false
	}
	return best, true
}
