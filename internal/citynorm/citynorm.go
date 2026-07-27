// Package citynorm normaliserer bynavne, så stavevarianter og
// postnummer-distrikt-suffikser (f.eks. "Århus", "Aarhus C", "aarhus")
// alle bliver til samme værdi. Bruges både når et opslag oprettes/redigeres
// og til at rette allerede gemte data, så by-filteret ikke viser den
// samme by flere gange.
package citynorm

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var aliases = map[string]string{
	"aarhus":     "Aarhus",
	"århus":      "Aarhus",
	"aalborg":    "Aalborg",
	"ålborg":     "Aalborg",
	"kobenhavn":  "København",
	"københavn":  "København",
	"koebenhavn": "København",
	"odense":     "Odense",
}

// Matcher et afsluttende postnummer-distrikt, f.eks. " C", " N", " SV".
var suffixPattern = regexp.MustCompile(`(?i)\s+(NV|NØ|SV|SØ|[NSVØC])$`)

func titleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) == 0 {
			continue
		}
		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		for j := 1; j < len(runes); j++ {
			runes[j] = unicode.ToLower(runes[j])
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// Normalize gør et bynavn ensartet: fjerner postnummer-distrikt-suffikser,
// slår kendte stavevarianter sammen, og ensretter store/små bogstaver.
func Normalize(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return trimmed
	}

	// Drop alt efter et komma — bydel/kvarter-tilføjelser som
	// "Aarhus C, Frederiksbjerg" skal normalisere til samme by som "Aarhus C".
	if idx := strings.Index(trimmed, ","); idx != -1 {
		if head := strings.TrimSpace(trimmed[:idx]); head != "" {
			trimmed = head
		}
	}

	city := strings.TrimSpace(suffixPattern.ReplaceAllString(trimmed, ""))
	if city == "" {
		city = trimmed
	}

	// Fritekstfeltet lader folk skrive hvad som helst efter bynavnet (fx
	// "Aarhus kommune", "Aarhus eller Trøjborg", "Aarhus/ Viby K") — for de
	// store byer i aliases-listen samler vi alt der starter med bynavnet
	// under samme filter, uanset hvad der følger efter. Matcher på prefix
	// (ikke kun hele ord), da folk ikke altid sætter mellemrum før resten
	// ("Aarhus/Viby").
	lowerCity := strings.ToLower(city)
	for key, canonical := range aliases {
		if !strings.HasPrefix(lowerCity, key) {
			continue
		}
		rest := lowerCity[len(key):]
		if rest == "" {
			return canonical
		}
		r, _ := utf8.DecodeRuneInString(rest)
		if !unicode.IsLetter(r) {
			return canonical
		}
	}

	return titleCase(city)
}
