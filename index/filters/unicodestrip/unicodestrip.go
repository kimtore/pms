package unicodestrip

import (
	"unicode"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const Name = "strip_unicode"

// StripUnicodeFilter is a Bleve keyword filter which decomposes unicode
// strings into their normalized form and strips away non-spacing marks.
// Effectively, this strips away diacritic marks so that searches may be done
// without entering them, e.g. "Télépopmusik" is indexed as "Telepopmusik".
type StripUnicodeFilter struct {
}

func New() (*StripUnicodeFilter, error) {
	return &StripUnicodeFilter{}, nil
}

func Constructor(config map[string]interface{}, cache *registry.Cache) (analysis.TokenFilter, error) {
	return New()
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func (s *StripUnicodeFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	chain := transform.Chain(norm.NFKD, transform.RemoveFunc(isMn), norm.NFC)
	for _, token := range input {
		token.Term, _, _ = transform.Bytes(chain, token.Term)
	}
	return input
}

func init() {
	registry.RegisterTokenFilter(Name, Constructor)
}
