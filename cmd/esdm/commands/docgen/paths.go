package docgen

import (
	"path"
	"path/filepath"
	"strings"
)

// segment is one level of an element's address: its kind
// and its name, as the reference notation writes them.
type segment struct {
	kind string
	name string
}

// reference renders the segments in the esdm: notation,
// kind=name per segment.
func reference(segments []segment) string {
	parts := make([]string, 0, len(segments))
	for _, s := range segments {
		parts = append(parts, s.kind+"="+s.name)
	}
	return "esdm:" + strings.Join(parts, "/")
}

// pagePath maps the segments to the page's path inside the
// output directory: one directory or file per segment, named
// kind_name - the reference with an underscore where it has
// an equals sign - and README.md for a container.
func pagePath(segments []segment) string {
	parts := make([]string, 0, len(segments)+1)
	for _, s := range segments {
		parts = append(parts, s.kind+"_"+s.name)
	}
	last := segments[len(segments)-1]
	if containerKinds[last.kind] {
		return path.Join(append(parts, "README.md")...)
	}
	parts[len(parts)-1] += ".md"
	return path.Join(parts...)
}

// relativeLink returns the link from one page to another,
// both given as slash-separated paths inside the output
// directory.
func relativeLink(fromPage, toPage string) string {
	fromDir := path.Dir(fromPage)
	rel, err := filepath.Rel(filepath.FromSlash(fromDir), filepath.FromSlash(toPage))
	if err != nil {
		return toPage
	}
	return filepath.ToSlash(rel)
}
