package resolver

import (
	"fmt"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
	"github.com/thenativeweb/esdm/parser"
)

const (
	coreAPIVersion               = "schema.esdm.io/core/v1"
	domainStorytellingAPIVersion = "schema.esdm.io/domain-storytelling/v1"
	givenWhenThenAPIVersion      = "schema.esdm.io/given-when-then/v1"
)

// Resolve collects all documents from the parsed files,
// dispatches them by (apiVersion, kind), and assembles a
// Model. Duplicate names within the same (apiVersion,
// kind, scope) namespace produce
// esdm/structure/duplicate-name diagnostics; the first
// occurrence wins in the Model. Two entities of the same
// kind that share a bare name but live in different
// scopes (e.g. two `register` commands belonging to
// different aggregates) are not considered duplicates -
// they coexist in the index under distinct composite
// keys.
//
// Documents whose apiVersion is unknown are skipped
// silently - the parser already emits
// esdm/structure/unknown-api-version for those. Same
// story for documents whose kind is unknown: the
// parser's schema validation catches kind typos as
// esdm/structure/constraint-violation (the kind enum is
// a closed set in every schema), so the resolver does
// not emit a second diagnostic.
//
// After indexing, Resolve walks every entity again and
// verifies its cross-references (scope triples,
// bare-name references, structured event and command
// references, mappingEndpoints, and intra-document
// field references). Failures become
// esdm/structure/unresolved-reference diagnostics with
// Levenshtein-based "did you mean?" hints.
func Resolve(files []*parser.ParsedFile) (*model.Model, []diag.Diagnostic) {
	m := model.NewModel()
	var diagnostics []diag.Diagnostic

	previous := make(map[string]map[string]ast.Node)

	for _, file := range files {
		for _, doc := range file.Documents {
			apiVersion, _ := doc.Field("apiVersion").Text()
			kindNode := doc.Field("kind")
			kind, ok := kindNode.Text()
			if !ok {
				continue
			}

			nameNode := doc.Field("name")
			name, ok := nameNode.Text()
			if !ok {
				continue
			}

			document := model.DocumentViewBase{Node: doc}
			key, ok := compositeKey(apiVersion, document, kind)
			if !ok {
				continue
			}

			namespaceKey := apiVersion + "/" + kind
			if previous[namespaceKey] == nil {
				previous[namespaceKey] = make(map[string]ast.Node)
			}

			if first, clash := previous[namespaceKey][key]; clash {
				diagnostics = append(diagnostics, diag.Diagnostic{
					RuleID:   "esdm/structure/duplicate-name",
					Severity: diag.SeverityError,
					Message:  fmt.Sprintf("duplicate %s %q", kind, name),
					Location: nameNode.Location(),
					Related: []diag.Related{
						{
							Message:  "first defined here",
							Location: first.Field("name").Location(),
						},
					},
				})
				continue
			}
			previous[namespaceKey][key] = doc

			switch apiVersion {
			case coreAPIVersion:
				dispatchCore(m, document, kind, key)
			case domainStorytellingAPIVersion:
				dispatchDomainStorytelling(m, document, kind, key)
			case givenWhenThenAPIVersion:
				dispatchGivenWhenThen(m, document, kind, key)
			default:
				// Unknown apiVersion - already flagged by
				// the parser; skip indexing.
			}
		}
	}

	diagnostics = append(diagnostics, resolveReferences(m)...)

	return m, diagnostics
}

// compositeKey routes a parsed document to its kind-
// specific key builder. Returns ("", false) for documents
// the resolver should not index (no readable name, or
// unsupported apiVersion+kind).
func compositeKey(apiVersion string, document model.DocumentViewBase, kind string) (string, bool) {
	switch apiVersion {
	case coreAPIVersion:
		return model.KeyForCoreDocument(document, kind)
	case domainStorytellingAPIVersion, givenWhenThenAPIVersion:
		// Both extensions key their documents by
		// domain/name: every scope variant carries `domain`
		// as its first field, so a domain-scoped composite
		// key identifies the document regardless of which
		// consistency unit it targets.
		name, ok := document.Name().Text()
		if !ok {
			return "", false
		}
		domain := model.ScopeText(document.Field("scope"), "domain")
		return domain + "/" + name, true
	}
	return "", false
}

// dispatchCore places a core-apiVersion document into
// the matching Model map. Unknown kinds are ignored: the
// parser has already flagged them as
// esdm/structure/constraint-violation because every
// schema's kind field is a closed enum.
func dispatchCore(m *model.Model, document model.DocumentViewBase, kind, key string) {
	switch kind {
	case "domain":
		m.Domains[key] = model.DomainView{DocumentViewBase: document}
	case "subdomain":
		m.Subdomains[key] = model.SubdomainView{DocumentViewBase: document}
	case "bounded-context":
		m.BoundedContexts[key] = model.BoundedContextView{DocumentViewBase: document}
	case "context-mapping":
		m.ContextMappings[key] = model.ContextMappingView{DocumentViewBase: document}
	case "aggregate":
		m.Aggregates[key] = model.AggregateView{DocumentViewBase: document}
	case "dynamic-consistency-boundary":
		m.DynamicConsistencyBoundaries[key] = model.DynamicConsistencyBoundaryView{DocumentViewBase: document}
	case "command":
		m.Commands[key] = model.CommandView{DocumentViewBase: document}
	case "event":
		m.Events[key] = model.EventView{DocumentViewBase: document}
	case "event-handler":
		m.EventHandlers[key] = model.EventHandlerView{DocumentViewBase: document}
	case "entity":
		m.Entities[key] = model.EntityView{DocumentViewBase: document}
	case "policy":
		m.Policies[key] = model.PolicyView{DocumentViewBase: document}
	case "process-manager":
		m.ProcessManagers[key] = model.ProcessManagerView{DocumentViewBase: document}
	case "read-model":
		m.ReadModels[key] = model.ReadModelView{DocumentViewBase: document}
	case "query":
		m.Queries[key] = model.QueryView{DocumentViewBase: document}
	case "value-object":
		m.ValueObjects[key] = model.ValueObjectView{DocumentViewBase: document}
	case "domain-service":
		m.DomainServices[key] = model.DomainServiceView{DocumentViewBase: document}
	case "actor":
		m.Actors[key] = model.ActorView{DocumentViewBase: document}
	case "external-system":
		m.ExternalSystems[key] = model.ExternalSystemView{DocumentViewBase: document}
	}
}

// dispatchDomainStorytelling places a
// domain-storytelling-apiVersion document into the
// matching Extensions map. Unknown kinds are ignored for
// the same reason as in dispatchCore.
func dispatchDomainStorytelling(m *model.Model, document model.DocumentViewBase, kind, key string) {
	switch kind {
	case "domain-story":
		m.Extensions.DomainStorytelling.Stories[key] = model.DomainStoryView{DocumentViewBase: document}
	}
}

// dispatchGivenWhenThen places a given-when-then-apiVersion
// document into the matching Extensions map. Unknown kinds
// are ignored for the same reason as in dispatchCore.
func dispatchGivenWhenThen(m *model.Model, document model.DocumentViewBase, kind, key string) {
	switch kind {
	case "feature":
		m.Extensions.GivenWhenThen.Features[key] = model.FeatureView{DocumentViewBase: document}
	}
}
