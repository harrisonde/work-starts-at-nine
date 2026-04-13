package main

import (
	"net/http"

	"wsan/handlers"

	"github.com/cidekar/adele-framework/mux"
)

func (a *application) ApiRoutes() http.Handler {

	r := mux.NewRouter()

	// WSAN API routes. This router is mounted at /api in routes.go, so
	// patterns here are relative (e.g. "/" becomes "/api/").
	r.Get("/", a.Handlers.WsanRoot)
	r.Get("/operations", a.Handlers.WsanOperations)
	r.NotFound(a.Handlers.WsanNotFound)

	// Each operation registers itself via init() in handlers/ops_*.go.
	for _, op := range handlers.AllOps() {
		r.Get(op.URL, a.Handlers.WsanServeOp(op))
	}

	return r
}
