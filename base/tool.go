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

// DeepEqual reports whether x and y are “deeply equal,” defined as follows.
// []any: len of []any is equal, values in same position are deeply equal
// map[string]any: len of map[string]any is equal,  values with same key are deeply equal
// others: return the result of ==
func DeepEqual(v1, v2 any) bool {
	switch v1 := v1.(type) {
	case []any:
		v2, ok := v2.([]any)
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		valid := true
		for i := range v1 {
			valid = valid && DeepEqual(v1[i], v2[i])
		}
		return valid
	case map[string]any:
		v2, ok := v2.(map[string]any)
		if !ok {
			return false
		}
		if len(v1) != len(v2) {
			return false
		}
		valid := true
		for i := range v1 {
			valid = valid && DeepEqual(v1[i], v2[i])
		}
		return valid
	default:
		return v1 == v2
	}
}
