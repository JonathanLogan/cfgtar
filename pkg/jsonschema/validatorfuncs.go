package jsonschema

func isString(s any) (interface{}, error) {
	if q, ok := s.(string); ok {
		if len(q) > 0 {
			return q, nil
		}
	}
	return nil, ErrViolationType
}

func isFloat(s any) (interface{}, error) {
	switch s.(type) {
	case float32, float64:
		return s, nil
	}
	return nil, ErrViolationType
}

func isInt(s any) (interface{}, error) {
	switch q := s.(type) {
	case float32:
		if s == float32(int(q)) {
			return int(q), nil
		}
	case float64:
		if s == float64(int(q)) {
			return int(q), nil
		}
	}
	return nil, ErrViolationType
}
