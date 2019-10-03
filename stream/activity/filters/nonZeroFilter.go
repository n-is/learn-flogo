package filters

type NonZeroFilter struct{}

func (*NonZeroFilter) FilterOut(val interface{}) (bool, interface{}) {
	return IsNonZero(val)
}

func IsNonZero(val interface{}) (bool, interface{}) {

	fOut := false
	switch t := val.(type) {
	case int:
		return t != 0, t
	case float64:
		return t != 0.0, t
	case []int:
		var vs []interface{} = make([]interface{}, len(t))
		vs_len := 0
		for _, v := range t {
			if v != 0 {
				vs[vs_len] = v
				vs_len++
			} else {
				fOut = true
			}
		}
		return fOut, vs[0:vs_len]
	case []float64:
		var vs []interface{} = make([]interface{}, len(t))
		vs_len := 0
		for _, v := range t {
			if v != 0.0 {
				vs[vs_len] = v
				vs_len++
			} else {
				fOut = true
			}
		}
		return fOut, vs[0:vs_len]
	}

	// For Unsupported Types
	return false, nil
}
