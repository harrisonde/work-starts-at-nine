package handlers

import (
	"testing"

	"wsan/handlers/generator"
)

func TestRegistryHasAllOps(t *testing.T) {
	expected := []string{
		"/nine/{name}/{from}",
		"/late/{name}/{from}",
		"/reminder/{name}/{from}",
		"/strict/{name}/{from}",
		"/boss/{name}/{from}",
		"/hr/{name}/{from}",
		"/polite/{name}/{from}",
		"/rude/{name}/{from}",
		"/monday/{name}/{from}",
		"/meeting/{name}/{from}",
		"/earlybird/{name}/{from}",
		"/coffee/{name}/{from}",
		"/standup/{name}/{from}",
		"/deadline/{name}/{from}",
		"/random/{name}/{from}",
	}
	ops := AllOps()
	if len(ops) != 15 {
		t.Errorf("expected 15 ops, got %d", len(ops))
	}
	expectedSet := make(map[string]bool, len(expected))
	for _, u := range expected {
		expectedSet[u] = true
	}
	registeredSet := make(map[string]bool, len(ops))
	for _, o := range ops {
		registeredSet[o.URL] = true
	}
	for u := range expectedSet {
		if !registeredSet[u] {
			t.Errorf("expected URL not registered: %q", u)
		}
	}
	for u := range registeredSet {
		if !expectedSet[u] {
			t.Errorf("unexpected URL registered: %q", u)
		}
	}
}

func TestRegistryNoDuplicateURLs(t *testing.T) {
	counts := make(map[string]int)
	for _, o := range AllOps() {
		counts[o.URL]++
	}
	for u, c := range counts {
		if c > 1 {
			t.Errorf("duplicate URL %q registered %d times", u, c)
		}
	}
}

func TestRegistryAllHaveRenders(t *testing.T) {
	g := generator.New(1)
	for _, o := range AllOps() {
		if o.Render == nil {
			t.Errorf("op %q has nil Render", o.URL)
			continue
		}
		if o.Name == "" {
			t.Errorf("op with URL %q has empty Name", o.URL)
		}
		if o.URL == "" {
			t.Errorf("op %q has empty URL", o.Name)
		}
		// Smoke check: every Render must work with a real generator and
		// produce a non-empty (message, subtitle) pair.
		msg, sub := o.Render(g, map[string]string{"name": "Alice", "from": "Bob"})
		if msg == "" || sub == "" {
			t.Errorf("op %q returned empty render: msg=%q sub=%q", o.URL, msg, sub)
		}
	}
}
