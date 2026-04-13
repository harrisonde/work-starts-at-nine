package generator

import (
	"strings"
	"testing"
)

// TestComposeDeterministicWithSeed verifies that two generators built with
// the same seed produce identical (msg, sub) pairs across several tones.
func TestComposeDeterministicWithSeed(t *testing.T) {
	tones := []string{"nine", "late", "boss", "hr", "rude", "polite"}
	for _, tone := range tones {
		tone := tone
		t.Run(tone, func(t *testing.T) {
			g1 := New(123)
			g2 := New(123)
			m1, s1 := g1.Compose(tone, "Alice", "Bob")
			m2, s2 := g2.Compose(tone, "Alice", "Bob")
			if m1 != m2 || s1 != s2 {
				t.Errorf("seeded compose differs: (%q,%q) vs (%q,%q)", m1, s1, m2, s2)
			}
		})
	}
}

// TestComposeDifferentSeeds asserts that at least one of five seed pairs
// produces different output. Statistical, not strict — picks may collide for
// any given pair, but five collisions across distinct seeds is implausible.
func TestComposeDifferentSeeds(t *testing.T) {
	pairs := [][2]int64{
		{1, 2},
		{10, 20},
		{100, 200},
		{1000, 2000},
		{42, 1337},
	}
	differs := 0
	for _, p := range pairs {
		ga := New(p[0])
		gb := New(p[1])
		ma, _ := ga.Compose("nine", "Alice", "Bob")
		mb, _ := gb.Compose("nine", "Alice", "Bob")
		if ma != mb {
			differs++
		}
	}
	if differs == 0 {
		t.Errorf("no seed pair produced different output across %d pairs", len(pairs))
	}
}

// TestComposeUnknownTone asserts the documented contract: unknown tones
// return ("", "").
func TestComposeUnknownTone(t *testing.T) {
	g := New(1)
	msg, sub := g.Compose("does-not-exist", "Alice", "Bob")
	if msg != "" || sub != "" {
		t.Errorf("unknown tone returned non-empty: msg=%q sub=%q", msg, sub)
	}
}

// TestComposeAllTonesNonEmpty walks every registered tone and asserts
// Compose returns a non-empty message AND subtitle, with the supplied name
// appearing somewhere in either string.
func TestComposeAllTonesNonEmpty(t *testing.T) {
	for tone := range Tones {
		tone := tone
		t.Run(tone, func(t *testing.T) {
			g := New(7)
			msg, sub := g.Compose(tone, "Alice", "Bob")
			if msg == "" {
				t.Errorf("empty message for tone %q", tone)
			}
			if sub == "" {
				t.Errorf("empty subtitle for tone %q", tone)
			}
			if !strings.Contains(msg, "Alice") && !strings.Contains(sub, "Alice") {
				t.Errorf("name %q not found in either msg=%q or sub=%q", "Alice", msg, sub)
			}
		})
	}
}

// TestAllTonesReferenceNine is the WSAN invariant test: every body fragment
// in every tone must reference 9 / nine. We compose 20 times per tone with
// different seeds and assert the resulting message contains "9" or "nine"
// every single time.
func TestAllTonesReferenceNine(t *testing.T) {
	for tone := range Tones {
		tone := tone
		t.Run(tone, func(t *testing.T) {
			for seed := int64(1); seed <= 20; seed++ {
				g := New(seed)
				msg, _ := g.Compose(tone, "Alice", "Bob")
				if !ContainsNineReference(msg) {
					t.Errorf("tone %q seed %d: msg %q lacks 9/nine reference", tone, seed, msg)
				}
			}
		})
	}
}

// TestFragmentsNonEmpty asserts every tone provides at least one of each
// fragment type so Compose can never index an empty slice.
func TestFragmentsNonEmpty(t *testing.T) {
	for name, tone := range Tones {
		if len(tone.Openers) == 0 {
			t.Errorf("tone %q: empty Openers", name)
		}
		if len(tone.Bodies) == 0 {
			t.Errorf("tone %q: empty Bodies", name)
		}
		if len(tone.Closers) == 0 {
			t.Errorf("tone %q: empty Closers", name)
		}
		if len(tone.SubtitlePrefixes) == 0 {
			t.Errorf("tone %q: empty SubtitlePrefixes", name)
		}
	}
}

// TestRandomToneIsDeterministic verifies that RandomTone, like Compose, is
// reproducible under a fixed seed. This is what makes /api/random?seed=N
// usable as a regression-test fixture.
func TestRandomToneIsDeterministic(t *testing.T) {
	g1 := New(99)
	g2 := New(99)
	for i := 0; i < 10; i++ {
		t1 := g1.RandomTone()
		t2 := g2.RandomTone()
		if t1 != t2 {
			t.Errorf("iter %d: RandomTone differs: %q vs %q", i, t1, t2)
		}
	}
}

// TestAllToneNamesSorted asserts the helper returns the keys of Tones in
// sorted order — relied on by RandomTone for seeded reproducibility.
func TestAllToneNamesSorted(t *testing.T) {
	names := AllToneNames()
	if len(names) != len(Tones) {
		t.Fatalf("AllToneNames len = %d, want %d", len(names), len(Tones))
	}
	for i := 1; i < len(names); i++ {
		if names[i-1] >= names[i] {
			t.Errorf("AllToneNames not sorted at index %d: %q >= %q", i, names[i-1], names[i])
		}
	}
}
