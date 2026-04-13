package handlers

import (
	"fmt"

	"wsan/handlers/generator"
)

// WsanField describes a single URL-parameter slot inside an operation's route.
// Name is a human-readable label; Field is the chi URL param key (e.g. "name"
// for a route like "/greet/{name}").
type WsanField struct {
	Name  string `json:"name"`
	Field string `json:"field"`
}

// WsanOperation is a single WSAN endpoint registered at init-time by an
// ops_*.go file. Render is invoked per request with the extracted URL params
// and must return the (message, subtitle) pair that gets JSON-encoded.
type WsanOperation struct {
	Name   string      `json:"name"`
	URL    string      `json:"url"`
	Fields []WsanField `json:"fields"`
	// Render is invoked per request. It receives the per-request generator
	// (either the package default or a seeded one when ?seed= is supplied)
	// and the URL params, and returns the (message, subtitle) pair.
	Render func(g *generator.Generator, params map[string]string) (message, subtitle string) `json:"-"`
}

var wsanOps []WsanOperation

// RegisterOp appends op to the package-level registry. It panics if op.URL
// collides with an already-registered operation so duplicates are caught at
// program start rather than silently shadowing each other.
//
// RegisterOp is intended to be called only from package init() functions and
// is NOT safe for concurrent use at runtime. No mutex is taken because the Go
// runtime serializes init functions within a package.
func RegisterOp(op WsanOperation) {
	for _, existing := range wsanOps {
		if existing.URL == op.URL {
			panic(fmt.Sprintf("wsan: duplicate operation URL %q (existing=%q, new=%q)", op.URL, existing.Name, op.Name))
		}
	}
	wsanOps = append(wsanOps, op)
}

// AllOps returns a defensive copy of the registered operations so callers
// cannot mutate the underlying slice.
func AllOps() []WsanOperation {
	out := make([]WsanOperation, len(wsanOps))
	copy(out, wsanOps)
	return out
}
