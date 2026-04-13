package handlers

import (
	"strings"
	"testing"

	"wsan/handlers/generator"
)

// findOp locates a registered op by URL or fails the test. Shared by every
// per-op test that wants to assert against AllOps().
func findOp(t *testing.T, url string) WsanOperation {
	t.Helper()
	for _, op := range AllOps() {
		if op.URL == url {
			return op
		}
	}
	t.Fatalf("op with URL %q not registered", url)
	return WsanOperation{}
}

// TestAllOpsRender exercises every registered op end-to-end via its Render
// func. We use a fixed seed so failures are reproducible, and assert SOFT
// properties (non-empty, contains name, contains a 9/nine reference) rather
// than exact strings — fragment edits should not break this test.
func TestAllOpsRender(t *testing.T) {
	cases := []struct {
		url string
	}{
		{"/nine/{name}/{from}"},
		{"/late/{name}/{from}"},
		{"/reminder/{name}/{from}"},
		{"/strict/{name}/{from}"},
		{"/boss/{name}/{from}"},
		{"/hr/{name}/{from}"},
		{"/polite/{name}/{from}"},
		{"/rude/{name}/{from}"},
		{"/monday/{name}/{from}"},
		{"/meeting/{name}/{from}"},
		{"/earlybird/{name}/{from}"},
		{"/coffee/{name}/{from}"},
		{"/standup/{name}/{from}"},
		{"/deadline/{name}/{from}"},
		{"/random/{name}/{from}"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.url, func(t *testing.T) {
			op := findOp(t, tc.url)
			g := generator.New(1)
			msg, sub := op.Render(g, map[string]string{"name": "Alice", "from": "Bob"})
			if msg == "" || sub == "" {
				t.Fatalf("empty render: msg=%q sub=%q", msg, sub)
			}
			if !strings.Contains(msg, "Alice") {
				t.Errorf("msg %q does not contain name", msg)
			}
			if !strings.Contains(sub, "Bob") {
				t.Errorf("sub %q does not contain from", sub)
			}
			if !generator.ContainsNineReference(msg) {
				t.Errorf("msg %q missing 9/nine reference", msg)
			}
		})
	}
}

// TestAllOpsRenderSeededDeterministic verifies that two generators built
// with the same seed produce identical Render output for the same op. This
// is the underlying guarantee that powers ?seed= in WsanServeOp.
func TestAllOpsRenderSeededDeterministic(t *testing.T) {
	for _, op := range AllOps() {
		op := op
		t.Run(op.URL, func(t *testing.T) {
			g1 := generator.New(42)
			g2 := generator.New(42)
			m1, s1 := op.Render(g1, map[string]string{"name": "Alice", "from": "Bob"})
			m2, s2 := op.Render(g2, map[string]string{"name": "Alice", "from": "Bob"})
			if m1 != m2 || s1 != s2 {
				t.Errorf("seeded output differs: (%q,%q) vs (%q,%q)", m1, s1, m2, s2)
			}
		})
	}
}
