package schema

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"sort"

	"gopkg.in/yaml.v3"
)

// currentVersion is the schema version every embedded
// schema currently pins to. When a new major version is
// published, add a parallel constant and a second embed
// directive; until then, a single version keeps the
// logic simple.
const currentVersion = "v1"

// schemaFS embeds the full schema tree. Each top-level
// directory is one schema group - "core" is reserved for
// the ESDM core schema; every other directory is an
// extension. Directories are expected to contain files
// named v<N>.yaml, one per major version.
//
// When a new extension directory is added to this
// package, extend the embed pattern below to include it.
// The list is short and explicit on purpose: embed does
// not offer dynamic discovery, and we prefer the
// compile-time signal of an extra line over a
// hand-maintained registry file.
//
//go:embed core
//go:embed domain-storytelling
//go:embed given-when-then
var schemaFS embed.FS

// Core returns the ESDM core schema as YAML bytes.
func Core() []byte {
	data, err := fs.ReadFile(schemaFS, path.Join("core", currentVersion+".yaml"))
	if err != nil {
		// Unreachable: the file is embedded at compile
		// time; any error here would indicate a broken
		// build, which we want to surface loudly.
		panic(err)
	}

	out := make([]byte, len(data))
	copy(out, data)
	return out
}

// Extension is a single extension schema: its name (equal
// to the directory name at the schema root) and the
// schema bytes for the current version.
type Extension struct {
	Name  string
	Bytes []byte
}

// Extensions returns all embedded extension schemas,
// meaning every top-level directory other than "core".
// Extensions are returned sorted by name so that callers
// emit deterministic output.
func Extensions() ([]Extension, error) {
	entries, err := fs.ReadDir(schemaFS, ".")
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if entry.Name() == "core" {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	out := make([]Extension, 0, len(names))
	for _, name := range names {
		data, err := fs.ReadFile(schemaFS, path.Join(name, currentVersion+".yaml"))
		if err != nil {
			return nil, err
		}

		clone := make([]byte, len(data))
		copy(clone, data)

		out = append(out, Extension{Name: name, Bytes: clone})
	}

	return out, nil
}

// CurrentVersion returns the schema version string used
// by every embedded schema. Consumers that need to
// reconstruct the "<name>/v<N>.yaml" path layout on disk
// use this together with the extension name.
func CurrentVersion() string {
	return currentVersion
}

// File describes one embedded schema file as it should
// appear under a project's `schemas/` directory. Path is
// relative to that directory, so writing the file is a
// straightforward join with the project root.
type File struct {
	Path  string
	Bytes []byte
}

// Files returns every embedded schema file - the core
// and every extension - with paths relative to a
// project's `schemas/` directory. Files are sorted by
// path so callers get deterministic output and can do
// stable byte-for-byte comparisons against on-disk
// copies.
func Files() []File {
	out := []File{
		{
			Path:  path.Join("core", currentVersion+".yaml"),
			Bytes: Core(),
		},
	}

	extensions, err := Extensions()
	if err != nil {
		// Same reasoning as Core: any error here means
		// the embedded FS is broken at build time.
		panic(err)
	}

	for _, extension := range extensions {
		out = append(out, File{
			Path:  path.Join(extension.Name, currentVersion+".yaml"),
			Bytes: extension.Bytes,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Path < out[j].Path
	})

	return out
}

// semverPattern matches the SemVer 2.0.0 core grammar
// (MAJOR.MINOR.PATCH); pre-release and build metadata are
// not used by ESDM revisions today, so we keep the
// pattern strict.
var semverPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

// Revision parses the `x-esdm-schema-revision` value out
// of a schema YAML and returns it as a SemVer string. The
// embedded schemas guarantee the field's presence; local
// copies that lack it (or carry a non-SemVer value) are
// rejected.
func Revision(yamlBytes []byte) (string, error) {
	var head struct {
		Revision string `yaml:"x-esdm-schema-revision"`
	}
	err := yaml.Unmarshal(yamlBytes, &head)
	if err != nil {
		return "", err
	}

	if head.Revision == "" {
		return "", errors.New("schema YAML is missing x-esdm-schema-revision")
	}
	if !semverPattern.MatchString(head.Revision) {
		return "", fmt.Errorf("schema revision %q is not a SemVer (MAJOR.MINOR.PATCH)", head.Revision)
	}

	return head.Revision, nil
}

// APIVersions returns the `apiVersion` value of every
// embedded schema - that is, the `$id` of each schema
// without the `https://` scheme. Documents declare which
// schema they validate against via `apiVersion`, and the
// linter consults this list to verify that every document
// targets a schema this binary actually carries.
func APIVersions() ([]string, error) {
	files := Files()

	out := make([]string, 0, len(files))
	for _, file := range files {
		var head struct {
			ID string `yaml:"$id"`
		}
		err := yaml.Unmarshal(file.Bytes, &head)
		if err != nil {
			return nil, fmt.Errorf("embedded schema %q: %w", file.Path, err)
		}

		if head.ID == "" {
			return nil, fmt.Errorf("embedded schema %q is missing $id", file.Path)
		}

		const httpsPrefix = "https://"
		if len(head.ID) <= len(httpsPrefix) || head.ID[:len(httpsPrefix)] != httpsPrefix {
			return nil, fmt.Errorf("embedded schema %q has $id %q without https:// prefix", file.Path, head.ID)
		}
		out = append(out, head.ID[len(httpsPrefix):])
	}

	sort.Strings(out)
	return out, nil
}
