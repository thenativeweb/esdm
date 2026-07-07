// Package modelpath parses the slash-separated model-path strings that
// the command-line tools accept to scope their output to part of a
// model. A path is a sequence of segments drawn from the natural model
// hierarchy, from the domain down toward an individual command or
// event; the empty path selects the whole model.
//
// The parser is deliberately strict about shape rather than about
// meaning: it rejects a leading slash and empty interior segments so a
// path always identifies a single, unambiguous node, but it does not
// know how deep any particular command interprets the segments. That
// interpretation stays with the caller, which is what lets commands
// with different reach share one path syntax.
package modelpath
