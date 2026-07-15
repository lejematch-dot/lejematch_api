package citynorm

import "testing"

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"Aarhus":                   "Aarhus",
		"Aarhus C":                 "Aarhus",
		"Aarhus C, Frederiksbjerg": "Aarhus",
		"aarhus c, frederiksbjerg": "Aarhus",
		"Århus C":                  "Aarhus",
		"København SV":             "København",
		"København SV, Valby":      "København",
		"København":                "København",
		"Aalborg Ø, Østre Havn":    "Aalborg",
		"Odense":                   "Odense",
		"":                         "",
		"   ":                      "",
	}

	for input, want := range cases {
		if got := Normalize(input); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", input, got, want)
		}
	}
}
