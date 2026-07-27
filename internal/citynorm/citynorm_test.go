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
		"Aarhus kommune":           "Aarhus",
		"Aarhus Kommune":           "Aarhus",
		"Aarhus eller Trøjborg":    "Aarhus",
		"Aalborg Øst":              "Aalborg",
		"København kommune":        "København",
		"Odense kommune":           "Odense",
		"":                         "",
		"   ":                      "",
	}

	for input, want := range cases {
		if got := Normalize(input); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", input, got, want)
		}
	}
}
