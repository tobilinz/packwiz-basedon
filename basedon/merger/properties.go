package merger

import (
	"bytes"
	"fmt"
	"regexp"
)

var properties = &Merger{
	Regex: regexp.MustCompile(`^config/(?:yosbr)?/(?:shaderpacks/.*\.txt)|(?:.*\.properties)$`),
	Function: func(base, this []byte, relPath string) ([]byte, error) {
		allLines := [][]byte{}
		allLines = append(allLines, bytes.Split(base, []byte{'\n'})...)
		allLines = append(allLines, bytes.Split(this, []byte{'\n'})...)

		mergedOptions := make(map[string][]byte)
		for _, line := range allLines {
			if len(line) == 0 || line[0] == '#' {
				continue
			}

			keyVal := bytes.SplitN(line, []byte{'='}, 2)

			if len(keyVal) != 2 {
				fmt.Println("Some properties file is misconfigured. Falling back to this packs config. Affected property: ", string(line), " File: ", relPath)
				return this, nil
			}

			mergedOptions[string(keyVal[0])] = keyVal[1]
		}

		merged := []byte{}
		for key, value := range mergedOptions {
			merged = append(merged, key...)
			merged = append(merged, '=')
			merged = append(merged, value...)
			merged = append(merged, '\n')
		}

		return merged, nil
	},
}
