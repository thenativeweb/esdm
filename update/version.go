package update

import (
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

// IsNewer reports whether latest is strictly greater than
// current using semver ordering. Both arguments may carry
// a leading "v" or omit it; they are normalized before
// comparison. An invalid input returns an error.
func IsNewer(current, latest string) (bool, error) {
	normalizedCurrent := normalizeVersion(current)
	if !semver.IsValid(normalizedCurrent) {
		return false, fmt.Errorf("invalid current version %q", current)
	}

	normalizedLatest := normalizeVersion(latest)
	if !semver.IsValid(normalizedLatest) {
		return false, fmt.Errorf("invalid latest version %q", latest)
	}

	return semver.Compare(normalizedCurrent, normalizedLatest) < 0, nil
}

// StripVPrefix returns the version without a leading "v",
// suitable for user-facing rendering. Inputs that do not
// start with "v" are returned unchanged.
func StripVPrefix(version string) string {
	return strings.TrimPrefix(version, "v")
}

func normalizeVersion(version string) string {
	if strings.HasPrefix(version, "v") {
		return version
	}
	return "v" + version
}
