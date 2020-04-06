package baseutils

// MergeMaps merges all maps in the list
func MergeMaps(MapList ...map[string]string) map[string]string {
	var baseMap = make(map[string]string)
	for _, imap := range MapList {
		for k, v := range imap {
			baseMap[k] = v
		}
	}
	return baseMap
}

// MergeMapsWithoutOverwrite merges all maps in the list without overwriting. The priority is the same as the sequence of the list.
func MergeMapsWithoutOverwrite(MapList ...map[string]string) map[string]string {
	var baseMap = make(map[string]string)
	for _, imap := range MapList {
		for k, v := range imap {
			if _, ok := baseMap[k]; !ok {
				baseMap[k] = v
			}

		}
	}
	return baseMap
}
