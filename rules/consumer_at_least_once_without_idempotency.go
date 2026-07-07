package rules

import (
	"context"
	"fmt"

	"github.com/thenativeweb/esdm/ast"
	"github.com/thenativeweb/esdm/diag"
	"github.com/thenativeweb/esdm/model"
)

const ruleIDConsumerAtLeastOnceWithoutIdempotency = "esdm/modeling/consumer-at-least-once-without-idempotency"

type consumerAtLeastOnceWithoutIdempotencyRule struct{}

func newConsumerAtLeastOnceWithoutIdempotencyRule() *consumerAtLeastOnceWithoutIdempotencyRule {
	return &consumerAtLeastOnceWithoutIdempotencyRule{}
}

func (*consumerAtLeastOnceWithoutIdempotencyRule) Meta() Meta {
	return Meta{
		ID:          ruleIDConsumerAtLeastOnceWithoutIdempotency,
		Severity:    diag.SeverityWarning,
		Description: "An event-handler, policy or process-manager declared with deliveryGuarantee: at-least-once must declare an idempotency strategy; mirrors the JSON Schema's allOf if-then constraint as defense in depth.",
	}
}

func (r *consumerAtLeastOnceWithoutIdempotencyRule) Check(ctx context.Context, m *model.Model, report diag.Reporter) {
	for _, eh := range sortedByName(m.EventHandlers) {
		name, _ := eh.Name().Text()
		r.checkConsumer("event-handler", name, eh.DeliveryGuarantee(), eh.Idempotency(), eh.Name().Location(), report)
	}
	for _, p := range sortedByName(m.Policies) {
		name, _ := p.Name().Text()
		r.checkConsumer("policy", name, p.DeliveryGuarantee(), p.Idempotency(), p.Name().Location(), report)
	}
	for _, pm := range sortedByName(m.ProcessManagers) {
		name, _ := pm.Name().Text()
		r.checkConsumer("process-manager", name, pm.DeliveryGuarantee(), pm.Idempotency(), pm.Name().Location(), report)
	}
}

func (*consumerAtLeastOnceWithoutIdempotencyRule) checkConsumer(kind, name string, deliveryGuarantee, idempotency ast.Node, loc diag.Location, report diag.Reporter) {
	guarantee, _ := deliveryGuarantee.Text()
	if guarantee != "at-least-once" {
		return
	}
	if idempotency.Exists() {
		return
	}
	report.Report(diag.Diagnostic{
		Message:  fmt.Sprintf("%s %q declares deliveryGuarantee: at-least-once but has no idempotency strategy", kind, name),
		Location: loc,
	})
}
