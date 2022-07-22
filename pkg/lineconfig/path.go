package lineconfig

import (
	"fmt"
	"strconv"
	"strings"
)

type ElementType int

const (
	ElementTypeValue    ElementType = 0
	ElementTypeMap      ElementType = 1
	ElementTypeArray    ElementType = 2
	ElementTypeMapArray ElementType = 3
)

type Path []PathElement

func (p Path) IsMap() bool {
	return p[len(p)-1].Type() == ElementTypeMap || p[len(p)-1].Type() == ElementTypeMapArray
}

func (p Path) IsContinue() bool {
	return len(p) > 0 && p[0].isContinue
}

func (p Path) String() string {
	return fmt.Sprintf("%s", []PathElement(p))
}

func (p Path) ToString() string {
	r := make([]string, 0, len(p))
	for _, e := range p {
		r = append(r, e.String())
	}
	return strings.Join(r, string(pathSeparator))
}

type PathElement struct {
	pos        int
	name       string
	isMap      bool
	isContinue bool
}

func (element PathElement) Name() string {
	return element.name
}

func (element PathElement) Type() ElementType {
	if element.pos >= 0 && element.isMap {
		return ElementTypeMapArray
	}
	if element.pos >= 0 {
		return ElementTypeArray
	}
	if element.isMap {
		return ElementTypeMap
	}
	return ElementTypeValue
}

func (element PathElement) Pos() int {
	return element.pos
}

func (element PathElement) String() string {
	if element.pos >= 0 {
		return fmt.Sprintf("%s[%d]", element.Name(), element.Pos())
	}
	return element.Name()
}

func (element PathElement) ToString() string {
	var cont string
	if element.isContinue {
		cont = "."
	}
	switch element.Type() {
	case ElementTypeArray:
		return fmt.Sprintf("%s%s[%d]=", cont, element.name, element.pos)
	case ElementTypeMap:
		return fmt.Sprintf("%s%s{}", cont, element.name)
	case ElementTypeMapArray:
		return fmt.Sprintf("%s%s[%d]{}", cont, element.name, element.pos)
	default:
		return fmt.Sprintf("%s%s=", cont, element.name)
	}
}

func asElement(s string, lastIsMap bool) (element *PathElement, err error) {
	b := strings.IndexRune(s, arrayCharBegin)
	if b < 0 {
		return &PathElement{pos: -1, name: s, isMap: lastIsMap}, nil
	}
	if b == 0 {
		return nil, ErrPath
	}
	if s[len(s)-1] != arrayCharEnd {
		return nil, ErrPath
	}
	posA, err := strconv.ParseUint(s[b+1:len(s)-1], 10, 32)
	if err != nil {
		return nil, ErrPath
	}
	return &PathElement{pos: int(posA), name: s[0:b], isMap: lastIsMap}, nil
}

func endWithMap(s string) (string, bool) {
	lastIsMap := s[len(s)-1] == pathSeparator
	if lastIsMap {
		s = s[:len(s)-1]
	}
	return s, lastIsMap
}

func beginContinue(s string) (string, bool) {
	isContinue := s[0] == pathSeparator
	if isContinue {
		s = s[1:]
	}
	return s, isContinue
}

func splitPath(s string) (Path, error) {
	if len(s) < 1 {
		return nil, ErrPath
	}
	s, lastIsMap := endWithMap(s)
	s, isContinue := beginContinue(s)
	f := strings.FieldsFunc(s, func(r rune) bool { return r == pathSeparator })
	r := make(Path, 0, len(f))
	for i := 0; i < len(f); i++ {
		isMap := (i == len(f)-1 && lastIsMap) || i < len(f)-1
		e, err := asElement(f[i], isMap)
		if err != nil {
			return nil, err
		}
		r = append(r, *e)
	}
	if len(r) > 0 {
		r[0].isContinue = isContinue
	}
	return r, nil
}
