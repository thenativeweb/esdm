package schema

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"
)

// CompareRevisions returns -1, 0, or +1 if a is older,
// equal, or newer than b. Both inputs must be SemVer
// strings of the form MAJOR.MINOR.PATCH; pre-release and
// build metadata are deliberately not supported because
// ESDM revisions do not use them.
func CompareRevisions(a, b string) (int, error) {
	aMajor, aMinor, aPatch, err := splitRevision(a)
	if err != nil {
		return 0, err
	}
	bMajor, bMinor, bPatch, err := splitRevision(b)
	if err != nil {
		return 0, err
	}

	switch {
	case aMajor != bMajor:
		return cmp.Compare(aMajor, bMajor), nil
	case aMinor != bMinor:
		return cmp.Compare(aMinor, bMinor), nil
	case aPatch != bPatch:
		return cmp.Compare(aPatch, bPatch), nil
	}
	return 0, nil
}

func splitRevision(revision string) (int, int, int, error) {
	parts := strings.Split(revision, ".")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("revision %q is not a SemVer (MAJOR.MINOR.PATCH)", revision)
	}

	out := [3]int{}
	for i, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 {
			return 0, 0, 0, fmt.Errorf("revision %q has non-numeric component %q", revision, part)
		}
		out[i] = n
	}
	return out[0], out[1], out[2], nil
}
