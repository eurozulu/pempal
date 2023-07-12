package utils

func MergeMap(temps ...map[string]string) map[string]string {
	m := map[string]string{}
	CopyMapInto(m, temps...)
	return m
}

func CopyMapInto(dst map[string]string, src ...map[string]string) {
	for _, t := range src {
		copyMapInto(dst, t)
	}
}

func copyMapInto(dst, src map[string]string) {
	for k, v := range src {
		_, exists := dst[k]
		if exists && v == "" {
			// skip overwriting with empty value
			continue
		}
		dst[k] = v
	}
}

func MapKeyNames(m map[string]string) []string {
	names := make([]string, len(m))
	var i int
	for k := range m {
		names[i] = k
		i++
	}
	return names
}
