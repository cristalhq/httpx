package httpx_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/cristalhq/httpx"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name      string
		header    string
		wantLang  string
		wantLangs []httpx.LangQ
	}{
		{
			name:     "empty header defaults to en",
			header:   "",
			wantLang: "en",
			wantLangs: []httpx.LangQ{
				{Lang: "en", Q: 1},
			},
		},
		{
			name:     "single language no q-value",
			header:   "what",
			wantLang: "what",
			wantLangs: []httpx.LangQ{
				{Lang: "what", Q: 1},
			},
		},
		{
			name:     "simple language code",
			header:   "uk",
			wantLang: "uk",
			wantLangs: []httpx.LangQ{
				{Lang: "uk", Q: 1},
			},
		},
		{
			name:     "multiple languages with q-values",
			header:   "en,en-GB;q=0.9,en-US;q=0.8",
			wantLang: "en",
			wantLangs: []httpx.LangQ{
				{Lang: "en", Q: 1},
				{Lang: "en-gb", Q: 0.9},
				{Lang: "en-us", Q: 0.8},
			},
		},
		{
			name:     "q-values with descending priority",
			header:   "fr;q=0.2,es;q=0.9,de;q=0.5",
			wantLang: "es",
			wantLangs: []httpx.LangQ{
				{Lang: "es", Q: 0.9},
				{Lang: "de", Q: 0.5},
				{Lang: "fr", Q: 0.2},
			},
		},
		{
			name:     "malformed q-value falls back to 1",
			header:   "jp;q=oops,cn;q=",
			wantLang: "jp",
			wantLangs: []httpx.LangQ{
				{Lang: "jp", Q: 1},
				{Lang: "cn", Q: 1},
			},
		},
		{
			name:     "trims spaces correctly",
			header:   " en ; q=0.5 , fr ",
			wantLang: "fr",
			wantLangs: []httpx.LangQ{
				{Lang: "fr", Q: 1},
				{Lang: "en", Q: 0.5},
			},
		},
		{
			name:     "invalid lang name still accepted",
			header:   "123-xyz;q=0.7",
			wantLang: "123-xyz",
			wantLangs: []httpx.LangQ{
				{Lang: "123-xyz", Q: 0.7},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &http.Request{
				Header: http.Header{
					"Accept-Language": []string{tc.header},
				},
			}

			lang := httpx.AcceptLanguage(r)
			if lang != tc.wantLang {
				t.Errorf("\nhave: %+v\nwant: %+v", lang, tc.wantLang)
			}

			langs := httpx.AcceptLanguages(r)
			if !reflect.DeepEqual(langs, tc.wantLangs) {
				t.Errorf("\nhave: %+v\nwant: %+v", langs, tc.wantLangs)
			}

			if len(langs) > 0 && lang != langs[0].Lang {
				t.Errorf("\nhave: %+v\nwant: %+v", lang, langs[0].Lang)
			}
		})
	}
}

func FuzzAcceptLanguages(f *testing.F) {
	// Seed with some interesting values
	seeds := []string{
		"", "en", "fr;q=0.9", "de;q=oops", "  es ; q= 0.3 , en ; q = 1 ",
		"zh;q=,ru;q=0.8", "abc-123;q=0.0001", "q;q=1.5", ";;;", "en;q=", "xx;q=0,yy;q=1",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, header string) {
		req := &http.Request{
			Header: http.Header{"Accept-Language": []string{header}},
		}

		langs := httpx.AcceptLanguages(req)
		if len(langs) == 0 {
			t.Errorf("no languages returned for header: %q", header)
		}

		// Ensure Q is always within [0,1] (though malformed inputs fallback to 1)
		for _, l := range langs {
			if l.Q < 0 || l.Q > 1 {
				t.Errorf("invalid Q=%v for lang=%q (header: %q)", l.Q, l.Lang, header)
			}
			if l.Lang == "" {
				t.Errorf("empty Lang returned (header: %q)", header)
			}
		}
	})
}
