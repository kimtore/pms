package index

import (
	"github.com/ambientsound/pms/index/filters/unicodestrip"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/token/edgengram"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
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

	err = indexMapping.AddCustomTokenFilter("unicodeStripper",
		map[string]interface{}{
			"type": unicodestrip.Name,
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
				`unicodeStripper`,
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
