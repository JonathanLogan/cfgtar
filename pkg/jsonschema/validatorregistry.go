package jsonschema

var validatorFuncMap = make(ValidatorFuncMap)

func RegisterValidatorFunc(typeDef string, f ValidatorFunc) {
	validatorFuncMap[typeDef] = f
}

func init() {
	RegisterValidatorFunc("string", isString)
	RegisterValidatorFunc("float", isFloat)
	RegisterValidatorFunc("int", isInt)
}
