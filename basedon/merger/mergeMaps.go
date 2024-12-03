package merger

// TODO: Doesn't seem to work
func merge(dst, src map[string]interface{}) {
	for key, value := range src {
		srcMap, ok := value.(map[string]interface{})

		if ok {
			dstMap, ok := dst[key].(map[string]interface{})

			if ok {
				merge(dstMap, srcMap)
				continue
			}
		}

		dst[key] = value
	}
}
