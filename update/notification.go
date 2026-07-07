package update

import "strings"

const (
	// UpgradeURL is where the notification points users
	// for upgrade instructions. It targets the dedicated
	// "Updating ESDM" page in the documentation site.
	UpgradeURL = "https://www.esdm.io/getting-started/updating-esdm/"

	// DefaultEndpoint is the URL the version check fetches
	// to learn about the latest released version.
	DefaultEndpoint = "https://www.esdm.io/version.json"

	// DisableEnvVar names the environment variable that
	// turns the version check off when set to "true".
	DisableEnvVar = "ESDM_DISABLE_UPDATE_CHECK"

	ansiBold  = "\x1b[1m"
	ansiDim   = "\x1b[2m"
	ansiReset = "\x1b[0m"
)

// RenderNotification builds the multi-line update hint
// that gets written to stderr after a normal command's
// output. When shouldIncludeDisableHint is true, an
// additional dimmed line explains how to turn the check
// off; this is used the first time the hint is shown.
func RenderNotification(currentVersion, latestVersion string, shouldIncludeDisableHint, isColorEnabled bool) string {
	current := StripVPrefix(currentVersion)
	latest := StripVPrefix(latestVersion)

	var builder strings.Builder

	if isColorEnabled {
		builder.WriteString(ansiBold)
	}
	builder.WriteString("⚡ A new version of esdm is available: ")
	builder.WriteString(current)
	builder.WriteString(" → ")
	builder.WriteString(latest)
	if isColorEnabled {
		builder.WriteString(ansiReset)
	}

	builder.WriteString("\n  See ")
	builder.WriteString(UpgradeURL)
	builder.WriteString(" for upgrade instructions.")

	if shouldIncludeDisableHint {
		builder.WriteString("\n  ")
		if isColorEnabled {
			builder.WriteString(ansiDim)
		}
		builder.WriteString("esdm checks for updates once a day. Set ")
		builder.WriteString(DisableEnvVar)
		builder.WriteString("=true to disable.")
		if isColorEnabled {
			builder.WriteString(ansiReset)
		}
	}

	return builder.String()
}
