package cmdutils

// Version contains the version of esdm. It is
// injected by the linker at compile time and is
// equal to the value of the Git tag, if the current
// commit is tagged. Otherwise, it defaults to the
// string "(version unavailable)".
var Version = "(version unavailable)"

// GitVersion is similar to Version, except that it
// doesn't refer to a version number, but to the Git
// commit hash of the current commit. This value, too,
// is injected by the linker at compile time. If this
// fails for any reason, GitVersion defaults to the
// string "(version unavailable)".
var GitVersion = "(version unavailable)"
