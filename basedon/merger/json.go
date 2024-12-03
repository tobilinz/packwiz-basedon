package merger

import (
	jsonEncoding "encoding/json"
	"regexp"
)

var json = &Merger{
	Regex: regexp.MustCompile(`^.*\.json$`),
	Function: func(base, this []byte, relPath string) ([]byte, error) {
		var baseMap, thisMap map[string]interface{}

		if err := jsonEncoding.Unmarshal(base, &baseMap); err != nil {
			return nil, err
		}
		if err := jsonEncoding.Unmarshal(this, &thisMap); err != nil {
			return nil, err
		}

		merge(baseMap, thisMap)

		return jsonEncoding.Marshal(baseMap)
	},
}
