package ast

import (
	"strconv"
	"strings"
)

// FollowPointer resolves a JSON Pointer (RFC 6901) against
// the receiver. The pointer "" returns the root; otherwise
// each "/"-separated segment is used to descend - by key
// in a mapping, by index in a sequence.
//
// If any segment cannot be followed (missing key, out-of-
// range index, wrong kind), FollowPointer returns the last
// successfully resolved Node. This matches the practical
// need of schema validators: if a validator says a
// required field is missing, the pointer walks to the
// parent mapping and stops there, which is exactly the
// position a diagnostic should point at.
func FollowPointer(n Node, pointer string) Node {
	if pointer == "" {
		return n
	}
	if !strings.HasPrefix(pointer, "/") {
		return n
	}

	current := n
	for _, rawSegment := range strings.Split(pointer[1:], "/") {
		segment := unescapePointerSegment(rawSegment)

		var next Node
		switch current.Kind() {
		case KindMapping:
			next = current.Field(segment)
		case KindSequence:
			i, err := strconv.Atoi(segment)
			if err != nil {
				return current
			}
			next = current.At(i)
		default:
			return current
		}

		if !next.Exists() {
			return current
		}
		current = next
	}

	return current
}

// unescapePointerSegment reverses the RFC 6901 escaping:
// ~1 -> /, ~0 -> ~. The order matters - ~1 must be expanded
// before ~0, otherwise "~01" would first become "~1" and
// then "/", losing the tilde.
func unescapePointerSegment(segment string) string {
	segment = strings.ReplaceAll(segment, "~1", "/")
	segment = strings.ReplaceAll(segment, "~0", "~")

	return segment
}
