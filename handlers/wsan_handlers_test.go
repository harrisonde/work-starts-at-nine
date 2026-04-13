package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"wsan/handlers/generator"

	"github.com/go-chi/chi/v5"
)

func TestWsanRoot(t *testing.T) {
	h := &Handlers{}
	r := chi.NewRouter()
	r.Get("/", h.WsanRoot)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("content-type = %q, want application/json", ct)
	}
	if !strings.Contains(rr.Body.String(), "Work Starts At Nine") {
		t.Fatalf("body missing banner: %s", rr.Body.String())
	}
}

func TestWsanNotFound(t *testing.T) {
	h := &Handlers{}
	mux := chi.NewRouter()
	mux.NotFound(h.WsanNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/nope", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("content-type = %q, want application/json", ct)
	}
	var body struct {
		Message  string `json:"message"`
		Subtitle string `json:"subtitle"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v (%s)", err, rr.Body.String())
	}
	if body.Message != "Not found. But work still starts at 9." {
		t.Errorf("message = %q", body.Message)
	}
	if body.Subtitle != "— WSAN" {
		t.Errorf("subtitle = %q", body.Subtitle)
	}
}

func TestWsanOperations(t *testing.T) {
	h := &Handlers{}
	r := chi.NewRouter()
	r.Get("/operations", h.WsanOperations)

	req := httptest.NewRequest(http.MethodGet, "/operations", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var out []WsanOperation
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatalf("body is not a JSON array: %v (%s)", err, rr.Body.String())
	}
	// WSAN ships exactly 14 tone ops + /random = 15.
	if len(out) != 15 {
		t.Errorf("expected exactly 15 ops in response, got %d", len(out))
	}
	// Verify JSON shape includes the expected keys.
	var raw []map[string]json.RawMessage
	if err := json.Unmarshal(rr.Body.Bytes(), &raw); err != nil {
		t.Fatalf("body is not a JSON array of objects: %v", err)
	}
	for i, obj := range raw {
		for _, key := range []string{"name", "url", "fields"} {
			if _, ok := obj[key]; !ok {
				t.Errorf("op[%d] missing key %q", i, key)
			}
		}
	}
}

func TestWsanServeOpEdgeCases(t *testing.T) {
	h := &Handlers{}
	op := WsanOperation{
		Name: "fake",
		URL:  "/edge/{name}",
		Fields: []WsanField{
			{Name: "Name", Field: "name"},
		},
		Render: func(_ *generator.Generator, params map[string]string) (string, string) {
			return "hi [" + params["name"] + "]", "— sub"
		},
	}
	r := chi.NewRouter()
	r.Get(op.URL, h.WsanServeOp(op))

	cases := []struct {
		label    string
		path     string
		wantName string
	}{
		{"unicode", "/edge/Jos%C3%A9", "José"},
		{"urlencoded-space", "/edge/Alice%20Smith", "Alice Smith"},
	}
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", rr.Code)
			}
			var body struct {
				Message  string `json:"message"`
				Subtitle string `json:"subtitle"`
			}
			if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode: %v (%s)", err, rr.Body.String())
			}
			if !strings.Contains(body.Message, tc.wantName) {
				t.Errorf("message = %q, want to contain %q", body.Message, tc.wantName)
			}
		})
	}

	t.Run("empty-name", func(t *testing.T) {
		// chi does not match /edge/ against /edge/{name}; expect a 404.
		req := httptest.NewRequest(http.MethodGet, "/edge/", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404 for empty name segment", rr.Code)
		}
	})
}

func TestWsanServeOp(t *testing.T) {
	h := &Handlers{}
	op := WsanOperation{
		Name: "fake",
		URL:  "/fake/{name}/{from}",
		Fields: []WsanField{
			{Name: "Name", Field: "name"},
			{Name: "From", Field: "from"},
		},
		Render: func(_ *generator.Generator, params map[string]string) (string, string) {
			return "hi " + params["name"], "— " + params["from"]
		},
	}

	r := chi.NewRouter()
	r.Get(op.URL, h.WsanServeOp(op))

	req := httptest.NewRequest(http.MethodGet, "/fake/Alice/Bob", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var body struct {
		Message  string `json:"message"`
		Subtitle string `json:"subtitle"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v (%s)", err, rr.Body.String())
	}
	if body.Message != "hi Alice" {
		t.Errorf("message = %q, want %q", body.Message, "hi Alice")
	}
	if body.Subtitle != "— Bob" {
		t.Errorf("subtitle = %q, want %q", body.Subtitle, "— Bob")
	}
}

// TestWsanServeOpSeedIsDeterministic hits every registered op twice with the
// same ?seed= and asserts the response bodies are byte-for-byte equal. This
// is the contract that justifies advertising ?seed= in the README.
func TestWsanServeOpSeedIsDeterministic(t *testing.T) {
	h := &Handlers{}
	for _, op := range AllOps() {
		op := op
		t.Run(op.URL, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get(op.URL, h.WsanServeOp(op))

			// Build a concrete URL by replacing {field} with sample values.
			path := op.URL
			for _, f := range op.Fields {
				path = strings.Replace(path, "{"+f.Field+"}", "Alice", 1)
			}
			path += "?seed=42"

			hit := func() string {
				req := httptest.NewRequest(http.MethodGet, path, nil)
				rr := httptest.NewRecorder()
				r.ServeHTTP(rr, req)
				if rr.Code != http.StatusOK {
					t.Fatalf("status = %d, want 200", rr.Code)
				}
				return rr.Body.String()
			}
			first := hit()
			second := hit()
			if first != second {
				t.Errorf("seeded responses differ:\nfirst:  %s\nsecond: %s", first, second)
			}
		})
	}
}

// TestWsanServeOpVariesWithoutSeed hits /nine/{name}/{from} without a seed
// 20 times and asserts we observe at least 2 distinct messages. With 15-20
// fragments per slot the false-failure probability is negligible.
func TestWsanServeOpVariesWithoutSeed(t *testing.T) {
	h := &Handlers{}
	op := findOp(t, "/nine/{name}/{from}")
	r := chi.NewRouter()
	r.Get(op.URL, h.WsanServeOp(op))

	seen := make(map[string]struct{}, 4)
	for i := 0; i < 20; i++ {
		req := httptest.NewRequest(http.MethodGet, "/nine/Alice/Bob", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
		seen[rr.Body.String()] = struct{}{}
	}
	if len(seen) < 2 {
		t.Errorf("unseeded op returned only %d distinct body across 20 calls", len(seen))
	}
}

// TestWsanServeOpInvalidSeed asserts that ?seed= with a bogus or zero value
// returns 400 with the WSAN envelope instead of silently falling back.
func TestWsanServeOpInvalidSeed(t *testing.T) {
	h := &Handlers{}
	op := findOp(t, "/nine/{name}/{from}")
	r := chi.NewRouter()
	r.Get(op.URL, h.WsanServeOp(op))

	for _, raw := range []string{"notanint", "0", ""} {
		raw := raw
		t.Run("seed="+raw, func(t *testing.T) {
			path := "/nine/Alice/Bob?seed=" + raw
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", rr.Code)
			}
			var body struct {
				Message  string `json:"message"`
				Subtitle string `json:"subtitle"`
			}
			if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if body.Message != "Invalid seed. Use a non-zero int64." {
				t.Errorf("message = %q", body.Message)
			}
			if body.Subtitle != "— WSAN" {
				t.Errorf("subtitle = %q", body.Subtitle)
			}
		})
	}
}

func TestWsanRegisterOpDuplicatePanics(t *testing.T) {
	// Snapshot/restore the package registry so this test doesn't leak a
	// synthetic op into sibling tests.
	old := wsanOps
	t.Cleanup(func() { wsanOps = old })

	url := "/wsan_test_dup/{x}"
	op := WsanOperation{
		Name:   "dup-test",
		URL:    url,
		Fields: []WsanField{{Name: "X", Field: "x"}},
		Render: func(_ *generator.Generator, _ map[string]string) (string, string) { return "", "" },
	}
	RegisterOp(op)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on duplicate RegisterOp, got none")
		}
	}()
	RegisterOp(op)
}
