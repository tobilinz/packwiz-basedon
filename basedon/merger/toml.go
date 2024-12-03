package merger

import (
	"bytes"
	"regexp"

	tomlEncoding "github.com/BurntSushi/toml"
)

var toml = &Merger{
	Regex: regexp.MustCompile(`^.*\.toml$`),
	Function: func(base, this []byte, relPath string) ([]byte, error) {
		var baseMap, thisMap map[string]interface{}

		if _, err := tomlEncoding.Decode(string(base), &baseMap); err != nil {
			return nil, err
		}
		if _, err := tomlEncoding.Decode(string(this), &thisMap); err != nil {
			return nil, err
		}

		merge(baseMap, thisMap)

		var buf bytes.Buffer
		encoder := tomlEncoding.NewEncoder(&buf)
		if err := encoder.Encode(baseMap); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	},
}
