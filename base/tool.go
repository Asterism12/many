package base

func deepCopy(src any) any {
	switch src := src.(type) {
	case []any:
		var dst []any
		for _, v := range src {
			dst = append(dst, deepCopy(v))
		}
		return dst
	case map[string]any:
		dst := map[string]any{}
		for k, v := range src {
			dst[k] = deepCopy(v)
		}
		return dst
	default:
		return src
	}
}

// Rest return rest of slices
func Rest[T any](s []T) []T {
	if len(s) <= 1 {
		return nil
	}
	return s[1:]
}
