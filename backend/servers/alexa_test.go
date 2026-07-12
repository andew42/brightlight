package servers

import (
	"github.com/andew42/brightlight/animations"
	"testing"
)

func namedButtons(names ...string) []animations.Button {

	buttons := make([]animations.Button, len(names))
	for i, n := range names {
		buttons[i] = animations.Button{Key: i + 1, Name: n}
	}
	return buttons
}

func TestCanonicaliseButtonName(t *testing.T) {

	cases := map[string]string{
		"Sweet Shop": "sweet shop",
		"  OFF ":     "off",
		"Two Tone!":  "2 tone",
		"2 tone":     "2 tone",
		"Baby-Bows":  "baby bows",
		"":           "",
		"Three":      "3",
	}
	for in, want := range cases {
		if got := canonicaliseButtonName(in); got != want {
			t.Errorf("canonicaliseButtonName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFindButtonByName(t *testing.T) {

	buttons := namedButtons("OFF", "Full", "Sweet Shop", "Two Tone", "Rainbow", "One")

	cases := map[string]string{
		"off":        "OFF",
		"sweet shop": "Sweet Shop",
		"Sweet Shop": "Sweet Shop",
		"two tone":   "Two Tone", // spoken "two" matches button "Two"
		"2 tone":     "Two Tone",
		"rainbows":   "Rainbow", // small speech to text error forgiven
		"sweet shot": "Sweet Shop",
		"1":          "One", // digit heard, word named button
	}
	for spoken, want := range cases {
		b := findButtonByName(buttons, spoken)
		if b == nil {
			t.Errorf("findButtonByName(%q) = nil, want %q", spoken, want)
		} else if b.Name != want {
			t.Errorf("findButtonByName(%q) = %q, want %q", spoken, b.Name, want)
		}
	}

	// Short names must not fuzzy match each other or gibberish
	for _, spoken := range []string{"on", "disco", ""} {
		if b := findButtonByName(buttons, spoken); b != nil {
			t.Errorf("findButtonByName(%q) = %q, want no match", spoken, b.Name)
		}
	}
}
