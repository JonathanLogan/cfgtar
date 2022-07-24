package jsonschema

import (
	"strconv"
	"strings"
)

type ParamMap map[string]interface{}

func (params ParamMap) AsFloat(key string) (float64, bool, error) {
	if params == nil {
		return 0, false, nil
	}
	if e, ok := params[key]; ok {
		if e == nil {
			return 0, true, nil
		}
		if eI, ok := e.(float64); ok {
			return eI, true, nil
		}
		if eI, ok := e.(float32); ok {
			return float64(eI), true, nil
		}
		if eI, ok := e.(string); ok {
			eX, err := strconv.ParseFloat(eI, 64)
			if err != nil {
				return 0, false, err
			}
			return eX, true, nil
		}
		return 0, true, ErrParamType
	}
	return 0, false, nil
}
func (params ParamMap) AsInt(key string) (int, bool, error) {
	if params == nil {
		return 0, false, nil
	}
	if e, ok := params[key]; ok {
		if e == nil {
			return 0, true, nil
		}
		if eI, ok := e.(int); ok {
			return eI, true, nil
		}
		if eI, ok := e.(string); ok {
			eX, err := strconv.ParseInt(eI, 10, 32)
			if err != nil {
				return 0, false, err
			}
			return int(eX), true, nil
		}
		return 0, true, ErrParamType
	}
	return 0, false, nil
}
func (params ParamMap) AsString(key string) (string, bool, error) {
	if params == nil {
		return "", false, nil
	}
	if e, ok := params[key]; ok {
		if e == nil {
			return "", true, nil
		}
		if eI, ok := e.(string); ok {
			return eI, true, nil
		}
		return "", true, ErrParamType
	}
	return "", false, nil
}

func extractParameters(s string) (funcName string, params ParamMap) {
	if pos := strings.Index(s, "("); pos >= 0 {
		if pos2 := strings.Index(s[1+pos:], ")"); pos2 >= 0 {
			ret := make(ParamMap)
			for _, e := range strings.Split(s[1+pos:1+pos+pos2], ",") {
				key := strings.ToLower(trimString(e))
				if posE := strings.Index(key, "="); posE >= 0 {
					value := trimString(key[posE+1:])
					key = trimString(key[:posE])
					ret[key] = value
				} else {
					ret[key] = nil
				}
			}
			return s[:pos], ret
		}
	}
	return s, nil
}
