// Package generator builds randomized WSAN messages from per-tone fragment
// pools. Each tone owns its own openers/bodies/closers/subtitle prefixes so
// that every request to /api/<tone>/{name}/{from} composes a slightly
// different sentence while still hammering home the point that work starts
// at nine.
package generator

// Tone is a bundle of fragment pools used to procedurally assemble a single
// WSAN message + subtitle pair. Each pool entry is a printf format string
// that takes a single %s placeholder for substitution at compose time.
//
//   - Openers receive the recipient's name.
//   - Bodies receive the recipient's name and MUST reference 9 / nine / 09:00.
//   - Closers receive the recipient's name (or are static).
//   - SubtitlePrefixes receive the From value.
type Tone struct {
	Openers          []string
	Bodies           []string
	Closers          []string
	SubtitlePrefixes []string
}

// Tones is the master registry of every supported WSAN tone. Keys must
// match the tone identifier passed to Generator.Compose.
var Tones = map[string]Tone{
	"nine": {
		Openers: []string{
			"Hey %s,",
			"Listen %s,",
			"%s,",
			"Attention %s:",
			"Quick one %s —",
			"Heads up %s,",
		},
		Bodies: []string{
			"work starts at 9",
			"the day kicks off at 9 AM",
			"9:00 means 9:00",
			"nine. Not nine-oh-five",
			"the office opens its arms at 9",
			"9 AM is the start, not a suggestion",
		},
		Closers: []string{
			".",
			", every weekday.",
			". Sharp.",
			". Always has been.",
			". That's the deal.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"— %s, ops",
			"signed, %s",
			"posted by %s",
		},
	},
	"late": {
		Openers: []string{
			"You're late, %s.",
			"%s, you're late.",
			"Tardy alert: %s.",
			"%s — late again.",
			"Late slip for %s:",
			"Clock check, %s:",
		},
		Bodies: []string{
			"Work started at 9",
			"The 9 AM bell rang without you",
			"Your 09:00 came and went",
			"Nine has already happened",
			"It is officially after 9",
			"9 AM was your call time",
		},
		Closers: []string{
			".",
			". Try earlier tomorrow.",
			". Again.",
			". Catch up fast.",
			". This is on the record.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"— %s, attendance",
			"logged by %s",
			"— %s (timekeeper)",
		},
	},
	"reminder": {
		Openers: []string{
			"Friendly reminder, %s:",
			"Just a nudge, %s —",
			"Quick reminder %s,",
			"Hi %s, gentle ping:",
			"Heads-up %s:",
			"%s, a soft reminder:",
		},
		Bodies: []string{
			"work starts at 9 AM sharp",
			"the day begins at 9",
			"9 AM is our official kickoff",
			"please plan around a 9:00 start",
			"the schedule still says 9 AM",
			"nine in the morning is go-time",
		},
		Closers: []string{
			".",
			". Thanks!",
			". Appreciate it.",
			". You got this.",
			". Cheers.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"— %s, friendly desk",
			"with care, %s",
			"— %s 🙂",
		},
	},
	"strict": {
		Openers: []string{
			"%s.",
			"Listen carefully, %s.",
			"%s — final warning.",
			"Pay attention, %s.",
			"%s. Read this once.",
			"%s, no debate:",
		},
		Bodies: []string{
			"Nine. AM. Every day. No exceptions",
			"09:00. Not a minute later",
			"You will be at your desk by 9",
			"9 AM is mandatory, not optional",
			"Be here by 9. Period",
			"Nine sharp. That is the rule",
		},
		Closers: []string{
			".",
			". Understood?",
			". End of discussion.",
			". Don't make me repeat it.",
			". That is all.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"— %s, management",
			"— %s (final notice)",
			"signed, %s",
		},
	},
	"boss": {
		Openers: []string{
			"%s, my office.",
			"%s — a word.",
			"Hey %s, real quick.",
			"%s, between you and me:",
			"%s, look —",
			"%s, I'm not gonna sugarcoat it.",
		},
		Bodies: []string{
			"We start at 9, not 9:15",
			"The team is here at 9. So are you, going forward",
			"9 AM is when we go live",
			"I expect you at 9 sharp",
			"You and I both know 9 means 9",
			"Around here, 9 AM is the floor",
		},
		Closers: []string{
			".",
			". We good?",
			". Let's fix it.",
			". Don't let it happen again.",
			". I'm counting on you.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"— %s, the boss",
			"— %s (corner office)",
			"— %s, your manager",
		},
	},
	"hr": {
		Openers: []string{
			"Dear %s,",
			"%s, per company policy:",
			"To: %s",
			"%s — formal notice:",
			"%s, this is a courtesy notification.",
			"Re: attendance — %s,",
		},
		Bodies: []string{
			"per the employee handbook §4.2, working hours commence at 09:00",
			"standard working hours begin at 9:00 AM local time",
			"the start-of-business is defined as 09:00 in §2.1",
			"per your employment agreement, you are expected on-site at 9 AM",
			"company policy stipulates a 09:00 start time",
			"working hours per §4 begin promptly at nine in the morning",
		},
		Closers: []string{
			".",
			". Please govern yourself accordingly.",
			". This notice will be filed.",
			". Compliance is appreciated.",
			". Let us know if you have questions.",
		},
		SubtitlePrefixes: []string{
			"— HR (%s)",
			"— %s, People Operations",
			"— %s, HR Business Partner",
			"signed, %s, HR",
		},
	},
	"polite": {
		Openers: []string{
			"Good morning %s —",
			"Hi %s,",
			"Hello %s,",
			"%s, hope you're well —",
			"Dear %s,",
			"Morning %s!",
		},
		Bodies: []string{
			"just a gentle note that we kick off at 9",
			"if it's not too much trouble, the day begins at 9",
			"a little reminder that we start at 9 AM",
			"whenever you can, please aim for a 9 AM arrival",
			"our official start is 9 in the morning",
			"we'd love to see you by 9 if at all possible",
		},
		Closers: []string{
			". Thank you!",
			". So appreciated.",
			". Many thanks!",
			". Have a wonderful day!",
			". Cheers.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"warmly, %s",
			"with thanks, %s",
			"— %s 🌼",
		},
	},
	"rude": {
		Openers: []string{
			"%s, get in here.",
			"Oi %s.",
			"Yo %s.",
			"%s — seriously?",
			"%s, are you kidding me?",
			"%s, wake up.",
		},
		Bodies: []string{
			"Work. Starts. At. Nine",
			"Nine AM, genius",
			"It's 9. NINE. Not whenever you feel like it",
			"Your bed is not the office. 9 AM",
			"NINE. AY. EM",
			"You know what time it is? Past 9, that's what",
		},
		Closers: []string{
			".",
			". Move it.",
			". Now.",
			". Unbelievable.",
			". Get a clock.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"— %s, no patience left",
			"sent from %s's last nerve",
			"— %s 🙄",
		},
	},
	"monday": {
		Openers: []string{
			"It's Monday, %s.",
			"Happy Monday %s.",
			"%s, it's Monday.",
			"Morning %s — yes, Monday.",
			"%s, the week begins.",
			"Brace yourself %s, Monday:",
		},
		Bodies: []string{
			"Yes, work still starts at 9",
			"9 AM hasn't moved just because it's Monday",
			"The weekend ends, 9 AM begins",
			"Monday or not, 9 means 9",
			"Even today, the start is 9",
			"Nine AM. Mondays included",
		},
		Closers: []string{
			".",
			". Coffee helps.",
			". You'll survive.",
			". One day at a time.",
			". Let's go.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"— %s, Monday desk",
			"— %s ☕",
			"sent reluctantly by %s",
		},
	},
	"meeting": {
		Openers: []string{
			"%s,",
			"Reminder %s:",
			"%s, calendar check —",
			"Hey %s,",
			"%s — quick one:",
			"FYI %s,",
		},
		Bodies: []string{
			"the 9 AM is at 9 AM. That's why it's called the 9 AM",
			"the 9 o'clock meeting starts at 9 o'clock",
			"the 0900 sync starts at 0900",
			"the 9 AM standup is, surprisingly, at 9 AM",
			"calendar says 9. Because the meeting is at 9",
			"the meeting named 9 AM begins at 9 AM",
		},
		Closers: []string{
			".",
			". Be there.",
			". Don't be that person.",
			". Cameras on.",
			". Link in the invite.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"— %s, organizer",
			"— %s, calendar",
			"sent by %s",
		},
	},
	"earlybird": {
		Openers: []string{
			"%s,",
			"Listen %s,",
			"Quick proverb, %s:",
			"%s — wisdom incoming:",
			"Heads up %s,",
			"FYI %s:",
		},
		Bodies: []string{
			"the early bird gets the worm. The 9 AM bird gets a warning",
			"the 8 AM bird thrives. The 9:01 bird is unemployed",
			"early birds eat. Nine-AM birds eat. 9:15 birds get nothing",
			"birds at 9 are fine. Birds after 9 are footnotes",
			"the worm is gone by 9:01",
			"be the bird at 9, not the bird at 9-something",
		},
		Closers: []string{
			".",
			". Be a bird.",
			". Worm awaits.",
			". Tweet tweet.",
			". You know the drill.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"— %s, ornithology dept",
			"— %s 🐦",
			"chirped by %s",
		},
	},
	"coffee": {
		Openers: []string{
			"%s,",
			"Hey %s,",
			"Real talk %s:",
			"%s — heads up:",
			"PSA %s:",
			"Quick one %s,",
		},
		Bodies: []string{
			"finish the coffee on your way in. 9 means 9",
			"sip and walk. 9 AM doesn't wait for your latte",
			"the coffee shop is not the office. 9 AM",
			"to-go cup, 9 AM start, simple math",
			"caffeine on the move — desk by 9",
			"drink it walking. Be at your desk at 9",
		},
		Closers: []string{
			".",
			". Lid on.",
			". No refills.",
			". Go go go.",
			". Tick tock.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"— %s, barista liaison",
			"— %s ☕",
			"sent from %s's travel mug",
		},
	},
	"standup": {
		Openers: []string{
			"%s,",
			"Reminder %s:",
			"Hey %s,",
			"FYI %s,",
			"%s — standup note:",
			"Quick one %s:",
		},
		Bodies: []string{
			"standup is at 9. Standing somewhere else doesn't count",
			"the 9 AM standup starts at 9 AM",
			"daily standup, 9 sharp, no exceptions",
			"if you're not standing at 9, you're not standing up",
			"standup at 9 — yesterday/today/blockers",
			"9 AM, three questions, two minutes",
		},
		Closers: []string{
			".",
			". Be ready.",
			". Camera on.",
			". Keep it short.",
			". See you there.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"— %s, scrum",
			"— %s (facilitator)",
			"sent by %s",
		},
	},
	"deadline": {
		Openers: []string{
			"%s,",
			"Final notice %s:",
			"%s — deadline alert:",
			"Heads up %s,",
			"%s, time check:",
			"Attention %s:",
		},
		Bodies: []string{
			"the deadline was 9 AM. It is now later than 9 AM",
			"9:00 was the cutoff. The cutoff has cut off",
			"the 9 AM deadline has officially passed",
			"deadline: 9 AM. Current state: not 9 AM anymore",
			"submission window closed at 9",
			"you had until 9. It is no longer until 9",
		},
		Closers: []string{
			".",
			". Please advise.",
			". Escalating now.",
			". Let me know your plan.",
			". Acknowledge receipt.",
		},
		SubtitlePrefixes: []string{
			"— %s",
			"— %s, project lead",
			"— %s (PMO)",
			"sent by %s",
		},
	},
}

// Get returns the Tone bundle for the given identifier and a boolean
// indicating whether the tone exists. It is a thin wrapper around the Tones
// map kept exported so callers don't need to import the map directly.
func Get(tone string) (Tone, bool) {
	t, ok := Tones[tone]
	return t, ok
}
