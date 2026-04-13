package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wsan/handlers/generator"

	"github.com/go-chi/chi/v5"
)

// wsanMessage is the envelope returned by WSAN JSON endpoints.
type wsanMessage struct {
	Message  string `json:"message"`
	Subtitle string `json:"subtitle"`
}

// writeJSON prefers the framework helper when the handler is wired to a real
// *adele.Adele; in tests (where App or Helpers is nil) it falls back to the
// stdlib encoder so tests can construct a zero-value Handlers.
func (h *Handlers) writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	if h != nil && h.App != nil && h.App.Helpers != nil {
		return h.App.Helpers.WriteJSON(w, status, data)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// WsanRoot returns the service banner.
func (h *Handlers) WsanRoot(w http.ResponseWriter, r *http.Request) {
	if err := h.writeJSON(w, http.StatusOK, wsanMessage{
		Message:  "WSAN v0.1.0 — Work Starts At Nine",
		Subtitle: "Reminding people since 9:00 AM",
	}); err != nil && h.App != nil && h.App.Log != nil {
		h.App.Log.Error("wsan writeJSON:", err)
	}
}

// WsanNotFound returns a 404 JSON envelope for unmatched /api routes.
func (h *Handlers) WsanNotFound(w http.ResponseWriter, r *http.Request) {
	if err := h.writeJSON(w, http.StatusNotFound, wsanMessage{
		Message:  "Not found. But work still starts at 9.",
		Subtitle: "— WSAN",
	}); err != nil && h.App != nil && h.App.Log != nil {
		h.App.Log.Error("wsan writeJSON:", err)
	}
}

// WsanOperations lists every operation registered via RegisterOp.
func (h *Handlers) WsanOperations(w http.ResponseWriter, r *http.Request) {
	if err := h.writeJSON(w, http.StatusOK, AllOps()); err != nil && h.App != nil && h.App.Log != nil {
		h.App.Log.Error("wsan writeJSON:", err)
	}
}

// WsanServeOp returns an http.HandlerFunc bound to a specific operation. It
// extracts each declared WsanField from the chi URL params, calls op.Render,
// and JSON-encodes the resulting message/subtitle pair.
func (h *Handlers) WsanServeOp(op WsanOperation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := make(map[string]string, len(op.Fields))
		for _, f := range op.Fields {
			params[f.Field] = chi.URLParam(r, f.Field)
		}
		// Seed semantics:
		//   - absent   -> use the package default (time-seeded) generator
		//   - present, parses as non-zero int64 -> per-request seeded generator
		//   - present but invalid or zero       -> HTTP 400 with WSAN envelope
		// seed=0 is reserved because generator.New treats 0 as "time-seeded".
		var g *generator.Generator
		if raw, ok := r.URL.Query()["seed"]; ok && len(raw) > 0 {
			seed, err := strconv.ParseInt(raw[0], 10, 64)
			if err != nil || seed == 0 {
				if werr := h.writeJSON(w, http.StatusBadRequest, wsanMessage{
					Message:  "Invalid seed. Use a non-zero int64.",
					Subtitle: "— WSAN",
				}); werr != nil && h.App != nil && h.App.Log != nil {
					h.App.Log.Error("wsan writeJSON:", werr)
				}
				return
			}
			g = generator.New(seed)
		} else {
			g = getDefaultGenerator()
		}
		msg, sub := op.Render(g, params)
		if err := h.writeJSON(w, http.StatusOK, wsanMessage{Message: msg, Subtitle: sub}); err != nil && h.App != nil && h.App.Log != nil {
			h.App.Log.Error("wsan writeJSON:", err)
		}
	}
}
