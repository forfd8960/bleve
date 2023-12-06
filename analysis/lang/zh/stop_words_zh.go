package zh

import (
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
)

const StopName = "stop_zh"

// this content was obtained from:
// lucene-4.7.2/analysis/common/src/resources/org/apache/lucene/analysis/snowball/
// ` was changed to ' to allow for literal string

var ChineseStopWords = []byte(``)

func TokenMapConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.TokenMap, error) {
	rv := analysis.NewTokenMap()
	err := rv.LoadBytes(ChineseStopWords)
	return rv, err
}

func init() {
	registry.RegisterTokenMap(StopName, TokenMapConstructor)
}
