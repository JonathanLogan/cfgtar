package jsonschema

var validatorFuncMap = make(ValidatorFuncMap)

func RegisterValidatorFunc(typeDef string, f ValidatorFunc) {
	validatorFuncMap[typeDef] = f
}

func init() {
	RegisterValidatorFunc("string", isString)
	RegisterValidatorFunc("float", isFloat)
	RegisterValidatorFunc("int", isInt)
	RegisterValidatorFunc("dir", isDir)
	RegisterValidatorFunc("file", isFile)
	RegisterValidatorFunc("duration", isDuration)
	RegisterValidatorFunc("hex", isHex)
	RegisterValidatorFunc("base64", isBase64)
	RegisterValidatorFunc("base58", isBase58)
	RegisterValidatorFunc("ipv4", isIPv4Addr)
	RegisterValidatorFunc("ipv6", isIPv6Addr)
	RegisterValidatorFunc("ipv4net", isIPv4Net)
	RegisterValidatorFunc("ipv6net", isIPv6Net)
	RegisterValidatorFunc("hostname", isHostname)
	RegisterValidatorFunc("nic", isNIC)
	RegisterValidatorFunc("nic4", isNIC4)
	RegisterValidatorFunc("nic6", isNIC6)
	RegisterValidatorFunc("lookup4", lookupIPv4)
	RegisterValidatorFunc("lookup6", lookupIPv6)
}
