package handlers

import "wsan/handlers/generator"

// /random picks a tone uniformly at random per request and delegates to the
// shared generator. With ?seed=<int>, both the tone choice and the fragment
// selection are reproducible because they share the same RNG.
func init() {
	RegisterOp(WsanOperation{
		Name: "Random",
		URL:  "/random/{name}/{from}",
		Fields: []WsanField{
			{Name: "Name", Field: "name"},
			{Name: "From", Field: "from"},
		},
		Render: func(g *generator.Generator, p map[string]string) (string, string) {
			return g.Compose(g.RandomTone(), p["name"], p["from"])
		},
	})
}
