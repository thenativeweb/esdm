package refgen

import (
	"fmt"
	"path"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/thenativeweb/esdm/schema"
)

// metaKeys are top-level schema bookkeeping fields that have no place
// in a reader-facing snippet. They describe how the schema is hosted,
// versioned, and laid out, not what an ESDM document looks like.
var metaKeys = []string{
	"$id",
	"$schema",
	"$defs",
	"x-esdm-schema-revision",
	"x-esdm-file-suffix",
	"x-esdm-document-separator",
	"x-esdm-project-layout",
	"description",
}

// Snippets returns every reference snippet that the documentation site
// embeds. Keys are paths relative to the snippet base directory
// (documentation/snippets), in forward-slash form. Values are the raw
// YAML bytes the snippet file must contain.
//
// The map is the single source of truth for both the cmd/refgen CLI
// (which writes the entries to disk) and the documentation sync test
// (which compares the entries against the committed files).
func Snippets() (map[string][]byte, error) {
	out := make(map[string][]byte)

	err := generateCore(out)
	if err != nil {
		return nil, fmt.Errorf("core: %w", err)
	}

	err = generateExtensions(out)
	if err != nil {
		return nil, fmt.Errorf("extensions: %w", err)
	}

	return out, nil
}

// SortedPaths returns the snippet paths in lexicographic order, useful
// when callers need a stable iteration order (e.g. for diff reporting).
func SortedPaths(snippets map[string][]byte) []string {
	keys := make([]string, 0, len(snippets))
	for k := range snippets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func generateCore(out map[string][]byte) error {
	root, err := parseSchema(schema.Core())
	if err != nil {
		return err
	}
	defs := collectDefs(root)
	allOfNode := mapValue(root, "allOf")
	if allOfNode == nil {
		return fmt.Errorf("missing allOf")
	}

	for _, entry := range allOfNode.Content {
		kindName := kindFromBranch(entry)
		if kindName == "" {
			continue
		}
		thenNode := mapValue(entry, "then")
		if thenNode == nil {
			continue
		}
		snippet := coreKindSchema(root, kindName, thenNode, defs)
		key := path.Join("reference", "core-schema", kindName+".yaml")
		data, err := marshalNode(snippet)
		if err != nil {
			return err
		}
		out[key] = data
	}
	return nil
}

// coreKindSchema returns the snippet for one core kind. When the
// kind's `then` block is a mapping (i.e. the kind has its own fields),
// the snippet is just that block - the common top-level shape is
// documented once on the Core Schema overview, so repeating it on
// every kind would be noise. When the `then` block is `true` (the
// kind has no kind-specific fields, as is the case for `domain`), the
// snippet falls back to the top-level shape with `kind.const` set to
// the kind name; otherwise the page would have nothing to show.
func coreKindSchema(root *yaml.Node, kindName string, thenNode *yaml.Node, defs map[string]*yaml.Node) *yaml.Node {
	if thenNode != nil && thenNode.Kind == yaml.MappingNode {
		return resolveRefs(thenNode, defs, map[string]bool{})
	}

	base := stripKeys(root, metaKeys...)
	base = stripKeys(base, "allOf")

	properties := mapValue(base, "properties")
	if properties != nil {
		properties = pinKindConst(properties, kindName)
		base = replaceProperty(base, "properties", properties)
	}
	return resolveRefs(base, defs, map[string]bool{})
}

// pinKindConst replaces the `kind` property's `enum` with a single
// `const` that pins the kind name. The schema's top-level `kind` is
// an enum across all valid kinds; per-kind snippets need it pinned to
// the one kind they describe.
func pinKindConst(properties *yaml.Node, kindName string) *yaml.Node {
	out := cloneNode(properties)
	for i := 0; i+1 < len(out.Content); i += 2 {
		if out.Content[i].Value != "kind" {
			continue
		}
		out.Content[i+1] = &yaml.Node{
			Kind: yaml.MappingNode,
			Tag:  "!!map",
			Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Tag: "!!str", Value: "const"},
				{Kind: yaml.ScalarNode, Tag: "!!str", Value: kindName},
			},
		}
		break
	}
	return out
}

// replaceProperty returns a clone of node with the value at key
// replaced by replacement.
func replaceProperty(node *yaml.Node, key string, replacement *yaml.Node) *yaml.Node {
	out := cloneNode(node)
	for i := 0; i+1 < len(out.Content); i += 2 {
		if out.Content[i].Value == key {
			out.Content[i+1] = replacement
			return out
		}
	}
	return out
}

func generateExtensions(out map[string][]byte) error {
	extensions, err := schema.Extensions()
	if err != nil {
		return err
	}
	for _, extension := range extensions {
		root, err := parseSchema(extension.Bytes)
		if err != nil {
			return fmt.Errorf("%s: %w", extension.Name, err)
		}
		defs := collectDefs(root)

		properties := mapValue(root, "properties")
		if properties == nil {
			return fmt.Errorf("%s: missing properties", extension.Name)
		}
		kindEntry := mapValue(properties, "kind")
		if kindEntry == nil {
			return fmt.Errorf("%s: missing kind", extension.Name)
		}
		kindConst := mapValue(kindEntry, "const")
		if kindConst == nil {
			return fmt.Errorf("%s: missing kind.const", extension.Name)
		}
		kindName := kindConst.Value

		stripped := stripKeys(root, metaKeys...)
		resolved := resolveRefs(stripped, defs, map[string]bool{})
		key := path.Join("reference", "extensions", extension.Name, kindName+".yaml")
		data, err := marshalNode(resolved)
		if err != nil {
			return err
		}
		out[key] = data

		if extension.Name == "given-when-then" {
			err := generateScenarioVariants(out, root, defs)
			if err != nil {
				return fmt.Errorf("%s: scenario variants: %w", extension.Name, err)
			}
		}
	}
	return nil
}

// generateScenarioVariants emits one snippet per variant of the
// `feature.scenarios[]` item shape. Each variant is the merge of the
// base item schema with the variant-specific override carried in the
// surrounding feature's allOf, with all $refs resolved. The result is
// what a single scenario looks like under each of the four feature
// variants the schema admits.
func generateScenarioVariants(out map[string][]byte, root *yaml.Node, defs map[string]*yaml.Node) error {
	properties := mapValue(root, "properties")
	if properties == nil {
		return fmt.Errorf("missing properties")
	}
	scenarios := mapValue(properties, "scenarios")
	if scenarios == nil {
		return fmt.Errorf("missing scenarios")
	}
	baseItem := mapValue(scenarios, "items")
	if baseItem == nil {
		return fmt.Errorf("missing scenarios.items")
	}

	allOfNode := mapValue(root, "allOf")
	if allOfNode == nil {
		return fmt.Errorf("missing allOf")
	}

	for _, branch := range allOfNode.Content {
		variant := scenarioVariantKey(branch)
		if variant == "" {
			continue
		}
		override := scenarioOverride(branch)
		merged := mergeScenarioItem(baseItem, override)
		resolved := resolveRefs(merged, defs, map[string]bool{})

		key := path.Join("reference", "extensions", "given-when-then", "scenario-"+variant+".yaml")
		data, err := marshalNode(resolved)
		if err != nil {
			return err
		}
		out[key] = data
	}
	return nil
}

// scenarioVariantKey reads the discriminator from a feature allOf
// branch. The branch's `if.properties.scope.required` contains exactly
// one variant marker; that marker (kebab-cased) keys the snippet.
func scenarioVariantKey(branch *yaml.Node) string {
	ifNode := mapValue(branch, "if")
	if ifNode == nil {
		return ""
	}
	scopeProperties := mapValue(mapValue(ifNode, "properties"), "scope")
	if scopeProperties == nil {
		return ""
	}
	required := mapValue(scopeProperties, "required")
	if required == nil || required.Kind != yaml.SequenceNode || len(required.Content) == 0 {
		return ""
	}
	return camelToKebab(required.Content[0].Value)
}

// scenarioOverride returns the scenarios.items node from a feature
// allOf branch's `then`. That node holds only the variant-specific
// shape of `given`, `when`, and `then`; the rest comes from the base.
func scenarioOverride(branch *yaml.Node) *yaml.Node {
	thenNode := mapValue(branch, "then")
	if thenNode == nil {
		return nil
	}
	scenarios := mapValue(mapValue(thenNode, "properties"), "scenarios")
	if scenarios == nil {
		return nil
	}
	return mapValue(scenarios, "items")
}

// mergeScenarioItem deep-merges the variant override into the base
// scenario item schema. Override keys win for nested mappings; arrays
// from the override replace the base outright.
func mergeScenarioItem(base, override *yaml.Node) *yaml.Node {
	if override == nil {
		return cloneNode(base)
	}
	if base == nil {
		return cloneNode(override)
	}
	return mergeNodes(base, override)
}

func mergeNodes(base, override *yaml.Node) *yaml.Node {
	if base == nil {
		return cloneNode(override)
	}
	if override == nil {
		return cloneNode(base)
	}
	if base.Kind != yaml.MappingNode || override.Kind != yaml.MappingNode {
		return cloneNode(override)
	}
	out := &yaml.Node{
		Kind:  yaml.MappingNode,
		Tag:   base.Tag,
		Style: base.Style,
	}
	seen := map[string]bool{}
	for i := 0; i+1 < len(base.Content); i += 2 {
		key := base.Content[i].Value
		seen[key] = true
		if overrideValue := mapValue(override, key); overrideValue != nil {
			out.Content = append(out.Content,
				cloneNode(base.Content[i]),
				mergeNodes(base.Content[i+1], overrideValue),
			)
			continue
		}
		out.Content = append(out.Content,
			cloneNode(base.Content[i]),
			cloneNode(base.Content[i+1]),
		)
	}
	for i := 0; i+1 < len(override.Content); i += 2 {
		key := override.Content[i].Value
		if seen[key] {
			continue
		}
		out.Content = append(out.Content,
			cloneNode(override.Content[i]),
			cloneNode(override.Content[i+1]),
		)
	}
	return out
}

func camelToKebab(s string) string {
	out := make([]byte, 0, len(s)+4)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				out = append(out, '-')
			}
			out = append(out, c+('a'-'A'))
			continue
		}
		out = append(out, c)
	}
	return string(out)
}

func parseSchema(data []byte) (*yaml.Node, error) {
	var root yaml.Node
	err := yaml.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}
	stripComments(&root)
	if len(root.Content) == 0 {
		return nil, fmt.Errorf("empty schema document")
	}
	return root.Content[0], nil
}

func kindFromBranch(branch *yaml.Node) string {
	ifNode := mapValue(branch, "if")
	if ifNode == nil {
		return ""
	}
	properties := mapValue(ifNode, "properties")
	if properties == nil {
		return ""
	}
	kind := mapValue(properties, "kind")
	if kind == nil {
		return ""
	}
	constNode := mapValue(kind, "const")
	if constNode == nil {
		return ""
	}
	return constNode.Value
}

func mapValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

func stripComments(node *yaml.Node) {
	if node == nil {
		return
	}
	node.HeadComment = ""
	node.LineComment = ""
	node.FootComment = ""
	for _, c := range node.Content {
		stripComments(c)
	}
}

func stripKeys(node *yaml.Node, keys ...string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return node
	}
	skip := make(map[string]bool, len(keys))
	for _, k := range keys {
		skip[k] = true
	}
	out := &yaml.Node{
		Kind:  node.Kind,
		Tag:   node.Tag,
		Style: node.Style,
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if skip[node.Content[i].Value] {
			continue
		}
		out.Content = append(out.Content, node.Content[i], node.Content[i+1])
	}
	return out
}

func marshalNode(node *yaml.Node) ([]byte, error) {
	return yaml.Marshal(node)
}

// collectDefs returns a lookup table from a `$defs` entry name to the
// node defining it. The reference snippets inline every `$ref:
// "#/$defs/<name>"` against this table so readers see the resolved
// shape directly, without having to chase references.
func collectDefs(root *yaml.Node) map[string]*yaml.Node {
	defs := map[string]*yaml.Node{}
	defsNode := mapValue(root, "$defs")
	if defsNode == nil || defsNode.Kind != yaml.MappingNode {
		return defs
	}
	for i := 0; i+1 < len(defsNode.Content); i += 2 {
		defs[defsNode.Content[i].Value] = defsNode.Content[i+1]
	}
	return defs
}

// resolveRefs walks the node tree and replaces every `$ref:
// "#/$defs/<name>"` with a deep copy of the referenced node, keeping
// any sibling keys around the `$ref` (rare in this schema, but legal
// in JSON Schema). A visited set guards against cycles even though
// the ESDM schema itself is acyclic - a future addition should not be
// able to make the generator hang.
func resolveRefs(node *yaml.Node, defs map[string]*yaml.Node, visited map[string]bool) *yaml.Node {
	if node == nil {
		return nil
	}
	switch node.Kind {
	case yaml.MappingNode:
		if name := refName(node); name != "" {
			target, ok := defs[name]
			if !ok || visited[name] {
				return cloneNode(node)
			}
			next := cloneVisited(visited)
			next[name] = true
			return resolveRefs(target, defs, next)
		}
		out := &yaml.Node{
			Kind:  node.Kind,
			Tag:   node.Tag,
			Style: node.Style,
		}
		for i := 0; i+1 < len(node.Content); i += 2 {
			out.Content = append(out.Content,
				cloneNode(node.Content[i]),
				resolveRefs(node.Content[i+1], defs, visited),
			)
		}
		return out
	case yaml.SequenceNode:
		out := &yaml.Node{
			Kind:  node.Kind,
			Tag:   node.Tag,
			Style: node.Style,
		}
		for _, item := range node.Content {
			out.Content = append(out.Content, resolveRefs(item, defs, visited))
		}
		return out
	default:
		return cloneNode(node)
	}
}

func refName(node *yaml.Node) string {
	if node == nil || node.Kind != yaml.MappingNode {
		return ""
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value != "$ref" {
			continue
		}
		const prefix = "#/$defs/"
		v := node.Content[i+1].Value
		if len(v) > len(prefix) && v[:len(prefix)] == prefix {
			return v[len(prefix):]
		}
	}
	return ""
}

func cloneNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	out := &yaml.Node{
		Kind:  node.Kind,
		Tag:   node.Tag,
		Value: node.Value,
		Style: node.Style,
	}
	for _, c := range node.Content {
		out.Content = append(out.Content, cloneNode(c))
	}
	return out
}

func cloneVisited(in map[string]bool) map[string]bool {
	out := make(map[string]bool, len(in)+1)
	for k, v := range in {
		out[k] = v
	}
	return out
}
