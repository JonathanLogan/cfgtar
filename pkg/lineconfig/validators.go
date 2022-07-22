package lineconfig

import (
	"fmt"
	"strconv"
)

type ValidatorFunc func(any) (interface{}, error)

var validatorMap = make(map[string]ValidatorFunc)

func init() {
	AddValidateFunc("string", StringTest)
	AddValidateFunc("int", IntTest)
}

func AddValidateFunc(valType string, valFunc ValidatorFunc) {
	validatorMap[valType] = valFunc
}

func Validate(valType string, value any) (interface{}, error) {
	if valFunc, ok := validatorMap[valType]; ok {
		return valFunc(value)
	}
	return nil, ErrUnknownType
}

func StringTest(i any) (interface{}, error) {
	if s, ok := i.(string); ok {
		return s, nil
	}
	return nil, fmt.Errorf("Expected 'string': %v", i)
}

func IntTest(i any) (interface{}, error) {
	if s, ok := i.(string); ok {
		return strconv.ParseInt(s, 10, 64)
	}
	return nil, fmt.Errorf("Expected 'integer': %v", i)
}
