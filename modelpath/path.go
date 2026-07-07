package modelpath

import (
	"fmt"
	"strings"
)

// Path is a parsed model-path: a sequence of segments drawn from the
// natural model hierarchy. The empty Path selects the whole model.
type Path struct {
	Segments []string
}

// ParsePath splits a slash-separated model-path string into its
// segments. The empty string yields an empty Path. Trailing slashes
// are tolerated; leading slashes and empty interior segments are
// rejected so the path always identifies a single, unambiguous model
// node.
func ParsePath(raw string) (Path, error) {
	if raw == "" {
		return Path{}, nil
	}
	if strings.HasPrefix(raw, "/") {
		return Path{}, fmt.Errorf("invalid path %q: leading slash", raw)
	}

	raw = strings.TrimSuffix(raw, "/")
	segments := strings.Split(raw, "/")
	for _, segment := range segments {
		if segment == "" {
			return Path{}, fmt.Errorf("invalid path %q: empty segment", raw)
		}
	}

	return Path{Segments: segments}, nil
}
