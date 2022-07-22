package lineconfig

import (
	"errors"
)

var (
	ErrPathTermination = errors.New("path terminates in non-value")
	ErrPathType        = errors.New("path with mismatched type")
	ErrPathValue       = errors.New("path has no value")
	ErrUnknownType     = errors.New("path contains unknown type")
)

type Tree map[string]interface{}

func makeTreeSlice(l int) []Tree {
	a := make([]Tree, l)
	for i := 0; i < l; i++ {
		a[i] = make(Tree)
	}
	return a
}

func (tree Tree) walk(path Path, setVal ...interface{}) (val interface{}, err error) {
	doSet := len(setVal) > 0
	if len(path) == 0 {
		return nil, ErrPathTermination
	}
	node := path[0]
	if e, ok := tree[node.Name()]; ok {
		switch node.Type() {
		case ElementTypeValue:
			switch e.(type) {
			case []Tree, Tree, []interface{}:
				return nil, ErrPathType
			}
			if doSet {
				tree[node.Name()] = setVal[0]
				return nil, nil
			}
			return e, nil
		case ElementTypeArray:
			switch e.(type) {
			case []Tree, Tree:
				return nil, ErrPathType
			}
			if eX, ok := e.([]interface{}); ok {
				// Find value or extend
				if len(eX) < node.Pos()+1 {
					eX = append(eX, make([]interface{}, (node.Pos()+1)-len(eX))...)
					tree[node.Name()] = eX
					return tree.walk(path, setVal...)
				}
				if doSet {
					tree[node.Name()].([]interface{})[node.Pos()] = setVal[0]
					return nil, nil
				}
				return eX[node.Pos()], nil
			}
			return nil, ErrPathType
		case ElementTypeMap:
			if eX, ok := e.(Tree); ok {
				return eX.walk(path[1:], setVal...)
			}
			return nil, ErrPathType
		case ElementTypeMapArray:
			if eX, ok := e.([]Tree); ok {
				if len(eX) < node.Pos()+1 {
					eX = append(eX, makeTreeSlice((node.Pos()+1)-len(eX))...)
					tree[node.Name()] = eX
					return tree.walk(path, setVal...)
				}
				return eX[node.Pos()].walk(path[1:], setVal...)
			}
			return nil, ErrPathType
		default:
			return nil, ErrUnknownType
		}
	}
	switch node.Type() {
	case ElementTypeValue:
		if doSet {
			tree[node.Name()] = setVal[0]
			return nil, nil
		}
		return nil, ErrPathValue
	case ElementTypeArray:
		tree[node.Name()] = make([]interface{}, node.Pos()+1)
	case ElementTypeMap:
		tree[node.Name()] = make(Tree)
	case ElementTypeMapArray:
		tree[node.Name()] = makeTreeSlice(node.Pos() + 1)
	default:
		return nil, ErrUnknownType
	}
	return tree.walk(path, setVal...)
}
