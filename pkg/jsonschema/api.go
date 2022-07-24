package jsonschema

// Validate that data conforms to schema. Returns error and violating path.
func Validate(schema, data interface{}) (errPath []string, modified interface{}, err error) {
	pErr, d, err := validate(schema, data, false)
	if err != nil {
		return reverseStringSlice(pErr), nil, err
	}
	return nil, d, nil
}
