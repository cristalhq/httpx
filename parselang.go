package httpx

import (
	"cmp"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

// LangQ represents a language and its quality factor (q-value)
// from the HTTP Accept-Language header.
type LangQ struct {
	Lang string  // The language tag (e.g., "en", "fr", "en-us").
	Q    float64 // The quality factor, ranging from 0 to 1.
}

// AcceptLanguage returns the highest-preference language from the 'Accept-Language' header.
// If the header is missing or empty, it defaults to "en".
func AcceptLanguage(req *http.Request) string {
	return AcceptLanguages(req)[0].Lang
}

// AcceptLanguages parses the 'Accept-Language' header and returns a slice of [LangQ] sorted by descending quality factor.
// If no valid languages are found, it returns a default slice containing {"en", 1}.
func AcceptLanguages(req *http.Request) []LangQ {
	lang := req.Header.Get("Accept-Language")
	if lang == "" {
		return []LangQ{{Lang: "en", Q: 1}}
	}

	var langs []LangQ
	parts := strings.Split(lang, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		langParts := strings.Split(strings.ToLower(part), ";")

		// Default to "en" if languages is empty
		lang := cmp.Or(strings.TrimSpace(langParts[0]), "en")

		// Default q=1 if no q-value is provided or invalid
		if len(langParts) == 1 {
			langs = append(langs, LangQ{Lang: lang, Q: 1})
		} else {
			qp := strings.SplitN(langParts[1], "=", 2)
			q := 1.0
			if len(qp) == 2 {
				if parsedQ, err := strconv.ParseFloat(qp[1], 64); err == nil {
					q = parsedQ
				}
			}
			langs = append(langs, LangQ{Lang: lang, Q: max(0, min(q, 1))})
		}
	}

	slices.SortStableFunc(langs, func(a, b LangQ) int {
		return cmp.Compare(b.Q, a.Q)
	})
	return langs
}
