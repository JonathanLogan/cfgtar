package jsonschema

import (
	"strings"
)

func nameRequired(name string) (string, bool) {
	name = trimString(name)
	if strings.HasSuffix(name, requiredSuffix) {
		return trimString(name[:len(name)-len(requiredSuffix)]), true
	}
	return name, false
}

func validationData(s interface{}) (valFunc ValidatorFunc, required bool, err error) {
	var ok bool
	var q, funcName string
	if q, ok = s.(string); !ok {
		return nil, true, ErrSchemaDefType
	}
	funcName, required = nameRequired(q)
	if funcName == "" {
		funcName = defaultType
	}
	funcName, parameters := extractParameters(funcName)
	if valFunc, ok := validatorFuncMap[funcName]; ok {
		if parameters != nil && len(parameters) > 0 {
			return func(i ...any) (interface{}, error) {
				if len(i) > 0 {
					return valFunc(i[0], parameters)
				}
				return nil, ErrViolationType
			}, required, nil
		}
		return valFunc, required, nil
	}
	return nil, true, ErrSchemaDefValidator
}

func compareType(schema, data interface{}, required bool) (interface{}, error) {
	valFunc, required2, err := validationData(schema)
	if err != nil {
		return nil, err
	}
	required = required2 || required
	if required && data == nil {
		return nil, ErrRequired
	}
	if data == nil {
		return nil, nil
	}
	return valFunc(data)
}

func validateMap(schema map[string]interface{}, data interface{}, required bool) ([]string, interface{}, error) {
	var expand bool
	var dataV interface{}
	var dataT map[string]interface{}
	ret := make(map[string]interface{})
	if data != nil {
		if dataT, expand = data.(map[string]interface{}); !expand {
			return nil, nil, ErrSchemaType
		}
	}
	if data == nil && required {
		return nil, nil, ErrRequired
	}
	for k, v := range schema {
		dataV = nil
		k, required = nameRequired(k)
		if expand || required {
			var ok bool
			dataV, ok = dataT[k]
			if !ok && required {
				return []string{k}, nil, ErrRequired
			}
		}
		if p, d, err := validate(v, dataV, required); err != nil {
			return append(p, k), nil, err
		} else if expand {
			ret[k] = d
		}
	}
	return nil, ret, nil
}

func validateArray(schema []interface{}, data interface{}, required bool) ([]string, interface{}, error) {
	if len(schema) != 1 {
		return nil, nil, ErrArraySchema
	}
	if data == nil {
		if required {
			return nil, nil, ErrRequired
		}
		return nil, nil, nil
	}
	if dataV, ok := data.([]interface{}); ok {
		if len(dataV) == 0 {
			if required {
				return nil, nil, ErrRequired
			}
			return nil, make([]interface{}, 0), nil
		}
		ret := make([]interface{}, len(dataV))
		for k, v := range dataV {
			if p, d, err := validate(schema[0], v, required); err != nil {
				return append(p, indexString(k)), nil, err
			} else {
				ret[k] = d
			}
		}
		return nil, ret, nil
	}
	return nil, nil, ErrSchemaType
}

func validate(schema, data interface{}, required bool) ([]string, interface{}, error) {
	switch m := schema.(type) {
	case map[string]interface{}:
		return validateMap(m, data, required)
	case []interface{}:
		return validateArray(m, data, required)
	case interface{}:
		d, err := compareType(m, data, required)
		if err != nil {
			return nil, nil, err
		}
		return nil, d, nil
	default:
		return nil, nil, ErrUnknownType
	}
}
