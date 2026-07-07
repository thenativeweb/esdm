package ast

import (
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/thenativeweb/esdm/diag"
)

// Kind classifies a Node by the underlying YAML kind.
type Kind int

const (
	KindInvalid Kind = iota
	KindMapping
	KindSequence
	KindScalar
)

// Node is a position-aware wrapper around a yaml.Node. A
// zero Node (Raw == nil) represents a missing path - all
// navigation methods return a zero Node rather than
// panicking, so expressions like
//
//	event.Field("name").Field("display").Text()
//
// are always safe to write.
type Node struct {
	Raw  *yaml.Node
	File string
}

// NewNode wraps a yaml.Node with its originating file
// name. If raw is a document node, it is unwrapped to its
// single content child.
func NewNode(file string, raw *yaml.Node) Node {
	if raw == nil {
		return Node{}
	}

	if raw.Kind == yaml.DocumentNode {
		if len(raw.Content) == 0 {
			return Node{}
		}
		raw = raw.Content[0]
	}

	return Node{Raw: raw, File: file}
}

// Exists reports whether the Node refers to an actual
// YAML value.
func (n Node) Exists() bool {
	return n.Raw != nil
}

// Kind returns the structural kind of the wrapped YAML
// node.
func (n Node) Kind() Kind {
	if n.Raw == nil {
		return KindInvalid
	}

	switch n.Raw.Kind {
	case yaml.MappingNode:
		return KindMapping
	case yaml.SequenceNode:
		return KindSequence
	case yaml.ScalarNode:
		return KindScalar
	default:
		return KindInvalid
	}
}

// Location returns the source position of the Node, or a
// zero Location if the Node is missing.
func (n Node) Location() diag.Location {
	if n.Raw == nil {
		return diag.Location{}
	}

	return diag.Location{
		File:   n.File,
		Line:   n.Raw.Line,
		Column: n.Raw.Column,
	}
}

// Field returns the value of the given key in a mapping.
// It returns a zero Node if the receiver is not a mapping
// or if the key is absent.
func (n Node) Field(key string) Node {
	if n.Raw == nil || n.Raw.Kind != yaml.MappingNode {
		return Node{}
	}

	content := n.Raw.Content
	for i := 0; i+1 < len(content); i += 2 {
		if content[i].Value == key {
			return Node{Raw: content[i+1], File: n.File}
		}
	}

	return Node{}
}

// HasField reports whether the mapping contains a value
// for the given key.
func (n Node) HasField(key string) bool {
	return n.Field(key).Exists()
}

// At returns the i-th element of a sequence. It returns a
// zero Node if the receiver is not a sequence or if i is
// out of range.
func (n Node) At(i int) Node {
	if n.Raw == nil || n.Raw.Kind != yaml.SequenceNode {
		return Node{}
	}
	if i < 0 || i >= len(n.Raw.Content) {
		return Node{}
	}

	return Node{Raw: n.Raw.Content[i], File: n.File}
}

// Seq returns all elements of a sequence, or nil if the
// receiver is not a sequence.
func (n Node) Seq() []Node {
	if n.Raw == nil || n.Raw.Kind != yaml.SequenceNode {
		return nil
	}

	out := make([]Node, len(n.Raw.Content))
	for i, c := range n.Raw.Content {
		out[i] = Node{Raw: c, File: n.File}
	}

	return out
}

// MapEntry is a single (key, value) pair from a mapping.
type MapEntry struct {
	Key   Node
	Value Node
}

// Entries returns all (key, value) pairs of a mapping, in
// file order. It returns nil if the receiver is not a
// mapping.
func (n Node) Entries() []MapEntry {
	if n.Raw == nil || n.Raw.Kind != yaml.MappingNode {
		return nil
	}

	content := n.Raw.Content
	out := make([]MapEntry, 0, len(content)/2)
	for i := 0; i+1 < len(content); i += 2 {
		out = append(out, MapEntry{
			Key:   Node{Raw: content[i], File: n.File},
			Value: Node{Raw: content[i+1], File: n.File},
		})
	}

	return out
}

// Text returns the string value of a scalar node. ok is
// true iff the receiver is a scalar with an explicit
// string tag (!!str) or an implicit/unspecified tag.
func (n Node) Text() (string, bool) {
	if n.Raw == nil || n.Raw.Kind != yaml.ScalarNode {
		return "", false
	}

	tag := n.Raw.Tag
	if tag != "" && tag != "!!str" {
		return "", false
	}

	return n.Raw.Value, true
}

// Int returns the integer value of a scalar node tagged
// as !!int. ok is false otherwise.
func (n Node) Int() (int64, bool) {
	if n.Raw == nil || n.Raw.Kind != yaml.ScalarNode {
		return 0, false
	}
	if n.Raw.Tag != "!!int" {
		return 0, false
	}

	v, err := strconv.ParseInt(n.Raw.Value, 10, 64)
	if err != nil {
		return 0, false
	}

	return v, true
}

// Bool returns the boolean value of a scalar node tagged
// as !!bool. ok is false otherwise.
func (n Node) Bool() (bool, bool) {
	if n.Raw == nil || n.Raw.Kind != yaml.ScalarNode {
		return false, false
	}
	if n.Raw.Tag != "!!bool" {
		return false, false
	}

	v, err := strconv.ParseBool(n.Raw.Value)
	if err != nil {
		return false, false
	}

	return v, true
}
