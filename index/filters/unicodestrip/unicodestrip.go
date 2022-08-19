// Package unicodestrip provides a Bleve keyword filter which decomposes unicode strings.
package unicodestrip

import (
	"unicode"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
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

// New returns a new instance of StripUnicodeFilter.
func New() (*StripUnicodeFilter, error) {
	return &StripUnicodeFilter{}, nil
}

// Constructor provides a constructor for Bleve.
func Constructor(config map[string]interface{}, cache *registry.Cache) (analysis.TokenFilter, error) {
	return New()
}

// isMn returns true if the provided rune is a unicode non-spacing mark.
func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

// Filter removes non-spacing marks from text in a token stream.
func (s *StripUnicodeFilter) Filter(input analysis.TokenStream) analysis.TokenStream {
	chain := transform.Chain(norm.NFKD, transform.RemoveFunc(isMn), norm.NFC)
	for _, token := range input {
		token.Term, _, _ = transform.Bytes(chain, token.Term)
	}
	return input
}

// init registers this plugin with Bleve.
func init() {
	registry.RegisterTokenFilter(Name, Constructor)
}
