package handlers

import (
	"strings"

	"wsan/handlers/generator"
)

// toneDisplayName overrides the default title-cased display name for tones
// whose canonical presentation differs from strings.Title.
var toneDisplayName = map[string]string{
	"earlybird": "EarlyBird",
	"hr":        "HR",
}

// displayNameForTone returns the human-readable Name field for a tone key.
// Unknown keys get a simple title-casing (first rune upper, rest unchanged).
func displayNameForTone(tone string) string {
	if n, ok := toneDisplayName[tone]; ok {
		return n
	}
	if tone == "" {
		return tone
	}
	return strings.ToUpper(tone[:1]) + tone[1:]
}

// init registers one WsanOperation per tone returned by generator.AllToneNames
// so that adding a new tone to the generator automatically produces a matching
// /{tone}/{name}/{from} endpoint. The /random op is registered separately in
// ops_random.go since its Render delegates to g.RandomTone.
func init() {
	for _, tone := range generator.AllToneNames() {
		tone := tone // capture
		RegisterOp(WsanOperation{
			Name: displayNameForTone(tone),
			URL:  "/" + tone + "/{name}/{from}",
			Fields: []WsanField{
				{Name: "Name", Field: "name"},
				{Name: "From", Field: "from"},
			},
			Render: func(g *generator.Generator, p map[string]string) (string, string) {
				return g.Compose(tone, p["name"], p["from"])
			},
		})
	}
}
