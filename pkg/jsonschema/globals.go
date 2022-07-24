package jsonschema

import "errors"

type ValidatorFunc func(i any) (interface{}, error)
type ValidatorFuncMap map[string]ValidatorFunc

var (
	ErrUnknownType        = errors.New("unknown type")
	ErrViolationType      = errors.New("data type violates schema")
	ErrRequired           = errors.New("required")
	ErrArraySchema        = errors.New("array schema more than one element")
	ErrSchemaType         = errors.New("schema type not matched")
	ErrSchemaDefType      = errors.New("schema definition is not string")
	ErrSchemaDefValidator = errors.New("schema definition contains unknown getValidatorFunc type")
)

const (
	requiredStr    = "required"
	optionalStr    = "optional"
	defaultType    = "string"
	requiredSuffix = "%required"
)
