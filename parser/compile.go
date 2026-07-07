package parser

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"

	"github.com/thenativeweb/esdm/schema"
)

var (
	schemasOnce sync.Once
	schemasMap  map[string]*jsonschema.Schema
	schemasErr  error
)

// schemasByAPIVersion returns every compiled schema the
// linter knows about - the core schema plus every
// embedded extension schema - keyed by the apiVersion
// string the document uses to identify itself (the
// schema's `$id` URL with the `https://` prefix
// stripped). Compilation happens once, lazily, on the
// first call.
func schemasByAPIVersion() (map[string]*jsonschema.Schema, error) {
	schemasOnce.Do(func() {
		schemasMap, schemasErr = compileAllSchemas()
	})
	return schemasMap, schemasErr
}

func compileAllSchemas() (map[string]*jsonschema.Schema, error) {
	out := make(map[string]*jsonschema.Schema)

	err := compileSchemaInto(out, schema.Core())
	if err != nil {
		return nil, fmt.Errorf("cannot compile core schema: %w", err)
	}

	extensions, err := schema.Extensions()
	if err != nil {
		return nil, fmt.Errorf("cannot enumerate extension schemas: %w", err)
	}

	for _, extension := range extensions {
		err := compileSchemaInto(out, extension.Bytes)
		if err != nil {
			return nil, fmt.Errorf("cannot compile extension schema %q: %w", extension.Name, err)
		}
	}

	return out, nil
}

// compileSchemaInto parses one schema YAML blob,
// extracts its `$id`, compiles it, and stores the
// compiled schema under the apiVersion key (the `$id`
// URL without the scheme).
func compileSchemaInto(out map[string]*jsonschema.Schema, schemaBytes []byte) error {
	var decoded any
	err := yaml.Unmarshal(schemaBytes, &decoded)
	if err != nil {
		return fmt.Errorf("cannot parse schema yaml: %w", err)
	}

	root, ok := decoded.(map[string]any)
	if !ok {
		return errors.New("schema root is not a mapping")
	}

	idRaw, ok := root["$id"]
	if !ok {
		return errors.New("schema is missing $id")
	}

	idURL, ok := idRaw.(string)
	if !ok {
		return errors.New("schema $id is not a string")
	}

	compiler := jsonschema.NewCompiler()
	err = compiler.AddResource(idURL, decoded)
	if err != nil {
		return fmt.Errorf("cannot register schema at %s: %w", idURL, err)
	}

	compiled, err := compiler.Compile(idURL)
	if err != nil {
		return fmt.Errorf("cannot compile schema at %s: %w", idURL, err)
	}

	apiVersion := strings.TrimPrefix(idURL, "https://")
	out[apiVersion] = compiled
	return nil
}
