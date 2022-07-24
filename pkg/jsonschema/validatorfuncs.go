package jsonschema

func isString(s ...interface{}) (interface{}, error) {
	var params ParamMap
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if len(s) > 1 {
		params = s[1].(ParamMap)
	}
	if q, ok := s[0].(string); ok {
		if min, ok, err := params.AsInt("min"); err != nil {
			return nil, err
		} else if ok {
			if len(q) < min {
				return nil, ErrParamConstraint
			}
		}
		if max, ok, err := params.AsInt("max"); err != nil {
			return nil, err
		} else if ok {
			if len(q) > max {
				return nil, ErrParamConstraint
			}
		}
		if l, ok, err := params.AsInt("len"); err != nil {
			return nil, err
		} else if ok {
			if len(q) != l {
				return nil, ErrParamConstraint
			}
		}
		if len(q) > 0 {
			return q, nil
		}
	}
	return nil, ErrViolationType
}

func isFloat(s ...interface{}) (interface{}, error) {
	var params ParamMap
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if len(s) > 1 {
		params = s[1].(ParamMap)
	}
	var q float64
	switch n := s[0].(type) {
	case float32:
		q = float64(n)
	case float64:
		q = n
	case int:
		q = float64(n)
	default:
		return nil, ErrViolationType
	}
	if min, ok, err := params.AsFloat("min"); err != nil {
		return nil, err
	} else if ok {
		if q < min {
			return nil, ErrParamConstraint
		}
	}
	if max, ok, err := params.AsFloat("max"); err != nil {
		return nil, err
	} else if ok {
		if q > max {
			return nil, ErrParamConstraint
		}
	}
	return q, nil

}

func isInt(s ...interface{}) (interface{}, error) {
	var params ParamMap
	if len(s) < 1 {
		return nil, ErrViolationType
	}
	if len(s) > 1 {
		params = s[1].(ParamMap)
	}
	var q int
	switch n := s[0].(type) {
	case float32:
		if s[0] == float32(int(n)) {
			q = int(n)
		} else {
			return nil, ErrViolationType
		}
	case float64:
		if s[0] == float64(int(n)) {
			q = int(n)
		} else {
			return nil, ErrViolationType
		}
	case int:
		q = n
	default:
		return nil, ErrViolationType
	}

	if min, ok, err := params.AsInt("min"); err != nil {
		return nil, err
	} else if ok {
		if q < min {
			return nil, ErrParamConstraint
		}
	}
	if max, ok, err := params.AsInt("max"); err != nil {
		return nil, err
	} else if ok {
		if q > max {
			return nil, ErrParamConstraint
		}
	}
	return q, nil
}
