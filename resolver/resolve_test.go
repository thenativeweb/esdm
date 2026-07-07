package resolver_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/parser"
	"github.com/thenativeweb/esdm/resolver"
)

const parentsYAML = `apiVersion: schema.esdm.io/core/v1
kind: domain
name: commerce
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: ordering
scope:
  domain: commerce
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: order
scope:
  domain: commerce
  boundedContext: ordering
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
`

const eventTemplate = `apiVersion: schema.esdm.io/core/v1
kind: event
name: %s
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
`

const aggregateTemplate = `apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: %s
scope:
  domain: commerce
  boundedContext: ordering
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
`

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
}

func writeParents(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "parents.esdm.yaml")
	writeFile(t, path, parentsYAML)
	return path
}

func writeEvent(t *testing.T, dir, fileName, name string) string {
	t.Helper()

	path := filepath.Join(dir, fileName)
	writeFile(t, path, fmt.Sprintf(eventTemplate, name))
	return path
}

func writeAggregate(t *testing.T, dir, fileName, name string) string {
	t.Helper()

	path := filepath.Join(dir, fileName)
	writeFile(t, path, fmt.Sprintf(aggregateTemplate, name))
	return path
}

func parseAll(t *testing.T, paths ...string) []*parser.ParsedFile {
	t.Helper()

	var files []*parser.ParsedFile
	for _, p := range paths {
		parsed, _, err := parser.Parse(p)
		require.NoError(t, err)
		files = append(files, parsed)
	}
	return files
}

// syntheticFile builds a ParsedFile from raw YAML
// without going through schema validation. Used to test
// resolver paths that the parser would otherwise reject
// (most importantly the unknown-kind case).
func syntheticFile(t *testing.T, path, content string) *parser.ParsedFile {
	t.Helper()

	var raw yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(content), &raw))
	return &parser.ParsedFile{
		Path:      path,
		Documents: []ast.Node{ast.NewNode(path, &raw)},
	}
}

func TestResolve(t *testing.T) {
	t.Run("indexes events by name and produces no diagnostics on a complete model", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		eventPathA := writeEvent(t, dir, "a.esdm.yaml", "order-placed")
		eventPathB := writeEvent(t, dir, "b.esdm.yaml", "order-shipped")
		files := parseAll(t, parents, eventPathA, eventPathB)

		m, diagnostics := resolver.Resolve(files)
		assert.Empty(t, diagnostics)
		_, ok := m.LookupEvent("commerce", "ordering", "order", "order-placed")
		assert.True(t, ok, "expected event order-placed in index")
		_, ok = m.LookupEvent("commerce", "ordering", "order", "order-shipped")
		assert.True(t, ok, "expected event order-shipped in index")
	})

	t.Run("indexes aggregates by name", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		aggregatePath := writeAggregate(t, dir, "order.esdm.yaml", "order")
		files := parseAll(t, parents, aggregatePath)

		m, _ := resolver.Resolve(files)
		// parents.yaml already declares the order aggregate,
		// so a second one with the same name produces a
		// duplicate-name diagnostic; we only assert that
		// the aggregate is in the index.
		_, ok := m.LookupAggregate("commerce", "ordering", "order")
		assert.True(t, ok, "expected aggregate order in index")
	})

	t.Run("reports a duplicate-name diagnostic when two events share a name", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		eventPathA := writeEvent(t, dir, "a.esdm.yaml", "order-placed")
		eventPathB := writeEvent(t, dir, "b.esdm.yaml", "order-placed")
		files := parseAll(t, parents, eventPathA, eventPathB)

		_, diagnostics := resolver.Resolve(files)

		var didFindDuplicate bool
		for _, d := range diagnostics {
			if d.RuleID == "esdm/structure/duplicate-name" {
				didFindDuplicate = true
				assert.Contains(t, d.Message, "order-placed")
				require.Len(t, d.Related, 1)
				assert.NotEmpty(t, d.Related[0].Location.File)
			}
		}
		assert.True(t, didFindDuplicate, "expected duplicate-name diagnostic, got %+v", diagnostics)
	})

	t.Run("does not consider duplicates across different kinds", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		eventPath := writeEvent(t, dir, "event.esdm.yaml", "shared")
		aggregatePath := writeAggregate(t, dir, "aggregate.esdm.yaml", "shared")
		files := parseAll(t, parents, eventPath, aggregatePath)

		_, diagnostics := resolver.Resolve(files)
		for _, d := range diagnostics {
			assert.NotEqual(t, "esdm/structure/duplicate-name", d.RuleID, "got %+v", diagnostics)
		}
	})

	t.Run("does not consider duplicates across different aggregates", func(t *testing.T) {
		// Two events with the same bare name in different
		// aggregates (within the same bounded context) live in
		// different DDD scopes and must coexist - the same
		// reasoning applies to commands across aggregates and
		// to actors / read-models / etc. across bounded
		// contexts. The resolver must not collapse them via
		// global-uniqueness-per-kind.
		dir := t.TempDir()
		parents := writeParents(t, dir)
		extras := filepath.Join(dir, "extras.esdm.yaml")
		writeFile(t, extras, `apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: shipment
scope:
  domain: commerce
  boundedContext: ordering
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
`)
		eventInOrder := filepath.Join(dir, "order-event.esdm.yaml")
		writeFile(t, eventInOrder, `apiVersion: schema.esdm.io/core/v1
kind: event
name: created
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
`)
		eventInShipment := filepath.Join(dir, "shipment-event.esdm.yaml")
		writeFile(t, eventInShipment, `apiVersion: schema.esdm.io/core/v1
kind: event
name: created
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: shipment
data:
  type: object
`)
		files := parseAll(t, parents, extras, eventInOrder, eventInShipment)

		m, diagnostics := resolver.Resolve(files)

		for _, d := range diagnostics {
			assert.NotEqual(t, "esdm/structure/duplicate-name", d.RuleID, "unexpected duplicate-name: %+v", d)
		}

		orderCreated, ok := m.LookupEvent("commerce", "ordering", "order", "created")
		require.True(t, ok, "expected event scoped to (commerce, ordering, order) named created")
		assert.Equal(t, eventInOrder, orderCreated.Name().Location().File)

		shipmentCreated, ok := m.LookupEvent("commerce", "ordering", "shipment", "created")
		require.True(t, ok, "expected event scoped to (commerce, ordering, shipment) named created")
		assert.Equal(t, eventInShipment, shipmentCreated.Name().Location().File)
	})

	t.Run("keeps the first occurrence in the index when a duplicate is found", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		eventPathA := writeEvent(t, dir, "a.esdm.yaml", "order-placed")
		eventPathB := writeEvent(t, dir, "b.esdm.yaml", "order-placed")
		files := parseAll(t, parents, eventPathA, eventPathB)

		m, _ := resolver.Resolve(files)
		require.Len(t, m.Events, 1)

		ev, ok := m.LookupEvent("commerce", "ordering", "order", "order-placed")
		require.True(t, ok)
		assert.Equal(t, eventPathA, ev.Name().Location().File)
	})

	t.Run("handles multi-document files", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		path := filepath.Join(dir, "multi.esdm.yaml")
		content := fmt.Sprintf(eventTemplate, "order-placed") + "---\n" + fmt.Sprintf(eventTemplate, "order-shipped")
		writeFile(t, path, content)

		files := parseAll(t, parents, path)

		m, diagnostics := resolver.Resolve(files)
		assert.Empty(t, diagnostics)
		assert.Len(t, m.Events, 2)
	})

	t.Run("silently skips documents with an unknown kind", func(t *testing.T) {
		// Typos in the kind field are caught by the
		// parser's schema validation as
		// esdm/structure/constraint-violation; the
		// resolver does not emit a second, misleading
		// "linter bug" diagnostic on top of that.
		file := syntheticFile(t, "bogus.esdm.yaml", "apiVersion: schema.esdm.io/core/v1\nkind: not-a-real-kind\nname: bogus\n")

		_, diagnostics := resolver.Resolve([]*parser.ParsedFile{file})
		assert.Empty(t, diagnostics)
	})

	t.Run("emits unresolved-reference for an event scoped to a missing aggregate, with a did-you-mean hint", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		// scope.aggregate "ordr" instead of "order"
		eventPath := filepath.Join(dir, "event.esdm.yaml")
		writeFile(t, eventPath, `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-placed
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: ordr
data:
  type: object
`)
		files := parseAll(t, parents, eventPath)

		_, diagnostics := resolver.Resolve(files)

		var didFindUnresolvedReference bool
		for _, d := range diagnostics {
			if d.RuleID != "esdm/structure/unresolved-reference" {
				continue
			}
			didFindUnresolvedReference = true
			assert.Contains(t, d.Message, "ordr")
			require.Len(t, d.Related, 1)
			assert.Contains(t, d.Related[0].Message, "order")
		}
		assert.True(t, didFindUnresolvedReference, "expected unresolved-reference, got %+v", diagnostics)
	})

	t.Run("emits unresolved-reference for an unknown domain", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		eventPath := filepath.Join(dir, "event.esdm.yaml")
		writeFile(t, eventPath, `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-placed
scope:
  domain: nonexistent
  boundedContext: ordering
  aggregate: order
data:
  type: object
`)
		files := parseAll(t, parents, eventPath)

		_, diagnostics := resolver.Resolve(files)

		var didFindUnresolvedReference bool
		for _, d := range diagnostics {
			if d.RuleID == "esdm/structure/unresolved-reference" && strings.Contains(d.Message, "nonexistent") {
				didFindUnresolvedReference = true
				assert.Contains(t, d.Message, "domain")
			}
		}
		assert.True(t, didFindUnresolvedReference, "expected unresolved-reference for the missing domain, got %+v", diagnostics)
	})

	t.Run("emits unresolved-reference when the aggregate exists but lives in a different bounded context", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		// Define a second BC and an aggregate inside it; then
		// reference that aggregate but with the wrong BC.
		extras := filepath.Join(dir, "extras.esdm.yaml")
		writeFile(t, extras, `apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: shipping
scope:
  domain: commerce
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: shipment
scope:
  domain: commerce
  boundedContext: shipping
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
`)
		eventPath := filepath.Join(dir, "event.esdm.yaml")
		writeFile(t, eventPath, `apiVersion: schema.esdm.io/core/v1
kind: event
name: shipment-recorded
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: shipment
data:
  type: object
`)
		files := parseAll(t, parents, extras, eventPath)

		_, diagnostics := resolver.Resolve(files)

		var didFindUnresolvedReference bool
		for _, d := range diagnostics {
			if d.RuleID == "esdm/structure/unresolved-reference" && strings.Contains(d.Message, "shipment") && strings.Contains(d.Message, "exists") {
				didFindUnresolvedReference = true
			}
		}
		assert.True(t, didFindUnresolvedReference, "expected mismatched-parent diagnostic, got %+v", diagnostics)
	})

	t.Run("flags subdomain.boundedContexts entries that do not exist", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		subPath := filepath.Join(dir, "sub.esdm.yaml")
		writeFile(t, subPath, `apiVersion: schema.esdm.io/core/v1
kind: subdomain
name: ordering-domain
scope:
  domain: commerce
type: core
boundedContexts:
  - ordering
  - phantom
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, subPath))

		assertHasUnresolved(t, diagnostics, "phantom")
	})

	t.Run("flags subdomain.boundedContexts entries that live in another domain", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		// Add a second domain plus a bounded context in it.
		extras := filepath.Join(dir, "extras.esdm.yaml")
		writeFile(t, extras, `apiVersion: schema.esdm.io/core/v1
kind: domain
name: shipping
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: warehousing
scope:
  domain: shipping
`)
		subPath := filepath.Join(dir, "sub.esdm.yaml")
		writeFile(t, subPath, `apiVersion: schema.esdm.io/core/v1
kind: subdomain
name: ordering-domain
scope:
  domain: commerce
type: core
boundedContexts:
  - ordering
  - warehousing
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, extras, subPath))

		assertHasMismatch(t, diagnostics, "warehousing", "domain", "commerce")
	})

	t.Run("flags actor.backedBy entries that do not exist", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		actorPath := filepath.Join(dir, "actor.esdm.yaml")
		writeFile(t, actorPath, `apiVersion: schema.esdm.io/core/v1
kind: actor
name: payment-system
scope:
  domain: commerce
  boundedContext: ordering
type: system
backedBy:
  - phantom-provider
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, actorPath))

		assertHasUnresolved(t, diagnostics, "phantom-provider")
	})

	t.Run("flags query.readModel pointing at a non-existent read-model", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		queryPath := filepath.Join(dir, "query.esdm.yaml")
		writeFile(t, queryPath, `apiVersion: schema.esdm.io/core/v1
kind: query
name: find-orders
scope:
  domain: commerce
  boundedContext: ordering
readModel: phantom-projection
result:
  type: array
parameters:
  type: object
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, queryPath))

		assertHasUnresolved(t, diagnostics, "phantom-projection")
	})

	t.Run("flags command.publishes when the event lives in a different aggregate", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		// Second aggregate in the same BC, plus an event scoped to it.
		extras := filepath.Join(dir, "extras.esdm.yaml")
		writeFile(t, extras, `apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: invoice
scope:
  domain: commerce
  boundedContext: ordering
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: invoice-issued
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: invoice
data:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: actor
name: customer
scope:
  domain: commerce
  boundedContext: ordering
type: human
`)
		// Command in `order` aggregate that wrongly publishes invoice-issued.
		cmdPath := filepath.Join(dir, "cmd.esdm.yaml")
		writeFile(t, cmdPath, `apiVersion: schema.esdm.io/core/v1
kind: command
name: place-order
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
publishes:
  - invoice-issued
actors:
  - customer
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, extras, cmdPath))

		assertHasMismatch(t, diagnostics, "invoice-issued", "aggregate", "order")
	})

	t.Run("does not flag command.publishes for DCB-bound commands publishing cross-aggregate events", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		// DCB plus an event under a different aggregate in the same BC.
		extras := filepath.Join(dir, "extras.esdm.yaml")
		writeFile(t, extras, `apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: enrollment
scope:
  domain: commerce
  boundedContext: ordering
identifiedBy:
  - name: studentId
    source: command-payload
    field: studentId
consults:
  - boundedContext: ordering
    aggregate: order
    event: order-placed
    criteria: same studentId
---
apiVersion: schema.esdm.io/core/v1
kind: actor
name: customer
scope:
  domain: commerce
  boundedContext: ordering
type: human
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: order-placed
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
`)
		// The DCB-bound command publishes order-placed, which
		// lives in the `order` aggregate. Strict aggregate
		// matching would fail; lax DCB matching (BC only)
		// should accept it.
		cmdPath := filepath.Join(dir, "cmd.esdm.yaml")
		writeFile(t, cmdPath, `apiVersion: schema.esdm.io/core/v1
kind: command
name: enroll-student
scope:
  domain: commerce
  boundedContext: ordering
  dynamicConsistencyBoundary: enrollment
data:
  type: object
publishes:
  - order-placed
actors:
  - customer
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, extras, cmdPath))

		for _, d := range diagnostics {
			if d.RuleID != "esdm/structure/unresolved-reference" {
				continue
			}

			assert.NotContains(t, d.Message, "order-placed", "DCB-bound command should not produce a mismatched-parent diagnostic for cross-aggregate events")
		}
	})

	t.Run("flags event-handler.handles entries that reference a non-existent event", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		handlerPath := filepath.Join(dir, "handler.esdm.yaml")
		writeFile(t, handlerPath, `apiVersion: schema.esdm.io/core/v1
kind: event-handler
name: send-confirmation
scope:
  domain: commerce
deliveryGuarantee: at-most-once
handles:
  - boundedContext: ordering
    aggregate: order
    event: phantom-event
sideEffects:
  - type: other
    rule: log
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, handlerPath))

		assertHasUnresolved(t, diagnostics, "phantom-event")
	})

	t.Run("flags event-handler.handles entries whose aggregate does not match the event's scope", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		// order-placed lives under aggregate `order`;
		// the handler references it under aggregate `invoice`.
		extras := filepath.Join(dir, "extras.esdm.yaml")
		writeFile(t, extras, `apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: invoice
scope:
  domain: commerce
  boundedContext: ordering
identifiedBy:
  source: generated
  generator: uuid
state:
  type: object
---
apiVersion: schema.esdm.io/core/v1
kind: event
name: order-placed
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
`)
		handlerPath := filepath.Join(dir, "handler.esdm.yaml")
		writeFile(t, handlerPath, `apiVersion: schema.esdm.io/core/v1
kind: event-handler
name: send-confirmation
scope:
  domain: commerce
deliveryGuarantee: at-most-once
handles:
  - boundedContext: ordering
    aggregate: invoice
    event: order-placed
sideEffects:
  - type: other
    rule: log
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, extras, handlerPath))

		assertHasMismatch(t, diagnostics, "order-placed", "aggregate", "invoice")
	})

	t.Run("flags policy.emits entries that reference a non-existent command", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		extras := filepath.Join(dir, "extras.esdm.yaml")
		writeFile(t, extras, `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-placed
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
`)
		policyPath := filepath.Join(dir, "policy.esdm.yaml")
		writeFile(t, policyPath, `apiVersion: schema.esdm.io/core/v1
kind: policy
name: reserve-on-placement
scope:
  domain: commerce
deliveryGuarantee: at-most-once
handles:
  - boundedContext: ordering
    aggregate: order
    event: order-placed
emits:
  - boundedContext: ordering
    aggregate: order
    command: phantom-reserve
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, extras, policyPath))

		assertHasUnresolved(t, diagnostics, "phantom-reserve")
	})

	t.Run("does not flag process-manager.reactions when the trigger is a timer reference", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		extras := filepath.Join(dir, "extras.esdm.yaml")
		writeFile(t, extras, `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-placed
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
  properties:
    orderId:
      type: string
`)
		pmPath := filepath.Join(dir, "pm.esdm.yaml")
		writeFile(t, pmPath, `apiVersion: schema.esdm.io/core/v1
kind: process-manager
name: shipping-saga
scope:
  domain: commerce
  boundedContext: ordering
deliveryGuarantee: at-most-once
correlatedBy:
  name: orderId
  field: orderId
state:
  type: object
  properties:
    done:
      type: boolean
startsWhen:
  - boundedContext: ordering
    aggregate: order
    event: order-placed
endsWhen:
  - name: done
    condition: state.done is true
timers:
  - name: expiry
    after:
      value: 30
      unit: minutes
reactions:
  - when:
      timer: expiry
    rule: expire the saga
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, extras, pmPath))

		for _, d := range diagnostics {
			assert.NotEqual(t, "esdm/structure/unresolved-reference", d.RuleID, "timer reference should not trigger unresolved-reference: %+v", d)
		}
	})

	t.Run("indexes a domain-storytelling document under Extensions and flags a missing domain", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		storyPath := filepath.Join(dir, "story.esdm.yaml")
		writeFile(t, storyPath, `apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: place-order
scope:
  domain: phantom-domain
sentences:
  - sequenceNumber: 1
    edges:
      - from:
          actor: customer
        to:
          workObject: order
`)
		m, diagnostics := resolver.Resolve(parseAll(t, parents, storyPath))

		assert.Contains(t, m.Extensions.DomainStorytelling.Stories, "phantom-domain/place-order")
		assertHasUnresolved(t, diagnostics, "phantom-domain")
	})

	t.Run("indexes a given-when-then feature document under Extensions", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		featurePath := filepath.Join(dir, "feature.esdm.yaml")
		writeFile(t, featurePath, `apiVersion: schema.esdm.io/given-when-then/v1
kind: feature
name: order-cancellation
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
scenarios:
  - name: cancels-an-open-order
    given:
      - event: order-placed
        data: {}
    when:
      command: cancel-order
      data: {}
    then:
      events:
        - event: order-canceled
          data: {}
`)
		m, _ := resolver.Resolve(parseAll(t, parents, featurePath))

		assert.Contains(t, m.Extensions.GivenWhenThen.Features, "commerce/order-cancellation")
	})

	t.Run("does not conflate an extension kind with a core kind of the same name", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		storyPath := filepath.Join(dir, "story.esdm.yaml")
		// A valid domain-story named "order" (same name as
		// the aggregate from parentsYAML) must not produce
		// a duplicate-name diagnostic - the two live in
		// separate namespaces.
		writeFile(t, storyPath, `apiVersion: schema.esdm.io/domain-storytelling/v1
kind: domain-story
name: order
scope:
  domain: commerce
sentences:
  - sequenceNumber: 1
    edges:
      - from:
          actor: customer
        to:
          workObject: order
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, storyPath))

		for _, d := range diagnostics {
			assert.NotEqual(t, "esdm/structure/duplicate-name", d.RuleID,
				"extension kind and core kind should not collide: %+v", d)
		}
	})

	t.Run("flags context-mapping endpoints that reference a non-existent bounded-context", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		mappingPath := filepath.Join(dir, "mapping.esdm.yaml")
		writeFile(t, mappingPath, `apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: ordering-shipping
type: customer-supplier
customer:
  domain: commerce
  boundedContext: ordering
supplier:
  domain: commerce
  boundedContext: phantom-context
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, mappingPath))

		assertHasUnresolved(t, diagnostics, "phantom-context")
	})

	t.Run("flags context-mapping endpoints that reference an external system in the wrong domain", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		extras := filepath.Join(dir, "extras.esdm.yaml")
		writeFile(t, extras, `apiVersion: schema.esdm.io/core/v1
kind: domain
name: shipping
---
apiVersion: schema.esdm.io/core/v1
kind: external-system
name: stripe
scope:
  domain: shipping
direction: outbound
`)
		mappingPath := filepath.Join(dir, "mapping.esdm.yaml")
		writeFile(t, mappingPath, `apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: ordering-stripe
type: anti-corruption-layer
downstream:
  domain: commerce
  boundedContext: ordering
upstream:
  domain: commerce
  externalSystem: stripe
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, extras, mappingPath))

		assertHasMismatch(t, diagnostics, "stripe", "domain", "commerce")
	})

	t.Run("resolves context-mapping participants when both bounded-contexts exist", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		extras := filepath.Join(dir, "extras.esdm.yaml")
		writeFile(t, extras, `apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: billing
scope:
  domain: commerce
`)
		mappingPath := filepath.Join(dir, "mapping.esdm.yaml")
		writeFile(t, mappingPath, `apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: ordering-billing
type: shared-kernel
participants:
  - domain: commerce
    boundedContext: ordering
  - domain: commerce
    boundedContext: billing
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, extras, mappingPath))

		for _, d := range diagnostics {
			assert.NotEqual(t, "esdm/structure/unresolved-reference", d.RuleID, "valid participants should not produce unresolved-reference: %+v", d)
		}
	})

	t.Run("flags context-mapping participants that reference a non-existent bounded-context", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		mappingPath := filepath.Join(dir, "mapping.esdm.yaml")
		writeFile(t, mappingPath, `apiVersion: schema.esdm.io/core/v1
kind: context-mapping
name: ordering-phantom
type: partnership
participants:
  - domain: commerce
    boundedContext: ordering
  - domain: commerce
    boundedContext: phantom-context
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, mappingPath))

		assertHasUnresolved(t, diagnostics, "phantom-context")
	})

	t.Run("flags dcb.consults entries that reference a non-existent event", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		dcbPath := filepath.Join(dir, "dcb.esdm.yaml")
		writeFile(t, dcbPath, `apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: enrollment
scope:
  domain: commerce
  boundedContext: ordering
identifiedBy:
  - name: studentId
    source: command-payload
    field: studentId
consults:
  - boundedContext: ordering
    aggregate: order
    event: phantom-event
    criteria: same studentId
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, dcbPath))

		assertHasUnresolved(t, diagnostics, "phantom-event")
	})

	t.Run("flags aggregate.identifiedBy.field when the named field is not declared in state.properties", func(t *testing.T) {
		dir := t.TempDir()
		// Redefine the domain + bc but skip the default
		// aggregate so we can supply our own with a bad
		// identifiedBy reference.
		own := filepath.Join(dir, "own.esdm.yaml")
		writeFile(t, own, `apiVersion: schema.esdm.io/core/v1
kind: domain
name: commerce
---
apiVersion: schema.esdm.io/core/v1
kind: bounded-context
name: ordering
scope:
  domain: commerce
---
apiVersion: schema.esdm.io/core/v1
kind: aggregate
name: order
scope:
  domain: commerce
  boundedContext: ordering
identifiedBy:
  source: state
  field: orderld
state:
  type: object
  properties:
    orderId:
      type: string
`)
		_, diagnostics := resolver.Resolve(parseAll(t, own))

		assertHasUnresolved(t, diagnostics, "orderld")
	})

	t.Run("flags process-manager.timers[].at when the named field is not in state.properties", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		extras := filepath.Join(dir, "extras.esdm.yaml")
		writeFile(t, extras, `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-placed
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
  properties:
    orderId:
      type: string
`)
		pmPath := filepath.Join(dir, "pm.esdm.yaml")
		writeFile(t, pmPath, `apiVersion: schema.esdm.io/core/v1
kind: process-manager
name: shipping-saga
scope:
  domain: commerce
  boundedContext: ordering
deliveryGuarantee: at-most-once
correlatedBy:
  name: orderId
  field: orderId
state:
  type: object
  properties:
    expiresAt:
      type: string
startsWhen:
  - boundedContext: ordering
    aggregate: order
    event: order-placed
endsWhen:
  - name: done
    condition: state.done is true
timers:
  - name: expiry
    at: expirsAt
reactions:
  - when:
      timer: expiry
    rule: expire the saga
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, extras, pmPath))

		assertHasUnresolved(t, diagnostics, "expirsAt")
	})

	t.Run("flags process-manager.correlatedBy.field when absent from a referenced event's data", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		extras := filepath.Join(dir, "extras.esdm.yaml")
		writeFile(t, extras, `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-placed
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
  properties:
    orderIdentifier:
      type: string
`)
		pmPath := filepath.Join(dir, "pm.esdm.yaml")
		writeFile(t, pmPath, `apiVersion: schema.esdm.io/core/v1
kind: process-manager
name: shipping-saga
scope:
  domain: commerce
  boundedContext: ordering
deliveryGuarantee: at-most-once
correlatedBy:
  name: orderId
  field: orderId
state:
  type: object
  properties:
    done:
      type: boolean
startsWhen:
  - boundedContext: ordering
    aggregate: order
    event: order-placed
endsWhen:
  - name: done
    condition: state.done is true
reactions:
  - when:
      boundedContext: ordering
      aggregate: order
      event: order-placed
    rule: mark placed
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, extras, pmPath))

		assertHasUnresolved(t, diagnostics, "orderId")
		didFindCorrelationDiagnostic := false
		for _, d := range diagnostics {
			if strings.Contains(d.Message, "order-placed") && strings.Contains(d.Message, "orderId") {
				didFindCorrelationDiagnostic = true
				// Related should suggest "orderIdentifier".
				if len(d.Related) > 0 {
					assert.Contains(t, d.Related[0].Message, "orderIdentifier")
				}
				break
			}
		}
		assert.True(t, didFindCorrelationDiagnostic, "expected correlation diagnostic mentioning event and field, got %+v", diagnostics)
	})

	t.Run("flags dcb.identifiedBy command-payload field when absent from a targeting command's data", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		dcbPath := filepath.Join(dir, "dcb.esdm.yaml")
		writeFile(t, dcbPath, `apiVersion: schema.esdm.io/core/v1
kind: dynamic-consistency-boundary
name: enrollment
scope:
  domain: commerce
  boundedContext: ordering
identifiedBy:
  - name: studentId
    source: command-payload
    field: studentId
consults: []
`)
		cmdPath := filepath.Join(dir, "cmd.esdm.yaml")
		writeFile(t, cmdPath, `apiVersion: schema.esdm.io/core/v1
kind: command
name: enroll-student
scope:
  domain: commerce
  boundedContext: ordering
  dynamicConsistencyBoundary: enrollment
data:
  type: object
  properties:
    learnerId:
      type: string
publishes:
  - student-enrolled
actors:
  - customer
`)
		actorPath := filepath.Join(dir, "actor.esdm.yaml")
		writeFile(t, actorPath, `apiVersion: schema.esdm.io/core/v1
kind: actor
name: customer
scope:
  domain: commerce
  boundedContext: ordering
type: human
`)
		eventPath := filepath.Join(dir, "event.esdm.yaml")
		writeFile(t, eventPath, `apiVersion: schema.esdm.io/core/v1
kind: event
name: student-enrolled
scope:
  domain: commerce
  boundedContext: ordering
data:
  type: object
`)
		_, diagnostics := resolver.Resolve(parseAll(t, parents, dcbPath, cmdPath, actorPath, eventPath))

		didFindIdentifierDiagnostic := false
		for _, d := range diagnostics {
			if strings.Contains(d.Message, "studentId") && strings.Contains(d.Message, "enroll-student") {
				didFindIdentifierDiagnostic = true
				if len(d.Related) > 0 {
					assert.Contains(t, d.Related[0].Message, "learnerId")
				}
				break
			}
		}
		assert.True(t, didFindIdentifierDiagnostic, "expected DCB identifiedBy diagnostic mentioning field and command, got %+v", diagnostics)
	})

	t.Run("flags command.actors entries that point to an unknown actor", func(t *testing.T) {
		dir := t.TempDir()
		parents := writeParents(t, dir)
		cmdPath := filepath.Join(dir, "cmd.esdm.yaml")
		writeFile(t, cmdPath, `apiVersion: schema.esdm.io/core/v1
kind: command
name: place-order
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
publishes:
  - order-placed
actors:
  - phantom-actor
`)
		eventPath := filepath.Join(dir, "event.esdm.yaml")
		writeFile(t, eventPath, validEventForReferenceTests())

		_, diagnostics := resolver.Resolve(parseAll(t, parents, cmdPath, eventPath))

		assertHasUnresolved(t, diagnostics, "phantom-actor")
	})
}

// validEventForReferenceTests returns a valid event in
// `order` aggregate so command.publishes can resolve when
// only command.actors is the focus of the test.
func validEventForReferenceTests() string {
	return `apiVersion: schema.esdm.io/core/v1
kind: event
name: order-placed
scope:
  domain: commerce
  boundedContext: ordering
  aggregate: order
data:
  type: object
`
}

func assertHasUnresolved(t *testing.T, diagnostics []diag.Diagnostic, name string) {
	t.Helper()
	for _, d := range diagnostics {
		if d.RuleID == "esdm/structure/unresolved-reference" && strings.Contains(d.Message, name) {
			return
		}
	}
	require.Failf(t, "unresolved-reference diagnostic missing", "expected a diagnostic mentioning %q, got %+v", name, diagnostics)
}

func assertHasMismatch(t *testing.T, diagnostics []diag.Diagnostic, name, parentKind, expectedParent string) {
	t.Helper()
	for _, d := range diagnostics {
		if d.RuleID != "esdm/structure/unresolved-reference" {
			continue
		}
		if !strings.Contains(d.Message, name) {
			continue
		}
		if !strings.Contains(d.Message, parentKind) {
			continue
		}
		if !strings.Contains(d.Message, expectedParent) {
			continue
		}
		return
	}
	require.Failf(t, "mismatch diagnostic missing", "expected a diagnostic mentioning %q with %s %q, got %+v", name, parentKind, expectedParent, diagnostics)
}
