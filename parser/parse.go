package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
)

// ParsedFile is the outcome of parsing one source file.
// Each YAML document inside the file becomes one element
// of Documents.
type ParsedFile struct {
	Path      string
	Documents []ast.Node
}

// Parse reads a file, splits it into YAML documents, and
// validates each document against the schema its
// apiVersion refers to (core or any embedded extension).
// YAML syntax errors, unknown apiVersions, and schema
// violations are returned as diagnostics. Only outright
// I/O or infrastructure failures become Go errors.
func Parse(path string) (*ParsedFile, []diag.Diagnostic, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	validators, err := schemasByAPIVersion()
	if err != nil {
		return nil, nil, err
	}

	parsed := &ParsedFile{Path: path}
	var diagnostics []diag.Diagnostic

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	for {
		var raw yaml.Node
		err := decoder.Decode(&raw)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			diagnostics = append(diagnostics, diag.Diagnostic{
				RuleID:   "esdm/structure/yaml-syntax-error",
				Severity: diag.SeverityError,
				Message:  err.Error(),
				Location: diag.Location{File: path},
			})
			break
		}

		node := ast.NewNode(path, &raw)
		parsed.Documents = append(parsed.Documents, node)

		apiVersionNode := node.Field("apiVersion")
		apiVersion, _ := apiVersionNode.Text()

		validator, isKnown := validators[apiVersion]
		if !isKnown {
			diagnostics = append(diagnostics, unknownAPIVersionDiagnostic(path, node, apiVersionNode, apiVersion, validators))
			continue
		}

		var generic any
		err = raw.Decode(&generic)
		if err != nil {
			continue
		}

		err = validator.Validate(generic)
		if err != nil {
			var validationError *jsonschema.ValidationError
			if errors.As(err, &validationError) {
				diagnostics = append(diagnostics, validationDiagnostics(validationError, node)...)
			}
		}
	}

	return parsed, diagnostics, nil
}

// unknownAPIVersionDiagnostic builds the diagnostic used
// when a document's apiVersion does not match any
// compiled schema - either because the field is missing
// or because it refers to a URL the linter does not
// know. The diagnostic's Location points at the
// apiVersion field when present, or at the document
// root as a fallback.
func unknownAPIVersionDiagnostic(path string, document ast.Node, apiVersionNode ast.Node, apiVersion string, validators map[string]*jsonschema.Schema) diag.Diagnostic {
	location := apiVersionNode.Location()
	if location.IsZero() {
		location = document.Location()
	}
	if location.IsZero() {
		location = diag.Location{File: path}
	}

	message := fmt.Sprintf("unknown apiVersion %q; expected one of: %s", apiVersion, knownAPIVersions(validators))
	if apiVersion == "" {
		message = fmt.Sprintf("document has no apiVersion; expected one of: %s", knownAPIVersions(validators))
	}

	return diag.Diagnostic{
		RuleID:   "esdm/structure/unknown-api-version",
		Severity: diag.SeverityError,
		Message:  message,
		Location: location,
	}
}

func knownAPIVersions(validators map[string]*jsonschema.Schema) string {
	names := make([]string, 0, len(validators))
	for k := range validators {
		names = append(names, k)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}
