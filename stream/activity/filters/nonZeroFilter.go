package filters

type NonZeroFilter struct{}

func (*NonZeroFilter) FilterOut(val interface{}) (bool, interface{}) {
	return IsNonZero(val)
}

func IsNonZero(val interface{}) (bool, interface{}) {

	fOut := false
	switch t := val.(type) {
	case int:
		return t != 0, (t + 1)
	case float64:
		return t != 0.0, (t + 0.1)
	case []int:
		for i, v := range t {
			if v != 0 {
				fOut = true
				t[i] = v + 1
			}
		}
		return fOut, t
	case []float64:
		for i, v := range t {
			if v != 0.0 {
				fOut = true
				t[i] = v + 0.1
			}
		}
		return fOut, t
	}

	// For Unsupported Types
	return false, nil
}
