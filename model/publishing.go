package model

import "sort"

// CommandMayPublish reports whether cmd is allowed to publish
// event, judged by their scopes alone. Both have to sit in
// the same domain and bounded context. An aggregate-bound
// command may only publish events owned by its own
// aggregate; a DCB-bound command may publish any event of
// its bounded context, free-standing or aggregate-owned,
// because a dynamic consistency boundary deliberately spans
// aggregates. The event's position is read from the event's
// own scope, never inferred from the command's.
//
// This is the one definition of "this command publishes
// this event". The resolver accepts references with it, the
// event-without-publisher rule and the view's publisher
// annotation derive from it, so the three cannot drift
// apart.
func CommandMayPublish(cmd CommandView, event EventView) bool {
	commandScope := cmd.Scope()
	eventScope := event.Scope()

	if ScopeText(commandScope, "domain") != ScopeText(eventScope, "domain") {
		return false
	}
	if ScopeText(commandScope, "boundedContext") != ScopeText(eventScope, "boundedContext") {
		return false
	}

	aggregate := ScopeText(commandScope, "aggregate")
	if aggregate == "" {
		return true
	}
	return aggregate == ScopeText(eventScope, "aggregate")
}

// PublishersOf returns every command that lists event in its
// publishes and is allowed to publish it, sorted by name.
func (m *Model) PublishersOf(event EventView) []CommandView {
	eventName, ok := event.Name().Text()
	if !ok {
		return nil
	}

	var out []CommandView
	for _, cmd := range m.Commands {
		if !commandListsEvent(cmd, eventName) {
			continue
		}
		if !CommandMayPublish(cmd, event) {
			continue
		}
		out = append(out, cmd)
	}
	sort.Slice(out, func(i, j int) bool {
		ni, _ := out[i].Name().Text()
		nj, _ := out[j].Name().Text()
		return ni < nj
	})
	return out
}

// PublishedEvent returns the event called name that cmd may
// publish, if the model holds one. It answers the resolver's
// question for a single publishes entry: does the bare name
// resolve to an event within the command's reach?
func (m *Model) PublishedEvent(cmd CommandView, name string) (EventView, bool) {
	for _, event := range m.Events {
		eventName, _ := event.Name().Text()
		if eventName != name {
			continue
		}
		if CommandMayPublish(cmd, event) {
			return event, true
		}
	}
	return EventView{}, false
}

// commandListsEvent reports whether name appears in
// cmd.publishes.
func commandListsEvent(cmd CommandView, name string) bool {
	for _, item := range cmd.Publishes().Seq() {
		listed, ok := item.Text()
		if ok && listed == name {
			return true
		}
	}
	return false
}
