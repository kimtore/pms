package index

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/token/edgengram"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/token/unicodenorm"
	"github.com/blevesearch/bleve/analysis/tokenizer/whitespace"
	"github.com/blevesearch/bleve/mapping"
)

// buildIndexMapping() returns an object that defines how input data is indexed in Bleve.
func buildIndexMapping() (mapping.IndexMapping, error) {
	indexMapping := bleve.NewIndexMapping()

	var err error

	err = indexMapping.AddCustomTokenFilter("songEdgeNgram",
		map[string]interface{}{
			"min":  float64(2),
			"max":  float64(25),
			"type": edgengram.Name,
		})
	if err != nil {
		return nil, err
	}

	err = indexMapping.AddCustomTokenFilter("unicodeNormalizer",
		map[string]interface{}{
			"form": unicodenorm.NFKD,
			"type": unicodenorm.Name,
		})
	if err != nil {
		return nil, err
	}

	err = indexMapping.AddCustomAnalyzer("songAnalyzer",
		map[string]interface{}{
			"type":         custom.Name,
			"char_filters": []interface{}{},
			"tokenizer":    whitespace.Name,
			"token_filters": []interface{}{
				`unicodeNormalizer`,
				lowercase.Name,
				`songEdgeNgram`,
			},
		})
	if err != nil {
		return nil, err
	}

	indexMapping.DefaultAnalyzer = "songAnalyzer"

	return indexMapping, nil
}
