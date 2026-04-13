package generator

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
	"sync"
	"time"
)

// Generator composes WSAN messages from the per-tone fragment pools. A zero
// Generator is NOT usable; always construct one with New.
//
// Generators created with a non-zero seed are deterministic — two generators
// with the same seed will produce identical output for identical input. This
// is the contract relied on by the ?seed= query parameter and by tests.
//
// Generator is safe for concurrent use; the embedded rand source is guarded
// by a mutex.
type Generator struct {
	mu  sync.Mutex
	rng *rand.Rand
}

// New constructs a Generator. If seed != 0, the underlying source is seeded
// deterministically with that seed. If seed == 0, a time-based seed is used
// so that successive process restarts produce different streams.
//
// math/rand/v2 (Go 1.22+) is used; the deprecated math/rand is not.
func New(seed int64) *Generator {
	var src rand.Source
	if seed != 0 {
		// PCG takes two uint64 seeds; derive a deterministic second word from
		// the first so callers can pass a single int64.
		s := uint64(seed)
		src = rand.NewPCG(s, s^0x9E3779B97F4A7C15)
	} else {
		now := uint64(time.Now().UnixNano())
		src = rand.NewPCG(now, now^0x9E3779B97F4A7C15)
	}
	return &Generator{rng: rand.New(src)}
}

// pickLocked returns a uniformly random element from pool. Caller must hold
// g.mu. pool is assumed non-empty (callers check Tone validity first).
func (g *Generator) pickLocked(pool []string) string {
	return pool[g.rng.IntN(len(pool))]
}

// sortedToneNames is a cached, pre-sorted snapshot of every registered tone
// key. It is populated once at package init from Tones so RandomTone and
// AllToneNames avoid per-call map iteration + sort allocations. Tones is
// populated at init in fragments.go; Go guarantees package-level init runs
// before our init() here (same package).
var sortedToneNames []string

func init() {
	sortedToneNames = make([]string, 0, len(Tones))
	for k := range Tones {
		sortedToneNames = append(sortedToneNames, k)
	}
	sort.Strings(sortedToneNames)
}

// Compose builds a (message, subtitle) pair for the given tone using the
// supplied name and from values. If tone is not registered in Tones,
// Compose returns ("", "").
//
// The message is assembled as: opener + " " + body + closer.
// The subtitle is the chosen subtitlePrefix template with %s replaced by
// from. Each fragment slot is sampled independently.
func (g *Generator) Compose(tone, name, from string) (message, subtitle string) {
	t, ok := Get(tone)
	if !ok {
		return "", ""
	}
	g.mu.Lock()
	opener := g.pickLocked(t.Openers)
	body := g.pickLocked(t.Bodies)
	closer := g.pickLocked(t.Closers)
	subPrefix := g.pickLocked(t.SubtitlePrefixes)
	g.mu.Unlock()

	openerStr := fmt.Sprintf(opener, name)
	// Body fragments are static (they reference 9 directly) so they don't
	// take a name placeholder; we still concatenate via Sprintf for clarity.
	message = fmt.Sprintf("%s %s%s", openerStr, body, closer)
	subtitle = fmt.Sprintf(subPrefix, from)
	return message, subtitle
}

// RandomTone returns the name of a uniformly random registered tone. It is
// used by /api/random/{name}/{from} and is exported so callers can build
// their own random delegations on top of the same generator.
func (g *Generator) RandomTone() string {
	// sortedToneNames is pre-computed at init, so no per-call allocation or
	// sort is needed. Only the RNG access requires the lock.
	g.mu.Lock()
	defer g.mu.Unlock()
	return sortedToneNames[g.rng.IntN(len(sortedToneNames))]
}

// AllToneNames returns a sorted slice of every registered tone identifier.
// Useful for documentation generators and tests. The returned slice is a
// copy so callers cannot mutate the cached package-level slice.
func AllToneNames() []string {
	out := make([]string, len(sortedToneNames))
	copy(out, sortedToneNames)
	return out
}

// ContainsNineReference reports whether s mentions "9" or "nine" (case
// insensitive). Used by tests to enforce the WSAN invariant: every assembled
// message must reference the magic hour.
func ContainsNineReference(s string) bool {
	low := strings.ToLower(s)
	return strings.Contains(low, "9") || strings.Contains(low, "nine")
}
