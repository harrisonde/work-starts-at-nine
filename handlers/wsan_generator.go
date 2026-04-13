package handlers

import (
	"sync"

	"wsan/handlers/generator"
)

// defaultGenerator is the package-level *generator.Generator used by every
// op.Render call when no per-request seed is supplied. It is initialized
// once at package init and accessed via getDefaultGenerator.
var (
	defaultGeneratorMu sync.RWMutex
	defaultGenerator   *generator.Generator
)

func init() {
	defaultGenerator = generator.New(0) // 0 == time-seeded
}

// getDefaultGenerator returns the current package-level generator. It is
// safe for concurrent use.
func getDefaultGenerator() *generator.Generator {
	defaultGeneratorMu.RLock()
	defer defaultGeneratorMu.RUnlock()
	return defaultGenerator
}
