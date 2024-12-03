package merger

import (
	"regexp"
)

var fabricLoaderDependencies = &Merger{
	Regex: regexp.MustCompile(`^config/fabric_loader_dependencies\.json$`),
	Function: func(base, this []byte, relPath string) ([]byte, error) {
		/*//TODO: version has to be first json value. < and > are replace with \u...
		var baseMap, thisMap map[string]interface{}

		if err := jsonEncoding.Unmarshal(base, &baseMap); err != nil {
			return nil, err
		}
		if err := jsonEncoding.Unmarshal(this, &thisMap); err != nil {
			return nil, err
		}

		merge(baseMap, thisMap)

		rawVersion, found := baseMap["version"]

		if !found {
			return this, nil
		}

		version, ok := rawVersion.([]byte)

		if !ok {
			return nil, errors.New("Could not convert version entry in fabric_loader_dependencies.json to a []byte.")
		}

		delete(baseMap, "version")

		json, err := jsonEncoding.Marshal(baseMap)

		if err != nil {
			return nil, err
		}

		index := bytes.IndexByte(json, '{')

		if index < 0 {
			return nil, errors.New("Could not find opening parenthesis in resulting fabric_loader_dependencies.json.")
		}

		result := make([]byte, 0, len(json)-1+len(version))
		result = append(result, json[:index]...)
		result = append(result, version...)
		result = append(result, json[index:]...)

		return result, nil
		*/

		return this, nil
	},
}
